---
description: go-zero REST API 模式与最佳实践，适用于 Go 源码文件。
applyTo: "**/*.go"
---
# go-zero REST API Patterns

## Handler 层

### 正确模式
Handler 仅处理 HTTP 关注点，不含业务逻辑：
```go
func SomeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req types.SomeRequest
        if err := httpx.Parse(r, &req); err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
            return
        }
        l := logic.NewSomeLogic(r.Context(), svcCtx)
        resp, err := l.DoSomething(&req)
        if err != nil {
            httpx.ErrorCtx(r.Context(), w, err)
        } else {
            httpx.OkJsonCtx(r.Context(), w, resp)
        }
    }
}
```

### 禁止
- Handler 中查询数据库
- Handler 中做复杂校验
- 使用 `httpx.Error(w, err)` 而忽略 context
- 手动 `json.NewEncoder(w).Encode()`，应使用 `httpx.OkJsonCtx`

## Logic 层

### 正确模式
所有业务逻辑在 Logic 层实现：
```go
type SomeLogic struct {
    logx.Logger
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func (l *SomeLogic) DoSomething(req *types.SomeRequest) (*types.SomeResponse, error) {
    // 通过 svcCtx 访问 Ent Client 等依赖
    result, err := l.svcCtx.Ent.SomeEntity.Query().
        Where(someentity.FieldEQ(req.Field)).
        Only(l.ctx)
    if err != nil {
        l.Logger.Errorf("query failed: %v", err)
        return nil, err
    }
    return &types.SomeResponse{Data: result}, nil
}
```

### 要点
- 始终传递和使用 `context.Context`
- 使用内嵌的 `logx.Logger` 做结构化日志
- 通过 `svcCtx` 访问依赖，不直接创建连接
- 返回领域错误，由中间件处理 HTTP 状态码

## ServiceContext

```go
type ServiceContext struct {
    Config config.Config
    Ent    *ent.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
    return &ServiceContext{
        Config: c,
        Ent:    initDB(c),
    }
}
```
所有共享依赖（配置、数据库客户端等）统一在 ServiceContext 中管理。

## 配置

```go
type Config struct {
    rest.RestConf
    // 带默认值
    MaxSize int64 `json:",default=10485760"`
    // 可选字段
    Feature string `json:",optional"`
    // 枚举约束
    Env string `json:",default=prod,options=[dev|test|prod]"`
}
```

## API 定义（.api）

### 基本路由
```api
type CreateRequest {
    Name string `json:"name"`
}
type CreateResponse {
    Id int64 `json:"id"`
}
service forma-api {
    @handler CreateHandler
    post /api/resource (CreateRequest) returns (CreateResponse)
}
```

### 路径参数
```api
type GetRequest {
    Id int64 `path:"id"`
}
service forma-api {
    @handler GetHandler
    get /api/resource/:id (GetRequest) returns (ResourceResponse)
}
```

### JWT 保护路由
```api
@server(
    jwt: Auth
    group: protected
)
service forma-api {
    @handler ProfileHandler
    get /api/profile returns (ProfileResponse)
}
```

## 统一错误处理

### 响应结构
- 成功：`ResponseBody{Code: "200", Message: "success", Data: ...}`
- 业务错误：`ResponseBody{Code: bizErr.Code, Message: bizErr.Message}`
- 系统错误：`ResponseBody{Code: "99999", Message: "系统内部错误"}`
- 所有请求统一返回 HTTP 200

### 业务错误（`internal/errorx/`）
```go
// 使用预定义错误
return nil, errorx.ErrInvalidParam   // Code: "10001"
return nil, errorx.ErrNotFound       // Code: "10002"

// 自定义错误
return nil, errorx.NewBizError("10005", "自定义消息")
return nil, errorx.NewBizErrorf("10005", "用户 %d 不存在", userID)
```

### 错误码规范
- 格式：数字字符串
- `"0"` 成功，`"10001"`-`"10004"` 常见业务错误，`"99999"` 系统内部错误
- 新增错误码在 `internal/errorx/code.go` 中定义常量

### 禁止
- 使用 `fmt.Errorf` 返回 API 业务错误
- 在 `SetErrorHandlerCtx` 中暴露非业务错误的内部细节给客户端

## 韧性特性（生产模式自动启用）
- **熔断器**：数据库和外部调用自动保护，错误率过高时快速失败
- **负载卸除**：CPU 过高时自动拒绝请求，`Mode: pro` 时启用
- **超时控制**：通过配置 `Timeout` 设置请求超时（毫秒）
