# Forma AI Instructions

## 规则

### Spec-First
- 修改接口前先更新 `.api` 文件
- 通过 goctl 生成代码，不手动编写 handler/types 样板
- 生成命令：`goctl api go -style go_zero -api forma.api -dir .`

### go-zero 约定
- Context 优先：`func(ctx context.Context, req *types.Request)`
- 错误处理：使用统一响应结构（见 `internal/svc/response.go`）
- 业务错误：使用 `internal/errorx` 包中的 `BizError`，不使用 `fmt.Errorf` 返回 API 错误
- 配置默认值：`json:",default=value"`
- 可选字段：`json:",optional"`

### Ent ORM
- Schema 定义在 `internal/ent/schema/`
- 生成命令：`go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert ./ent/schema`（在 `internal/` 下执行）
- 新建实体：`go run -mod=mod entgo.io/ent/cmd/ent new <Name>`（在 `internal/` 下执行）
- 不使用 go-zero 自带的 sqlx/model 生成

### 三层架构
```
HTTP Request → Handler（HTTP 关注点） → Logic（业务逻辑） → Ent Client
                                          ↓
                                   ServiceContext（依赖注入）
                                          ↓
                                   Service（公共业务函数）
```
- Handler：仅处理请求解析和响应，不含业务逻辑
- Logic：业务编排，通过 `svcCtx` 访问依赖，调用 service 层公共函数
- Service（`internal/service/`）：跨 Logic 复用的公共业务函数，纯函数，不建 struct，如校验、实体转换等
- ServiceContext：统一管理配置、数据库连接等依赖

## 工作流

### 新增 API 端点
1. 在 `forma.api` 中添加类型和路由定义
2. 运行 goctl 生成 handler/logic/types
3. 在 logic 层实现业务逻辑

### 修改 API
1. 编辑 `forma.api`
2. 重新运行 goctl（不会覆盖已有 logic 实现）
3. 更新 logic 层代码

### 新增数据实体
1. 运行 `ent new <Name>` 创建 schema
2. 在 `internal/ent/schema/` 中定义字段和边
3. 运行 ent generate 生成代码
4. 在 ServiceContext 中注入 Ent Client

### 错误处理
- 业务错误定义在 `internal/errorx/` 包
- `BizError` 携带 `Code`（数字字符串）和 `Message`
- 错误码常量定义在 `internal/errorx/code.go`，格式为数字字符串（如 `"10001"`）
- Logic 层返回业务错误：使用预定义变量（如 `errorx.ErrNotFound`）或 `errorx.NewBizError(code, msg)`
- `SetErrorHandlerCtx` 通过 `errors.As` 区分 `BizError` 和普通错误
- `BizError` → 返回其 Code/Message；普通错误 → 返回 `CodeInternal`（`"99999"`）+ 通用提示，日志记录原始错误
- 所有请求统一返回 HTTP 200，通过 `ResponseBody.Code` 区分业务状态

## 避免
- Handler 中写业务逻辑
- 绕过 Ent 直接写 SQL
- 手动编写 goctl 应生成的样板代码
- 使用 `fmt.Errorf` 返回 API 错误（使用 `errorx.BizError`）
- 非业务错误暴露内部信息给客户端
