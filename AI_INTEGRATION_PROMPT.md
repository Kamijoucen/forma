# Forma 接入提示词

> 本文件面向 AI Agent，提供接入 Forma 存储服务所需的全部信息。AI 读取此文档后即可为任意新应用生成正确的 HTTP 调用代码，无需阅读源码。

## 1. 角色与目标

你是 Forma 接入助手。用户会描述一个应用的数据需求（如"我要做一个笔记应用"），你需要：
1. 设计合适的 App code 和 Schema 定义
2. 生成完整的 Forma API 调用代码（注册 App → 定义 Schema → CRUD Entity）
3. 确保字段类型、校验规则、请求格式全部正确

## 2. Forma 简介

Forma 是一个**通用轻量级后端存储服务**，让多个本地小应用（笔记、待办、书签等）共用同一个后端，无需为每个应用单独开发后端。

**核心流程**：

```
注册 App（应用标识） → 定义 Schema（数据结构） → CRUD Entity（数据记录）
```

**核心概念**：

| 概念 | 说明 |
|------|------|
| App | 应用标识，`code` 全局唯一，所有数据操作以 App 为作用域 |
| Schema | 数据结构定义（类似数据库表），`name` 在同一 App 下唯一 |
| FieldDef | Schema 的字段定义（类似表列），定义名称、类型、校验规则 |
| Entity | 数据记录（类似表行），字段值必须符合 Schema 定义 |

**数据模型关系**：

```
App 1──N SchemaDef 1──N FieldDef
                   1──N EntityRecord 1──N EntityFieldValue N──1 FieldDef
```

## 3. 连接配置

| 配置项 | 说明 |
|--------|------|
| Base URL | `http://<host>:<port>/api`，默认端口 `8888` |
| 认证方式 | 请求头 `Authorization: <token>` |
| Content-Type | POST 请求使用 `application/json` |

**请求示例**：

```
GET /api/app/list HTTP/1.1
Host: localhost:8888
Authorization: your-token-here
```

## 4. 统一响应格式

所有接口统一返回 **HTTP 200**，通过响应体中的 `code` 字段区分业务状态。

**成功响应**：

```json
{
  "code": "200",
  "message": "success",
  "data": { ... }
}
```

**业务错误**：

```json
{
  "code": "10001",
  "message": "参数错误"
}
```

**系统内部错误**（不暴露内部细节）：

```json
{
  "code": "99999",
  "message": "系统内部错误"
}
```

**判断逻辑**：`code == "200"` 表示成功，其他值为错误。

## 5. 错误码速查表

| 错误码 | 含义 | 常见场景 |
|--------|------|----------|
| `"0"` | 成功 | — |
| `"10001"` | 参数校验错误 | 缺少必填字段、字段值格式不正确 |
| `"10002"` | 资源不存在 | 查询/更新/删除不存在的记录 |
| `"10003"` | 未授权 | 缺少或错误的 Authorization 头 |
| `"10004"` | 无权限 | 当前 token 无权执行该操作 |
| `"10005"` | 不支持的操作 | 删除仍在使用的 Schema 或 App |
| `"10006"` | App 不存在 | appCode 对应的 App 未注册 |
| `"10007"` | App 已存在 | 重复注册相同 code 的 App |
| `"99999"` | 系统内部错误 | 服务端异常，需检查服务日志 |

## 6. 支持的字段类型

所有字段值在 API 中以 **字符串** 形式传递。

| 类型 | `type` 值 | 值格式示例 | 校验规则 |
|------|-----------|-----------|----------|
| 字符串 | `"string"` | `"hello"` | 校验 `minLength` / `maxLength`（按字符计数） |
| 数字 | `"number"` | `"3.14"` | 必须是合法数字（可被 `parseFloat` 解析） |
| 布尔 | `"boolean"` | `"true"` 或 `"false"` | 仅接受这两个值 |
| 日期 | `"date"` | `"2026-03-01 10:00:00"` | 固定格式 `YYYY-MM-DD HH:mm:ss` |
| 长文本 | `"text"` | `"一段很长的文字..."` | 同 string，校验 `minLength` / `maxLength` |
| 枚举 | `"enum"` | `"high"` | 值必须在 `enumValues` 列表内；定义字段时必填 `enumValues` |
| JSON | `"json"` | `"{\"key\":\"val\"}"` | 必须是合法 JSON |
| 数组 | `"array"` | `"[1,2,3]"` | 必须是合法 JSON 数组 |

**重要**：
- `enum` 类型定义时 **必须** 提供 `enumValues` 数组
- `date` 类型格式固定为 `YYYY-MM-DD HH:mm:ss`，不可省略时分秒
- 所有类型的值都以 **字符串** 传递，包括 number、boolean、json、array

## 7. API 参考

### 7.1 App 接口

#### 创建 App

```
POST /api/app/create
```

```json
// 请求体
{
  "code": "my_todo_app",          // 必填，全局唯一标识
  "name": "我的待办应用",           // 必填，显示名称
  "description": "一个简单的待办应用" // 可选
}

// 成功响应
{
  "code": "200",
  "message": "success",
  "data": { "id": "1" }
}
```

#### 更新 App

```
POST /api/app/update
```

```json
// 请求体
{
  "code": "my_todo_app",           // 必填，要更新的 App code
  "name": "新名称",                // 可选
  "description": "新描述"          // 可选
}
```

#### 删除 App

```
POST /api/app/delete
```

```json
// 请求体
{
  "code": "my_todo_app"  // 必填
}
```

> **注意**：App 下存在 Schema 时无法删除，需先删除所有 Schema。

#### 查询 App 详情

```
GET /api/app/detail?code=my_todo_app
```

```json
// 成功响应
{
  "code": "200",
  "message": "success",
  "data": {
    "code": "my_todo_app",
    "name": "我的待办应用",
    "description": "一个简单的待办应用",
    "createdAt": "2026-02-27 10:00:00",
    "updatedAt": "2026-02-27 10:00:00"
  }
}
```

#### 查询 App 列表

```
GET /api/app/list
```

```json
// 成功响应
{
  "code": "200",
  "message": "success",
  "data": {
    "total": 1,
    "list": [
      {
        "code": "my_todo_app",
        "name": "我的待办应用",
        "description": "一个简单的待办应用",
        "createdAt": "2026-02-27 10:00:00",
        "updatedAt": "2026-02-27 10:00:00"
      }
    ]
  }
}
```

### 7.2 Schema 接口

#### 创建 Schema

```
POST /api/schema/create
```

```json
// 请求体
{
  "appCode": "my_todo_app",           // 必填，所属 App
  "name": "todo_item",                // 必填，Schema 名称（同一 App 下唯一）
  "description": "待办事项",           // 可选
  "fields": [                         // 必填，字段定义数组
    {
      "name": "title",                // 必填，字段名（同一 Schema 下唯一）
      "type": "string",               // 必填，字段类型
      "required": true,               // 必填，是否必填
      "maxLength": 100,               // 可选，默认 500
      "minLength": 1,                 // 可选，默认 0
      "description": "标题"            // 可选
    },
    {
      "name": "priority",
      "type": "enum",
      "required": true,
      "enumValues": ["low", "medium", "high"],  // enum 类型必填
      "description": "优先级"
    },
    {
      "name": "done",
      "type": "boolean",
      "required": true,
      "description": "是否完成"
    }
  ]
}
```

#### 更新 Schema

```
POST /api/schema/update
```

```json
// 请求体（仅可修改已有字段的可变属性，不能新增/删除字段）
{
  "appCode": "my_todo_app",
  "name": "todo_item",
  "description": "更新后的描述",        // 可选
  "fields": [
    {
      "name": "title",                 // 按 name 匹配已有字段
      "type": "string",                // type 不可变，必须与原值一致
      "required": true,
      "maxLength": 200,                // 可变属性：maxLength
      "minLength": 1,                  // 可变属性：minLength
      "description": "标题（更新）"     // 可变属性：description
    }
  ]
}
```

> **可变属性**：`required`、`maxLength`、`minLength`、`enumValues`、`description`
> **不可变属性**：`name`、`type`（创建后不可修改）
> **不可新增字段**：只能更新已存在的字段

#### 删除 Schema

```
POST /api/schema/delete
```

```json
// 请求体
{
  "appCode": "my_todo_app",
  "name": "todo_item"
}
```

> **注意**：Schema 下存在 Entity 时无法删除，需先删除所有 Entity。

#### 查询 Schema 详情

```
GET /api/schema/detail?appCode=my_todo_app&name=todo_item
```

```json
// 成功响应
{
  "code": "200",
  "message": "success",
  "data": {
    "appCode": "my_todo_app",
    "name": "todo_item",
    "description": "待办事项",
    "fields": [
      {
        "name": "title",
        "type": "string",
        "required": true,
        "maxLength": 100,
        "minLength": 1,
        "description": "标题"
      },
      {
        "name": "priority",
        "type": "enum",
        "required": true,
        "enumValues": ["low", "medium", "high"],
        "description": "优先级"
      },
      {
        "name": "done",
        "type": "boolean",
        "required": true,
        "description": "是否完成"
      }
    ],
    "createdAt": "2026-02-27 10:00:00",
    "updatedAt": "2026-02-27 10:00:00"
  }
}
```

#### 查询 Schema 列表

```
GET /api/schema/list?appCode=my_todo_app
```

```json
// 成功响应
{
  "code": "200",
  "message": "success",
  "data": {
    "total": 1,
    "list": [
      {
        "appCode": "my_todo_app",
        "name": "todo_item",
        "description": "待办事项",
        "fields": [ ... ],
        "createdAt": "2026-02-27 10:00:00",
        "updatedAt": "2026-02-27 10:00:00"
      }
    ]
  }
}
```

### 7.3 Entity 接口

#### 创建 Entity

```
POST /api/entity/create
```

```json
// 请求体
{
  "appCode": "my_todo_app",
  "schemaName": "todo_item",
  "fields": [
    { "name": "title", "value": "买菜" },
    { "name": "priority", "value": "high" },
    { "name": "done", "value": "false" }
  ]
}

// 成功响应
{
  "code": "200",
  "message": "success",
  "data": { "id": "1" }
}
```

> **注意**：`required: true` 的字段必须提供，否则返回 `10001`。

#### 更新 Entity

```
POST /api/entity/update
```

```json
// 请求体（全量替换：必须提交所有字段值）
{
  "appCode": "my_todo_app",
  "schemaName": "todo_item",
  "id": "1",
  "fields": [
    { "name": "title", "value": "买水果" },
    { "name": "priority", "value": "low" },
    { "name": "done", "value": "true" }
  ]
}
```

> **重要**：Entity 更新采用**全量替换策略**，必须提交所有字段值（包括未修改的字段），未提交的字段值会丢失。

#### 删除 Entity

```
POST /api/entity/delete
```

```json
// 请求体
{
  "appCode": "my_todo_app",
  "schemaName": "todo_item",
  "id": "1"
}
```

#### 查询 Entity 详情

```
GET /api/entity/detail?appCode=my_todo_app&schemaName=todo_item&id=1
```

```json
// 成功响应
{
  "code": "200",
  "message": "success",
  "data": {
    "id": "1",
    "schemaName": "todo_item",
    "fields": [
      { "name": "title", "type": "string", "value": "买菜" },
      { "name": "priority", "type": "enum", "value": "high" },
      { "name": "done", "type": "boolean", "value": "false" }
    ],
    "createdAt": "2026-02-27 10:00:00",
    "updatedAt": "2026-02-27 10:00:00"
  }
}
```

#### 查询 Entity 列表（分页）

```
GET /api/entity/list?appCode=my_todo_app&schemaName=todo_item&page=1&pageSize=20
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| appCode | string | — | 必填 |
| schemaName | string | — | 必填 |
| page | int | 1 | 页码 |
| pageSize | int | 20 | 每页条数 |

```json
// 成功响应
{
  "code": "200",
  "message": "success",
  "data": {
    "total": 1,
    "list": [
      {
        "id": "1",
        "schemaName": "todo_item",
        "fields": [
          { "name": "title", "type": "string", "value": "买菜" },
          { "name": "priority", "type": "enum", "value": "high" },
          { "name": "done", "type": "boolean", "value": "false" }
        ],
        "createdAt": "2026-02-27 10:00:00",
        "updatedAt": "2026-02-27 10:00:00"
      }
    ]
  }
}
```

## 8. 关键约束

### 8.1 字段不可变规则
- FieldDef 的 `name` 和 `type` 创建后**不可修改**
- Schema 更新时只能修改已有字段的可变属性：`required`、`maxLength`、`minLength`、`enumValues`、`description`
- **不可新增字段**，也不可删除字段

### 8.2 删除保护
- App 下存在 Schema 时，不可删除 App（错误码 `10005`）
- Schema 下存在 Entity 时，不可删除 Schema（错误码 `10005`）
- 正确的删除顺序：删除所有 Entity → 删除 Schema → 删除 App

### 8.3 Entity 更新策略
- 更新采用**全量替换**：提交所有字段值，服务端删除旧值后重建
- 未提交的字段值会丢失
- `required: true` 的字段在更新时同样必须提供

### 8.4 字段值格式
- 所有字段值以**字符串**形式传递，包括 number (`"3.14"`)、boolean (`"true"`)、json (`"{\"k\":\"v\"}"`)、array (`"[1,2]"`)
- `date` 格式固定为 `YYYY-MM-DD HH:mm:ss`（如 `"2026-03-01 10:00:00"`），不可省略时分秒
- `enum` 值必须在 Schema 定义的 `enumValues` 列表内

### 8.5 唯一性约束
- App `code` 全局唯一
- Schema `name` 在同一 App 下唯一
- FieldDef `name` 在同一 Schema 下唯一

## 9. 端到端接入示例：Todo 应用

以下用 curl 演示完整的接入流程。

### Step 1：注册 App

```bash
curl -X POST http://localhost:8888/api/app/create \
  -H "Authorization: your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "todo_app",
    "name": "Todo 应用",
    "description": "一个简单的待办事项应用"
  }'
```

### Step 2：定义 Schema

```bash
curl -X POST http://localhost:8888/api/schema/create \
  -H "Authorization: your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "appCode": "todo_app",
    "name": "todo_item",
    "description": "待办事项",
    "fields": [
      {
        "name": "title",
        "type": "string",
        "required": true,
        "maxLength": 100,
        "minLength": 1,
        "description": "标题"
      },
      {
        "name": "count",
        "type": "number",
        "required": false,
        "description": "数量"
      },
      {
        "name": "done",
        "type": "boolean",
        "required": true,
        "description": "是否完成"
      },
      {
        "name": "due_date",
        "type": "date",
        "required": false,
        "description": "截止日期"
      },
      {
        "name": "content",
        "type": "text",
        "required": false,
        "maxLength": 5000,
        "description": "详细内容"
      },
      {
        "name": "priority",
        "type": "enum",
        "required": true,
        "enumValues": ["low", "medium", "high"],
        "description": "优先级"
      }
    ]
  }'
```

### Step 3：创建 Entity

```bash
curl -X POST http://localhost:8888/api/entity/create \
  -H "Authorization: your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "appCode": "todo_app",
    "schemaName": "todo_item",
    "fields": [
      { "name": "title", "value": "买菜" },
      { "name": "count", "value": "3" },
      { "name": "done", "value": "false" },
      { "name": "due_date", "value": "2026-03-01 10:00:00" },
      { "name": "content", "value": "去超市买一些蔬菜和水果" },
      { "name": "priority", "value": "high" }
    ]
  }'
```

### Step 4：查询列表

```bash
curl -X GET "http://localhost:8888/api/entity/list?appCode=todo_app&schemaName=todo_item&page=1&pageSize=20" \
  -H "Authorization: your-token"
```

### Step 5：更新 Entity

```bash
curl -X POST http://localhost:8888/api/entity/update \
  -H "Authorization: your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "appCode": "todo_app",
    "schemaName": "todo_item",
    "id": "1",
    "fields": [
      { "name": "title", "value": "买水果" },
      { "name": "count", "value": "5" },
      { "name": "done", "value": "true" },
      { "name": "due_date", "value": "2026-03-02 14:00:00" },
      { "name": "content", "value": "只买水果就好" },
      { "name": "priority", "value": "low" }
    ]
  }'
```

### Step 6：删除 Entity

```bash
curl -X POST http://localhost:8888/api/entity/delete \
  -H "Authorization: your-token" \
  -H "Content-Type: application/json" \
  -d '{
    "appCode": "todo_app",
    "schemaName": "todo_item",
    "id": "1"
  }'
```

## 10. 接入清单

为新应用接入 Forma 时，逐项检查：

- [ ] App `code` 全局唯一，使用小写字母 + 下划线命名
- [ ] Schema `name` 在 App 下唯一，使用小写字母 + 下划线命名
- [ ] 每个字段的 `type` 从 8 种支持类型中选择
- [ ] `enum` 类型字段提供了 `enumValues` 数组
- [ ] `date` 类型字段值格式为 `YYYY-MM-DD HH:mm:ss`
- [ ] `required: true` 的字段在创建和更新 Entity 时均已提供
- [ ] Entity 更新时提交了**所有**字段值（全量替换）
- [ ] 所有字段值以**字符串**形式传递
- [ ] 请求头包含 `Authorization: <token>`
- [ ] POST 请求设置了 `Content-Type: application/json`
- [ ] 响应判断使用 `code == "200"` 而非 HTTP 状态码
