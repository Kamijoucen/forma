# Forma API 接口测试

基于 Python unittest + requests 的接口集成测试套件。

## 前置条件

- Python 3.10+
- Forma 服务已启动（默认 `http://localhost:8888`）
- 数据库已就绪

## 安装依赖

```bash
pip install -r tests/requirements.txt
```

## 运行测试

```bash
# 运行全部测试
python -m unittest discover -s tests -p "test_*.py" -v

# 运行单个测试文件
python -m unittest tests.test_schema -v
python -m unittest tests.test_entity -v
python -m unittest tests.test_field_types -v

# 运行单个测试类
python -m unittest tests.test_schema.TestSchemaCreate -v

# 运行单个测试方法
python -m unittest tests.test_schema.TestSchemaCreate.test_create_success -v
```

## 环境变量

| 变量 | 默认值 | 说明 |
|---|---|---|
| `FORMA_BASE_URL` | `http://localhost:8888/api` | 服务地址 |
| `FORMA_TIMEOUT` | `10` | 请求超时秒数 |

## 目录结构

```
tests/
├── __init__.py
├── config.py            # 配置常量（BASE_URL、错误码等）
├── client.py            # API Client 封装（SchemaAPI / EntityAPI）
├── test_schema.py       # Schema CRUD 接口测试
├── test_entity.py       # Entity CRUD 接口测试
├── test_field_types.py  # 8 种字段类型值校验测试
├── requirements.txt
└── README.md
```

## 测试覆盖

### test_schema.py
- 创建：正常、全字段类型、空 name、空 fields、非法 type、enum 无枚举值、重复 name
- 详情：正常查询、字段一致性、不存在
- 列表：正常列表
- 更新：可变属性、不存在 Schema、修改 type、不存在字段
- 删除：正常删除、不存在
- 完整生命周期流程

### test_entity.py
- 创建：正常、含可选字段、缺必填字段、不存在 Schema、未定义字段、类型不匹配、值校验失败
- 详情：正常查询、不存在、非法 ID
- 列表：正常、分页、空列表
- 更新：正常、不存在、非法 ID、值校验失败
- 删除：正常、不存在、非法 ID
- 完整生命周期流程

### test_field_types.py
- string/text：正常值、minLength/maxLength 边界、unicode
- number：整数、浮点、负数、科学计数法、非法值
- boolean：true/false、大写/其它非法值
- date：合法格式、斜杠格式、仅日期、ISO 格式
- enum：合法/非法枚举值
- json：对象/数组/嵌套/非法 JSON
- array：各类合法数组/非数组 JSON/非法 JSON

## 数据清理

测试使用 UUID 前缀生成唯一 Schema 名称，在 `tearDown`/`tearDownClass` 中通过 API 自动清理测试数据。
