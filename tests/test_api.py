"""
Forma API 接口测试脚本

使用方式:
    1. 安装依赖: pip install -r requirements.txt
    2. 启动 Forma 服务确保 localhost:8888 可访问
    3. 运行: python test_api.py
    4. 如需只测试部分功能，注释掉 main 中不需要的调用即可
"""

import os
import time
import requests

# ========== 配置 ==========
BASE_URL = os.getenv("FORMA_BASE_URL", "http://localhost:8888/api")
TEMP_TOKEN = os.getenv("FORMA_TOKEN", "123")
HEADERS = {"Authorization": TEMP_TOKEN}

# 测试用的 App 标识
APP_CODE = f"test_app_{int(time.time())}"
APP_NAME = "测试应用"

# 测试用的 Schema 字段定义 —— 覆盖 string/number/boolean/date/text/enum 六种类型
SCHEMA_FIELDS = [
    {
        "name": "title",
        "type": "string",
        "required": True,
        "maxLength": 100,
        "minLength": 1,
        "description": "标题",
    },
    {
        "name": "count",
        "type": "number",
        "required": False,
        "description": "数量",
    },
    {
        "name": "done",
        "type": "boolean",
        "required": True,
        "description": "是否完成",
    },
    {
        "name": "due_date",
        "type": "date",
        "required": False,
        "description": "截止日期",
    },
    {
        "name": "content",
        "type": "text",
        "required": False,
        "maxLength": 5000,
        "description": "详细内容",
    },
    {
        "name": "priority",
        "type": "enum",
        "required": True,
        "enumValues": ["low", "medium", "high"],
        "description": "优先级",
    },
]

# 测试用的 Entity 字段值 —— 与上面的 SCHEMA_FIELDS 对应
ENTITY_FIELDS = [
    {"name": "title", "type": "string", "value": "买菜"},
    {"name": "count", "type": "number", "value": "3"},
    {"name": "done", "type": "boolean", "value": "false"},
    {"name": "due_date", "type": "date", "value": "2026-03-01 10:00:00"},
    {"name": "content", "type": "text", "value": "去超市买一些蔬菜和水果"},
    {"name": "priority", "type": "enum", "value": "high"},
]

# 更新后的 Entity 字段值
ENTITY_FIELDS_UPDATED = [
    {"name": "title", "type": "string", "value": "买水果"},
    {"name": "count", "type": "number", "value": "5"},
    {"name": "done", "type": "boolean", "value": "true"},
    {"name": "due_date", "type": "date", "value": "2026-03-02 14:00:00"},
    {"name": "content", "type": "text", "value": "只买水果就好"},
    {"name": "priority", "type": "enum", "value": "low"},
]

# ========== 通过/失败计数 ==========
_passed = 0
_failed = 0


# ========== 工具函数 ==========


def gen_schema_name():
    """生成唯一的 Schema 名称，避免与已有数据冲突"""
    return f"test_schema_{int(time.time())}"


def post(path, body=None):
    """发送 POST 请求，返回解析后的 JSON"""
    url = f"{BASE_URL}{path}"
    resp = requests.post(url, json=body, headers=HEADERS, timeout=10)
    print(f"  Response: {resp.text}")
    return resp.json()


def get(path, params=None):
    """发送 GET 请求，返回解析后的 JSON"""
    url = f"{BASE_URL}{path}"
    resp = requests.get(url, params=params, headers=HEADERS, timeout=10)
    return resp.json()


def assert_success(resp, label):
    """断言响应成功（code == "200"），打印结果"""
    global _passed, _failed
    print(f"  Response: {resp}")
    code = resp.get("code")
    if code == "200":
        _passed += 1
        print(f"  [PASS] {label}")
        return True
    else:
        _failed += 1
        print(f"  [FAIL] {label}")
        print(f"         期望 code=\"200\", 实际: {resp}")
        return False


def assert_code(resp, expected_code, label):
    """断言响应返回特定业务码"""
    global _passed, _failed
    code = resp.get("code")
    if code == expected_code:
        _passed += 1
        print(f"  [PASS] {label}")
        return True
    else:
        _failed += 1
        print(f"  [FAIL] {label}")
        print(f"         期望 code=\"{expected_code}\", 实际: {resp}")
        return False


def assert_equal(actual, expected, label):
    """断言两个值相等"""
    global _passed, _failed
    if actual == expected:
        _passed += 1
        print(f"  [PASS] {label}")
        return True
    else:
        _failed += 1
        print(f"  [FAIL] {label}")
        print(f"         期望: {expected}, 实际: {actual}")
        return False


def assert_true(condition, label):
    """断言条件为真"""
    global _passed, _failed
    if condition:
        _passed += 1
        print(f"  [PASS] {label}")
        return True
    else:
        _failed += 1
        print(f"  [FAIL] {label}")
        return False


def print_summary():
    """打印测试汇总"""
    total = _passed + _failed
    print("\n" + "=" * 50)
    print(f"测试完成: 共 {total} 项, 通过 {_passed} 项, 失败 {_failed} 项")
    if _failed == 0:
        print("All tests passed!")
    else:
        print(f"{_failed} test(s) FAILED")
    print("=" * 50)


# ========== App 测试 ==========


def test_app_create(app_code):
    """创建 App"""
    print(f"\n--- test_app_create (code={app_code}) ---")
    body = {
        "code": app_code,
        "name": APP_NAME,
        "description": "接口测试自动创建的App",
    }
    resp = post("/app/create", body)
    assert_success(resp, "创建 App")

    data = resp.get("data", {})
    app_id = data.get("id", "")
    assert_true(len(app_id) > 0, f"返回了 app id: {app_id}")
    return app_id


def test_app_create_duplicate(app_code):
    """重复创建 App 应返回已存在错误"""
    print(f"\n--- test_app_create_duplicate (code={app_code}) ---")
    body = {
        "code": app_code,
        "name": "重复的App",
    }
    resp = post("/app/create", body)
    assert_code(resp, "10007", "重复创建 App 返回 10007")


def test_app_detail(app_code):
    """查询 App 详情"""
    print(f"\n--- test_app_detail (code={app_code}) ---")
    resp = get("/app/detail", {"code": app_code})
    assert_success(resp, "查询 App 详情")

    data = resp.get("data", {})
    assert_equal(data.get("code"), app_code, "App code 一致")
    assert_equal(data.get("name"), APP_NAME, "App name 一致")
    return data


def test_app_list(app_code):
    """查询 App 列表，验证包含指定 App"""
    print(f"\n--- test_app_list (code={app_code}) ---")
    resp = get("/app/list")
    assert_success(resp, "查询 App 列表")

    data = resp.get("data", {})
    total = data.get("total", 0)
    assert_true(total >= 1, f"App 总数 >= 1 (实际: {total})")

    codes = [a.get("code") for a in data.get("list", [])]
    assert_true(app_code in codes, f"列表中包含 {app_code}")


def test_app_update(app_code):
    """更新 App 名称和描述"""
    print(f"\n--- test_app_update (code={app_code}) ---")
    body = {
        "code": app_code,
        "name": "更新后的应用名",
        "description": "更新后的描述",
    }
    resp = post("/app/update", body)
    assert_success(resp, "更新 App")

    # 查询验证
    resp = get("/app/detail", {"code": app_code})
    assert_success(resp, "查询更新后的 App 详情")
    data = resp.get("data", {})
    assert_equal(data.get("name"), "更新后的应用名", "App name 已更新")


def test_app_delete(app_code):
    """删除 App"""
    print(f"\n--- test_app_delete (code={app_code}) ---")
    resp = post("/app/delete", {"code": app_code})
    assert_success(resp, "删除 App")

    # 查询已删除的 App 应返回 10006
    resp = get("/app/detail", {"code": app_code})
    assert_code(resp, "10006", "已删除 App 查询返回 10006")


def test_app_delete_in_use(app_code):
    """有 Schema 关联时删除 App 应被拒绝"""
    print(f"\n--- test_app_delete_in_use (code={app_code}) ---")
    resp = post("/app/delete", {"code": app_code})
    assert_code(resp, "10005", "App 下有 Schema 时删除返回 10005")


# ========== Schema 测试 ==========


def test_schema_create(app_code, schema_name):
    """创建 Schema —— 包含 string/number/boolean/date/text/enum 六种字段类型"""
    print(f"\n--- test_schema_create (app={app_code}, name={schema_name}) ---")
    body = {
        "appCode": app_code,
        "name": schema_name,
        "description": "接口测试自动创建的Schema",
        "fields": SCHEMA_FIELDS,
    }
    resp = post("/schema/create", body)
    assert_success(resp, "创建 Schema")


def test_schema_detail(app_code, schema_name):
    """查询 Schema 详情，验证字段列表与提交一致"""
    print(f"\n--- test_schema_detail (app={app_code}, name={schema_name}) ---")
    resp = get("/schema/detail", {"appCode": app_code, "name": schema_name})
    assert_success(resp, "查询 Schema 详情")

    data = resp.get("data", {})
    assert_equal(data.get("name"), schema_name, "Schema name 一致")
    assert_equal(data.get("appCode"), app_code, "Schema appCode 一致")

    fields = data.get("fields", [])
    expected_names = {f["name"] for f in SCHEMA_FIELDS}
    actual_names = {f["name"] for f in fields}
    assert_equal(actual_names, expected_names, "字段名称集合一致")

    return data


def test_schema_list(app_code, schema_name):
    """查询 Schema 列表，验证包含指定 Schema"""
    print(f"\n--- test_schema_list (app={app_code}, name={schema_name}) ---")
    resp = get("/schema/list", {"appCode": app_code})
    assert_success(resp, "查询 Schema 列表")

    data = resp.get("data", {})
    total = data.get("total", 0)
    assert_true(total >= 1, f"Schema 总数 >= 1 (实际: {total})")

    names = [s.get("name") for s in data.get("list", [])]
    assert_true(schema_name in names, f"列表中包含 {schema_name}")


def test_schema_update(app_code, schema_name):
    """更新 Schema 的 description 和字段属性，再查询验证"""
    print(f"\n--- test_schema_update (app={app_code}, name={schema_name}) ---")

    # 更新字段：修改 title 的 maxLength，修改 done 的 required
    updated_fields = []
    for f in SCHEMA_FIELDS:
        field = dict(f)
        if field["name"] == "title":
            field["maxLength"] = 200
        if field["name"] == "done":
            field["required"] = False
        updated_fields.append(field)

    body = {
        "appCode": app_code,
        "name": schema_name,
        "description": "更新后的描述",
        "fields": updated_fields,
    }
    resp = post("/schema/update", body)
    assert_success(resp, "更新 Schema")

    # 查询验证更新生效
    resp = get("/schema/detail", {"appCode": app_code, "name": schema_name})
    assert_success(resp, "查询更新后的 Schema 详情")

    data = resp.get("data", {})
    assert_equal(data.get("description"), "更新后的描述", "description 已更新")

    # 验证字段属性更新
    fields_map = {f["name"]: f for f in data.get("fields", [])}
    if "title" in fields_map:
        assert_equal(fields_map["title"].get("maxLength"), 200, "title.maxLength 已更新为 200")
    if "done" in fields_map:
        assert_equal(fields_map["done"].get("required"), False, "done.required 已更新为 false")


def test_schema_delete(app_code, schema_name):
    """删除 Schema，随后查询确认返回资源不存在错误"""
    print(f"\n--- test_schema_delete (app={app_code}, name={schema_name}) ---")
    resp = post("/schema/delete", {"appCode": app_code, "name": schema_name})
    assert_success(resp, "删除 Schema")

    # 查询已删除的 Schema 应返回 10002（资源不存在）
    resp = get("/schema/detail", {"appCode": app_code, "name": schema_name})
    assert_code(resp, "10002", "已删除 Schema 查询返回 10002")


# ========== Entity 测试 ==========


def test_entity_create(app_code, schema_name):
    """创建实体记录，返回实体 ID"""
    print(f"\n--- test_entity_create (app={app_code}, schema={schema_name}) ---")
    body = {
        "appCode": app_code,
        "schemaName": schema_name,
        "fields": ENTITY_FIELDS,
    }
    resp = post("/entity/create", body)
    assert_success(resp, "创建 Entity")

    data = resp.get("data", {})
    entity_id = data.get("id", "")
    assert_true(len(entity_id) > 0, f"返回了 entity id: {entity_id}")

    return entity_id


def test_entity_detail(app_code, schema_name, entity_id):
    """查询实体详情，验证字段值正确"""
    print(f"\n--- test_entity_detail (app={app_code}, schema={schema_name}, id={entity_id}) ---")
    resp = get("/entity/detail", {"appCode": app_code, "schemaName": schema_name, "id": entity_id})
    assert_success(resp, "查询 Entity 详情")

    data = resp.get("data", {})
    assert_equal(data.get("id"), entity_id, "Entity ID 一致")
    assert_equal(data.get("schemaName"), schema_name, "Schema name 一致")

    # 验证字段值
    fields_map = {f["name"]: f["value"] for f in data.get("fields", [])}
    expected_map = {f["name"]: f["value"] for f in ENTITY_FIELDS}
    for name, expected_val in expected_map.items():
        actual_val = fields_map.get(name)
        assert_equal(actual_val, expected_val, f"字段 {name} 值一致")

    return data


def test_entity_list(app_code, schema_name):
    """查询实体列表，验证 total >= 1"""
    print(f"\n--- test_entity_list (app={app_code}, schema={schema_name}) ---")
    resp = get("/entity/list", {"appCode": app_code, "schemaName": schema_name, "page": 1, "pageSize": 20})
    assert_success(resp, "查询 Entity 列表")

    data = resp.get("data", {})
    total = data.get("total", 0)
    assert_true(total >= 1, f"Entity 总数 >= 1 (实际: {total})")


def test_entity_update(app_code, schema_name, entity_id):
    """更新实体字段值，再查询验证"""
    print(f"\n--- test_entity_update (app={app_code}, schema={schema_name}, id={entity_id}) ---")
    body = {
        "appCode": app_code,
        "schemaName": schema_name,
        "id": entity_id,
        "fields": ENTITY_FIELDS_UPDATED,
    }
    resp = post("/entity/update", body)
    assert_success(resp, "更新 Entity")

    # 查询验证更新生效
    resp = get("/entity/detail", {"appCode": app_code, "schemaName": schema_name, "id": entity_id})
    assert_success(resp, "查询更新后的 Entity 详情")

    fields_map = {f["name"]: f["value"] for f in resp.get("data", {}).get("fields", [])}
    expected_map = {f["name"]: f["value"] for f in ENTITY_FIELDS_UPDATED}
    for name, expected_val in expected_map.items():
        actual_val = fields_map.get(name)
        assert_equal(actual_val, expected_val, f"字段 {name} 更新后值一致")


def test_entity_delete(app_code, schema_name, entity_id):
    """删除实体"""
    print(f"\n--- test_entity_delete (app={app_code}, schema={schema_name}, id={entity_id}) ---")
    resp = post("/entity/delete", {"appCode": app_code, "schemaName": schema_name, "id": entity_id})
    assert_success(resp, "删除 Entity")

    # 查询列表验证删除
    resp = get("/entity/list", {"appCode": app_code, "schemaName": schema_name, "page": 1, "pageSize": 20})
    assert_success(resp, "删除后查询 Entity 列表")

    data = resp.get("data", {})
    ids = [e.get("id") for e in data.get("list", [])]
    assert_true(entity_id not in ids, f"Entity {entity_id} 已不在列表中")


# ========== 入口 ==========

if __name__ == "__main__":
    print(f"Forma API 接口测试")
    print(f"Base URL: {BASE_URL}")
    print("=" * 50)

    # 生成唯一标识
    app_code = APP_CODE
    schema_name = gen_schema_name()

    # ===== App 测试 =====
    test_app_create(app_code)
    test_app_create_duplicate(app_code)
    test_app_detail(app_code)
    test_app_list(app_code)
    test_app_update(app_code)

    # ===== Schema 测试 =====
    test_schema_create(app_code, schema_name)
    test_schema_detail(app_code, schema_name)
    test_schema_list(app_code, schema_name)
    test_schema_update(app_code, schema_name)

    # ===== App 删除（有 Schema 时应被拒绝） =====
    test_app_delete_in_use(app_code)

    # ===== Entity 测试 =====
    entity_id = test_entity_create(app_code, schema_name)
    test_entity_detail(app_code, schema_name, entity_id)
    test_entity_list(app_code, schema_name)
    test_entity_update(app_code, schema_name, entity_id)
    test_entity_delete(app_code, schema_name, entity_id)

    # ===== 清理：删除 Schema 和 App =====
    test_schema_delete(app_code, schema_name)
    test_app_delete(app_code)

    print_summary()
