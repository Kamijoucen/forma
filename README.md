# Forma

Forma 是一个通用轻量级后端存储服务，让多个本地小应用（笔记、待办、书签等）共用同一个后端，无需为每个应用单独开发后端服务。

## 核心思路

1. **定义 Schema** — 通过 API 描述你的数据结构（有哪些字段、什么类型、是否必填等）
2. **读写数据** — 基于已定义的 Schema，直接通过统一 API 进行增删改查

所有数据按 Schema 隔离，不同应用各自定义自己的 Schema，互不干扰。

## 使用示例

### 创建一个 "todo" Schema

```
POST /api/schema/create
```

```json
{
  "name": "todo",
  "displayName": "待办事项",
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
GET /api/entity/list?schemaName=todo&page=1&pageSize=20
```

支持的字段类型：`string` `number` `boolean` `date` `text` `enum` `json` `array`

## 技术选型

- 语言：Go
- Web 框架：[go-zero](https://go-zero.dev)
- ORM：[Ent](https://entgo.io)
- 数据库：PostgreSQL

## License

MIT
