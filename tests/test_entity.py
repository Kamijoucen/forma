"""Entity 接口测试 — 覆盖正向、异常、边界场景"""

import unittest
import uuid

from tests.client import SchemaAPI, EntityAPI
from tests.config import CODE_SUCCESS, CODE_INVALID_PARAM, CODE_NOT_FOUND


def _unique_name(prefix: str = "ent_test") -> str:
    return f"{prefix}_{uuid.uuid4().hex[:8]}"


def _safe_delete_entity(schema_name: str, entity_id: str):
    try:
        EntityAPI.delete(schema_name, entity_id)
    except Exception:
        pass


def _safe_delete_schema(name: str):
    try:
        SchemaAPI.delete(name)
    except Exception:
        pass


class TestEntityCreate(unittest.TestCase):
    """Entity 创建接口测试"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("ent_create")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "title", "type": "string", "required": True, "maxLength": 200, "minLength": 1},
                {"name": "count", "type": "number", "required": False},
                {"name": "done", "type": "boolean", "required": True},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def setUp(self):
        self.created_ids = []

    def tearDown(self):
        for eid in self.created_ids:
            _safe_delete_entity(self.schema_name, eid)

    def _create_entity(self, fields: list) -> dict:
        body = EntityAPI.create({"schemaName": self.schema_name, "fields": fields})
        if body["code"] == CODE_SUCCESS and body.get("data"):
            self.created_ids.append(body["data"]["id"])
        return body

    def test_create_success(self):
        """正常创建实体"""
        body = self._create_entity([
            {"name": "title", "type": "string", "value": "Test Item"},
            {"name": "done", "type": "boolean", "value": "false"},
        ])
        self.assertEqual(body["code"], CODE_SUCCESS)
        self.assertIn("id", body["data"])

    def test_create_with_optional_field(self):
        """创建时提供可选字段"""
        body = self._create_entity([
            {"name": "title", "type": "string", "value": "With Count"},
            {"name": "count", "type": "number", "value": "42"},
            {"name": "done", "type": "boolean", "value": "true"},
        ])
        self.assertEqual(body["code"], CODE_SUCCESS)

    def test_create_missing_required_field(self):
        """缺少必填字段 title"""
        body = self._create_entity([
            {"name": "done", "type": "boolean", "value": "false"},
        ])
        self.assertEqual(body["code"], CODE_INVALID_PARAM)

    def test_create_nonexistent_schema(self):
        """引用不存在的 schemaName"""
        body = EntityAPI.create({
            "schemaName": "nonexistent_schema_xyz",
            "fields": [{"name": "title", "type": "string", "value": "x"}],
        })
        self.assertEqual(body["code"], CODE_NOT_FOUND)

    def test_create_undefined_field(self):
        """传入 Schema 中未定义的字段"""
        body = self._create_entity([
            {"name": "title", "type": "string", "value": "Test"},
            {"name": "done", "type": "boolean", "value": "false"},
            {"name": "unknown_field", "type": "string", "value": "x"},
        ])
        self.assertEqual(body["code"], CODE_INVALID_PARAM)

    def test_create_type_mismatch(self):
        """字段 type 与 Schema 定义不匹配"""
        body = self._create_entity([
            {"name": "title", "type": "number", "value": "123"},
            {"name": "done", "type": "boolean", "value": "false"},
        ])
        self.assertEqual(body["code"], CODE_INVALID_PARAM)

    def test_create_value_validation_fail(self):
        """值校验失败（boolean 值非法）"""
        body = self._create_entity([
            {"name": "title", "type": "string", "value": "Test"},
            {"name": "done", "type": "boolean", "value": "yes"},
        ])
        self.assertEqual(body["code"], CODE_INVALID_PARAM)


class TestEntityDetail(unittest.TestCase):
    """Entity 详情查询接口测试"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("ent_detail")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "title", "type": "string", "required": True},
                {"name": "score", "type": "number", "required": False},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def setUp(self):
        self.created_ids = []

    def tearDown(self):
        for eid in self.created_ids:
            _safe_delete_entity(self.schema_name, eid)

    def test_detail_success(self):
        """正常查询实体详情"""
        create_body = EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [
                {"name": "title", "type": "string", "value": "Detail Test"},
                {"name": "score", "type": "number", "value": "99.5"},
            ],
        })
        eid = create_body["data"]["id"]
        self.created_ids.append(eid)

        body = EntityAPI.detail(self.schema_name, eid)
        self.assertEqual(body["code"], CODE_SUCCESS)
        data = body["data"]
        self.assertEqual(data["id"], eid)
        self.assertEqual(data["schemaName"], self.schema_name)
        self.assertIn("createdAt", data)
        self.assertIn("updatedAt", data)

        field_map = {f["name"]: f for f in data["fields"]}
        self.assertEqual(field_map["title"]["value"], "Detail Test")
        self.assertEqual(field_map["score"]["value"], "99.5")

    def test_detail_not_found(self):
        """查询不存在的实体"""
        body = EntityAPI.detail(self.schema_name, "999999999")
        self.assertEqual(body["code"], CODE_NOT_FOUND)

    def test_detail_invalid_id(self):
        """ID 格式非法"""
        body = EntityAPI.detail(self.schema_name, "not_a_number")
        self.assertEqual(body["code"], CODE_INVALID_PARAM)


class TestEntityList(unittest.TestCase):
    """Entity 列表查询接口测试"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("ent_list")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "title", "type": "string", "required": True},
            ],
        })
        cls.entity_ids = []
        for i in range(5):
            body = EntityAPI.create({
                "schemaName": cls.schema_name,
                "fields": [{"name": "title", "type": "string", "value": f"Item {i}"}],
            })
            cls.entity_ids.append(body["data"]["id"])

    @classmethod
    def tearDownClass(cls):
        for eid in cls.entity_ids:
            _safe_delete_entity(cls.schema_name, eid)
        _safe_delete_schema(cls.schema_name)

    def test_list_success(self):
        """列表查询返回正确数量"""
        body = EntityAPI.list(self.schema_name)
        self.assertEqual(body["code"], CODE_SUCCESS)
        self.assertEqual(body["data"]["total"], 5)
        self.assertEqual(len(body["data"]["list"]), 5)

    def test_list_pagination(self):
        """分页查询"""
        body = EntityAPI.list(self.schema_name, page=1, page_size=2)
        self.assertEqual(body["code"], CODE_SUCCESS)
        self.assertEqual(body["data"]["total"], 5)
        self.assertEqual(len(body["data"]["list"]), 2)

        body2 = EntityAPI.list(self.schema_name, page=2, page_size=2)
        self.assertEqual(body2["code"], CODE_SUCCESS)
        self.assertEqual(len(body2["data"]["list"]), 2)

        body3 = EntityAPI.list(self.schema_name, page=3, page_size=2)
        self.assertEqual(body3["code"], CODE_SUCCESS)
        self.assertEqual(len(body3["data"]["list"]), 1)

    def test_list_empty(self):
        """Schema 下无实体时返回空列表"""
        empty_schema = _unique_name("ent_empty")
        SchemaAPI.create({
            "name": empty_schema,
            "fields": [{"name": "x", "type": "string", "required": True}],
        })
        try:
            body = EntityAPI.list(empty_schema)
            self.assertEqual(body["code"], CODE_SUCCESS)
            self.assertEqual(body["data"]["total"], 0)
            self.assertEqual(len(body["data"]["list"]), 0)
        finally:
            _safe_delete_schema(empty_schema)


class TestEntityUpdate(unittest.TestCase):
    """Entity 更新接口测试"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("ent_update")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "title", "type": "string", "required": True},
                {"name": "done", "type": "boolean", "required": True},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def setUp(self):
        self.created_ids = []

    def tearDown(self):
        for eid in self.created_ids:
            _safe_delete_entity(self.schema_name, eid)

    def test_update_success(self):
        """正常更新实体"""
        create_body = EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [
                {"name": "title", "type": "string", "value": "Original"},
                {"name": "done", "type": "boolean", "value": "false"},
            ],
        })
        eid = create_body["data"]["id"]
        self.created_ids.append(eid)

        update_body = EntityAPI.update({
            "schemaName": self.schema_name,
            "id": eid,
            "fields": [
                {"name": "title", "type": "string", "value": "Updated"},
                {"name": "done", "type": "boolean", "value": "true"},
            ],
        })
        self.assertEqual(update_body["code"], CODE_SUCCESS)

        # 验证更新生效
        detail = EntityAPI.detail(self.schema_name, eid)
        field_map = {f["name"]: f for f in detail["data"]["fields"]}
        self.assertEqual(field_map["title"]["value"], "Updated")
        self.assertEqual(field_map["done"]["value"], "true")

    def test_update_not_found(self):
        """更新不存在的实体"""
        body = EntityAPI.update({
            "schemaName": self.schema_name,
            "id": "999999999",
            "fields": [
                {"name": "title", "type": "string", "value": "x"},
                {"name": "done", "type": "boolean", "value": "false"},
            ],
        })
        self.assertEqual(body["code"], CODE_NOT_FOUND)

    def test_update_invalid_id(self):
        """ID 格式非法"""
        body = EntityAPI.update({
            "schemaName": self.schema_name,
            "id": "abc",
            "fields": [
                {"name": "title", "type": "string", "value": "x"},
                {"name": "done", "type": "boolean", "value": "false"},
            ],
        })
        self.assertEqual(body["code"], CODE_INVALID_PARAM)

    def test_update_validation_fail(self):
        """更新时值校验失败"""
        create_body = EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [
                {"name": "title", "type": "string", "value": "Test"},
                {"name": "done", "type": "boolean", "value": "false"},
            ],
        })
        eid = create_body["data"]["id"]
        self.created_ids.append(eid)

        body = EntityAPI.update({
            "schemaName": self.schema_name,
            "id": eid,
            "fields": [
                {"name": "title", "type": "string", "value": "Updated"},
                {"name": "done", "type": "boolean", "value": "not_bool"},
            ],
        })
        self.assertEqual(body["code"], CODE_INVALID_PARAM)


class TestEntityDelete(unittest.TestCase):
    """Entity 删除接口测试"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("ent_delete")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "title", "type": "string", "required": True},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def test_delete_success(self):
        """正常删除实体"""
        create_body = EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "title", "type": "string", "value": "To Delete"}],
        })
        eid = create_body["data"]["id"]

        body = EntityAPI.delete(self.schema_name, eid)
        self.assertEqual(body["code"], CODE_SUCCESS)

        # 删除后查询应返回 not found
        detail = EntityAPI.detail(self.schema_name, eid)
        self.assertEqual(detail["code"], CODE_NOT_FOUND)

    def test_delete_not_found(self):
        """删除不存在的实体"""
        body = EntityAPI.delete(self.schema_name, "999999999")
        self.assertEqual(body["code"], CODE_NOT_FOUND)

    def test_delete_invalid_id(self):
        """ID 格式非法"""
        body = EntityAPI.delete(self.schema_name, "not_number")
        self.assertEqual(body["code"], CODE_INVALID_PARAM)


class TestEntityCRUDFlow(unittest.TestCase):
    """Entity 完整 CRUD 流程集成测试"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("ent_flow")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "title", "type": "string", "required": True, "maxLength": 100},
                {"name": "tags", "type": "array", "required": False},
                {"name": "meta", "type": "json", "required": False},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def test_full_lifecycle(self):
        """创建 → 查询 → 列表 → 更新 → 验证 → 删除 → 确认删除"""
        # 创建
        create_body = EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [
                {"name": "title", "type": "string", "value": "Lifecycle"},
                {"name": "tags", "type": "array", "value": '["a","b"]'},
                {"name": "meta", "type": "json", "value": '{"key":"val"}'},
            ],
        })
        self.assertEqual(create_body["code"], CODE_SUCCESS)
        eid = create_body["data"]["id"]

        # 详情
        detail = EntityAPI.detail(self.schema_name, eid)
        self.assertEqual(detail["code"], CODE_SUCCESS)
        self.assertEqual(detail["data"]["id"], eid)

        # 列表
        list_body = EntityAPI.list(self.schema_name)
        self.assertEqual(list_body["code"], CODE_SUCCESS)
        self.assertGreaterEqual(list_body["data"]["total"], 1)

        # 更新
        update_body = EntityAPI.update({
            "schemaName": self.schema_name,
            "id": eid,
            "fields": [
                {"name": "title", "type": "string", "value": "Updated Lifecycle"},
                {"name": "tags", "type": "array", "value": '["c"]'},
                {"name": "meta", "type": "json", "value": '{"key":"new"}'},
            ],
        })
        self.assertEqual(update_body["code"], CODE_SUCCESS)

        # 验证更新
        detail2 = EntityAPI.detail(self.schema_name, eid)
        field_map = {f["name"]: f for f in detail2["data"]["fields"]}
        self.assertEqual(field_map["title"]["value"], "Updated Lifecycle")
        self.assertEqual(field_map["tags"]["value"], '["c"]')

        # 删除
        del_body = EntityAPI.delete(self.schema_name, eid)
        self.assertEqual(del_body["code"], CODE_SUCCESS)

        # 确认删除
        detail3 = EntityAPI.detail(self.schema_name, eid)
        self.assertEqual(detail3["code"], CODE_NOT_FOUND)


if __name__ == "__main__":
    unittest.main()
