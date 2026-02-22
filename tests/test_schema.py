"""Schema 接口测试 — 覆盖正向、异常、边界场景"""

import unittest
import uuid

from tests.client import SchemaAPI
from tests.config import CODE_SUCCESS, CODE_INVALID_PARAM, CODE_NOT_FOUND


def _unique_name(prefix: str = "test") -> str:
    """生成唯一 Schema 名称，避免用例间冲突"""
    return f"{prefix}_{uuid.uuid4().hex[:8]}"


def _make_schema(name: str, fields: list | None = None) -> dict:
    """构造创建 Schema 请求体"""
    if fields is None:
        fields = [
            {"name": "title", "type": "string", "required": True, "maxLength": 200, "minLength": 1},
            {"name": "count", "type": "number", "required": False},
        ]
    return {"name": name, "fields": fields}


def _safe_delete_schema(name: str):
    """安全删除 Schema，忽略错误"""
    try:
        SchemaAPI.delete(name)
    except Exception:
        pass


class TestSchemaCreate(unittest.TestCase):
    """Schema 创建接口测试"""

    def setUp(self):
        self.schema_name = _unique_name()

    def tearDown(self):
        _safe_delete_schema(self.schema_name)

    def test_create_success(self):
        """正常创建"""
        body = SchemaAPI.create(_make_schema(self.schema_name))
        self.assertEqual(body["code"], CODE_SUCCESS)

    def test_create_with_all_field_types(self):
        """覆盖所有字段类型"""
        fields = [
            {"name": "f_string", "type": "string", "required": True},
            {"name": "f_number", "type": "number", "required": False},
            {"name": "f_boolean", "type": "boolean", "required": False},
            {"name": "f_date", "type": "date", "required": False},
            {"name": "f_text", "type": "text", "required": False},
            {"name": "f_enum", "type": "enum", "required": False, "enumValues": ["a", "b", "c"]},
            {"name": "f_json", "type": "json", "required": False},
            {"name": "f_array", "type": "array", "required": False},
        ]
        body = SchemaAPI.create(_make_schema(self.schema_name, fields))
        self.assertEqual(body["code"], CODE_SUCCESS)

    def test_create_empty_name(self):
        """name 为空"""
        body = SchemaAPI.create(_make_schema(""))
        self.assertNotEqual(body["code"], CODE_SUCCESS)

    def test_create_empty_fields(self):
        """fields 为空列表"""
        body = SchemaAPI.create(_make_schema(self.schema_name, fields=[]))
        self.assertEqual(body["code"], CODE_INVALID_PARAM)

    def test_create_invalid_field_type(self):
        """字段 type 非法"""
        fields = [{"name": "bad", "type": "invalid_type", "required": True}]
        body = SchemaAPI.create(_make_schema(self.schema_name, fields))
        self.assertEqual(body["code"], CODE_INVALID_PARAM)

    def test_create_enum_without_values(self):
        """enum 类型不提供 enumValues"""
        fields = [{"name": "status", "type": "enum", "required": True}]
        body = SchemaAPI.create(_make_schema(self.schema_name, fields))
        self.assertEqual(body["code"], CODE_INVALID_PARAM)

    def test_create_duplicate_name(self):
        """重复 name 创建"""
        SchemaAPI.create(_make_schema(self.schema_name))
        body = SchemaAPI.create(_make_schema(self.schema_name))
        # 重复创建应返回非成功
        self.assertNotEqual(body["code"], CODE_SUCCESS)


class TestSchemaDetail(unittest.TestCase):
    """Schema 详情查询接口测试"""

    def setUp(self):
        self.schema_name = _unique_name()
        self.fields = [
            {"name": "title", "type": "string", "required": True, "maxLength": 200, "minLength": 1},
            {"name": "done", "type": "boolean", "required": False},
        ]
        SchemaAPI.create(_make_schema(self.schema_name, self.fields))

    def tearDown(self):
        _safe_delete_schema(self.schema_name)

    def test_detail_success(self):
        """正常查询详情"""
        body = SchemaAPI.detail(self.schema_name)
        self.assertEqual(body["code"], CODE_SUCCESS)
        data = body["data"]
        self.assertEqual(data["name"], self.schema_name)
        self.assertEqual(len(data["fields"]), 2)
        self.assertIn("createdAt", data)
        self.assertIn("updatedAt", data)

    def test_detail_fields_match(self):
        """详情字段与创建时一致"""
        body = SchemaAPI.detail(self.schema_name)
        data = body["data"]
        field_names = {f["name"] for f in data["fields"]}
        self.assertEqual(field_names, {"title", "done"})

        title_field = next(f for f in data["fields"] if f["name"] == "title")
        self.assertEqual(title_field["type"], "string")
        self.assertTrue(title_field["required"])
        self.assertEqual(title_field["maxLength"], 200)
        self.assertEqual(title_field["minLength"], 1)

    def test_detail_not_found(self):
        """查询不存在的 Schema"""
        body = SchemaAPI.detail("nonexistent_schema_xyz")
        self.assertEqual(body["code"], CODE_NOT_FOUND)


class TestSchemaList(unittest.TestCase):
    """Schema 列表查询接口测试"""

    def setUp(self):
        self.created_names = []
        for _ in range(3):
            name = _unique_name()
            SchemaAPI.create(_make_schema(name))
            self.created_names.append(name)

    def tearDown(self):
        for name in self.created_names:
            _safe_delete_schema(name)

    def test_list_success(self):
        """列表查询包含已创建的 Schema"""
        body = SchemaAPI.list()
        self.assertEqual(body["code"], CODE_SUCCESS)
        data = body["data"]
        self.assertIn("total", data)
        self.assertIn("list", data)
        self.assertGreaterEqual(data["total"], 3)
        listed_names = {s["name"] for s in data["list"]}
        for name in self.created_names:
            self.assertIn(name, listed_names)


class TestSchemaUpdate(unittest.TestCase):
    """Schema 更新接口测试"""

    def setUp(self):
        self.schema_name = _unique_name()
        self.fields = [
            {"name": "title", "type": "string", "required": True, "maxLength": 200, "minLength": 1},
            {"name": "status", "type": "enum", "required": False, "enumValues": ["open", "closed"]},
        ]
        SchemaAPI.create({
            "name": self.schema_name,
            "description": "original desc",
            "fields": self.fields,
        })

    def tearDown(self):
        _safe_delete_schema(self.schema_name)

    def test_update_mutable_fields(self):
        """更新可变属性：required, maxLength, minLength, enumValues, description"""
        updated_fields = [
            {
                "name": "title", "type": "string", "required": False,
                "maxLength": 500, "minLength": 0, "description": "updated title desc",
            },
            {
                "name": "status", "type": "enum", "required": True,
                "enumValues": ["open", "closed", "pending"],
            },
        ]
        body = SchemaAPI.update({
            "name": self.schema_name,
            "description": "updated desc",
            "fields": updated_fields,
        })
        self.assertEqual(body["code"], CODE_SUCCESS)

        # 验证更新生效
        detail = SchemaAPI.detail(self.schema_name)["data"]
        self.assertEqual(detail["description"], "updated desc")

        title_f = next(f for f in detail["fields"] if f["name"] == "title")
        self.assertFalse(title_f["required"])
        self.assertEqual(title_f["maxLength"], 500)
        self.assertEqual(title_f["description"], "updated title desc")

        status_f = next(f for f in detail["fields"] if f["name"] == "status")
        self.assertTrue(status_f["required"])
        self.assertIn("pending", status_f["enumValues"])

    def test_update_nonexistent_schema(self):
        """更新不存在的 Schema"""
        body = SchemaAPI.update({
            "name": "nonexistent_xyz",
            "fields": [{"name": "title", "type": "string", "required": True}],
        })
        self.assertEqual(body["code"], CODE_NOT_FOUND)

    def test_update_change_field_type(self):
        """尝试修改字段 type（不允许）"""
        body = SchemaAPI.update({
            "name": self.schema_name,
            "fields": [{"name": "title", "type": "number", "required": True}],
        })
        self.assertEqual(body["code"], CODE_INVALID_PARAM)

    def test_update_nonexistent_field(self):
        """更新时传入不存在的字段"""
        body = SchemaAPI.update({
            "name": self.schema_name,
            "fields": [{"name": "nonexistent_field", "type": "string", "required": True}],
        })
        self.assertEqual(body["code"], CODE_INVALID_PARAM)


class TestSchemaDelete(unittest.TestCase):
    """Schema 删除接口测试"""

    def test_delete_success(self):
        """正常删除"""
        name = _unique_name()
        SchemaAPI.create(_make_schema(name))
        body = SchemaAPI.delete(name)
        self.assertEqual(body["code"], CODE_SUCCESS)

        # 删除后查询应返回 not found
        detail = SchemaAPI.detail(name)
        self.assertEqual(detail["code"], CODE_NOT_FOUND)

    def test_delete_not_found(self):
        """删除不存在的 Schema"""
        body = SchemaAPI.delete("nonexistent_schema_xyz")
        self.assertEqual(body["code"], CODE_NOT_FOUND)


class TestSchemaCRUDFlow(unittest.TestCase):
    """Schema 完整 CRUD 流程集成测试"""

    def test_full_lifecycle(self):
        """创建 → 查询 → 列表 → 更新 → 验证 → 删除 → 确认删除"""
        name = _unique_name("lifecycle")
        fields = [
            {"name": "content", "type": "text", "required": True, "minLength": 1},
            {"name": "priority", "type": "enum", "required": True, "enumValues": ["low", "medium", "high"]},
        ]

        # 创建
        body = SchemaAPI.create({"name": name, "description": "lifecycle test", "fields": fields})
        self.assertEqual(body["code"], CODE_SUCCESS)

        # 详情
        detail = SchemaAPI.detail(name)
        self.assertEqual(detail["code"], CODE_SUCCESS)
        self.assertEqual(detail["data"]["name"], name)
        self.assertEqual(len(detail["data"]["fields"]), 2)

        # 列表包含
        list_body = SchemaAPI.list()
        listed_names = {s["name"] for s in list_body["data"]["list"]}
        self.assertIn(name, listed_names)

        # 更新
        update_body = SchemaAPI.update({
            "name": name,
            "description": "updated lifecycle",
            "fields": [
                {"name": "content", "type": "text", "required": False, "minLength": 0},
                {"name": "priority", "type": "enum", "required": True, "enumValues": ["low", "medium", "high", "critical"]},
            ],
        })
        self.assertEqual(update_body["code"], CODE_SUCCESS)

        # 验证更新
        detail2 = SchemaAPI.detail(name)
        self.assertEqual(detail2["data"]["description"], "updated lifecycle")
        content_f = next(f for f in detail2["data"]["fields"] if f["name"] == "content")
        self.assertFalse(content_f["required"])

        # 删除
        del_body = SchemaAPI.delete(name)
        self.assertEqual(del_body["code"], CODE_SUCCESS)

        # 确认删除
        detail3 = SchemaAPI.detail(name)
        self.assertEqual(detail3["code"], CODE_NOT_FOUND)


if __name__ == "__main__":
    unittest.main()
