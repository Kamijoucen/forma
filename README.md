# Forma

Forma 是一个通用轻量级后端存储服务，让多个本地小应用（笔记、待办、书签等）共用同一个后端，无需为每个应用单独开发后端服务。

## 核心思路

1. **注册 App** — 每个应用注册为一个 App，通过唯一的 `appCode` 标识
2. **定义 Schema** — 在 App 下通过 API 描述你的数据结构（有哪些字段、什么类型、是否必填等）
3. **读写数据** — 基于已定义的 Schema，直接通过统一 API 进行增删改查

所有数据按 App 隔离，不同应用各自定义自己的 Schema，互不干扰。同一 App 下 Schema 名称唯一，不同 App 可以有同名 Schema。

## 使用示例

### 注册一个 App

```
POST /api/app/create
```

```json
{
  "code": "my-todo",
  "name": "我的待办应用",
  "description": "一个简单的待办事项应用"
}
```

### 创建一个 "todo" Schema

```
POST /api/schema/create
```

```json
{
  "appCode": "my-todo",
  "name": "todo",
  "fields": [
    { "name": "title", "type": "string", "required": true, "maxLength": 200 },
    { "name": "done", "type": "boolean", "defaultValue": "false" },
    { "name": "priority", "type": "enum", "enumValues": ["low", "medium", "high"] }
  ]
}
```

### 写入一条待办

```
POST /api/entity/create
```

```json
{
  "appCode": "my-todo",
  "schemaName": "todo",
  "fields": [
    { "name": "title", "value": "买牛奶" },
    { "name": "done", "value": "false" },
    { "name": "priority", "value": "medium" }
  ]
}
```

### 查询待办列表

```
GET /api/entity/list?appCode=my-todo&schemaName=todo&page=1&pageSize=20
```

支持的字段类型：`string` `number` `boolean` `date` `text` `enum` `json` `array`

## 技术选型

- 语言：Go
- Web 框架：[go-zero](https://go-zero.dev)
- ORM：[Ent](https://entgo.io)
- 数据库：PostgreSQL

## License

MIT
