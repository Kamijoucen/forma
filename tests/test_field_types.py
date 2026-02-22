"""8 种字段类型的值校验测试 — 使用 subTest 数据驱动"""

import unittest
import uuid

from tests.client import SchemaAPI, EntityAPI
from tests.config import CODE_SUCCESS, CODE_INVALID_PARAM


def _unique_name(prefix: str = "ftype") -> str:
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


class TestStringFieldType(unittest.TestCase):
    """string 类型字段值校验"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("str")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "content", "type": "string", "required": True, "maxLength": 10, "minLength": 2},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def _create(self, value: str) -> dict:
        return EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "content", "type": "string", "value": value}],
        })

    def test_valid_values(self):
        valid_cases = [
            ("normal", "hello"),
            ("min_boundary", "ab"),         # 正好 minLength=2
            ("max_boundary", "a" * 10),     # 正好 maxLength=10
            ("unicode", "你好世界"),          # unicode 4 字符，在范围内
        ]
        for label, value in valid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_SUCCESS, f"value={value!r}")
                _safe_delete_entity(self.schema_name, body["data"]["id"])

    def test_invalid_values(self):
        invalid_cases = [
            ("too_short", "a"),             # 长度 1 < minLength 2
            ("too_long", "a" * 11),         # 长度 11 > maxLength 10
        ]
        for label, value in invalid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_INVALID_PARAM, f"value={value!r}")


class TestTextFieldType(unittest.TestCase):
    """text 类型字段值校验（与 string 共享校验逻辑）"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("txt")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "body", "type": "text", "required": True, "minLength": 1, "maxLength": 500},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def test_valid_text(self):
        body = EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "body", "type": "text", "value": "A" * 500}],
        })
        self.assertEqual(body["code"], CODE_SUCCESS)
        _safe_delete_entity(self.schema_name, body["data"]["id"])

    def test_text_too_long(self):
        body = EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "body", "type": "text", "value": "A" * 501}],
        })
        self.assertEqual(body["code"], CODE_INVALID_PARAM)


class TestNumberFieldType(unittest.TestCase):
    """number 类型字段值校验"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("num")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "amount", "type": "number", "required": True},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def _create(self, value: str) -> dict:
        return EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "amount", "type": "number", "value": value}],
        })

    def test_valid_numbers(self):
        valid_cases = [
            ("integer", "42"),
            ("negative", "-10"),
            ("float", "3.14"),
            ("neg_float", "-0.5"),
            ("zero", "0"),
            ("scientific", "1e10"),
        ]
        for label, value in valid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_SUCCESS, f"value={value!r}")
                _safe_delete_entity(self.schema_name, body["data"]["id"])

    def test_invalid_numbers(self):
        invalid_cases = [
            ("text", "hello"),
            ("bool_str", "true"),
            ("empty", ""),
            ("special", "NaN_custom"),
        ]
        for label, value in invalid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_INVALID_PARAM, f"value={value!r}")


class TestBooleanFieldType(unittest.TestCase):
    """boolean 类型字段值校验"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("bool")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "flag", "type": "boolean", "required": True},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def _create(self, value: str) -> dict:
        return EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "flag", "type": "boolean", "value": value}],
        })

    def test_valid_booleans(self):
        for value in ["true", "false"]:
            with self.subTest(value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_SUCCESS)
                _safe_delete_entity(self.schema_name, body["data"]["id"])

    def test_invalid_booleans(self):
        invalid_cases = ["True", "FALSE", "yes", "no", "1", "0", ""]
        for value in invalid_cases:
            with self.subTest(value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_INVALID_PARAM, f"value={value!r}")


class TestDateFieldType(unittest.TestCase):
    """date 类型字段值校验（格式：2006-01-02 15:04:05）"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("date")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "event_time", "type": "date", "required": True},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def _create(self, value: str) -> dict:
        return EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "event_time", "type": "date", "value": value}],
        })

    def test_valid_dates(self):
        valid_cases = [
            ("normal", "2026-01-15 10:30:00"),
            ("midnight", "2026-12-31 00:00:00"),
            ("noon", "2026-06-15 12:00:00"),
        ]
        for label, value in valid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_SUCCESS, f"value={value!r}")
                _safe_delete_entity(self.schema_name, body["data"]["id"])

    def test_invalid_dates(self):
        invalid_cases = [
            ("slash_format", "2026/01/15 10:30:00"),
            ("date_only", "2026-01-15"),
            ("iso_format", "2026-01-15T10:30:00Z"),
            ("timestamp", "1706000000"),
            ("empty", ""),
            ("garbage", "not a date"),
        ]
        for label, value in invalid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_INVALID_PARAM, f"value={value!r}")


class TestEnumFieldType(unittest.TestCase):
    """enum 类型字段值校验"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("enum")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {
                    "name": "priority",
                    "type": "enum",
                    "required": True,
                    "enumValues": ["low", "medium", "high"],
                },
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def _create(self, value: str) -> dict:
        return EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "priority", "type": "enum", "value": value}],
        })

    def test_valid_enum_values(self):
        for value in ["low", "medium", "high"]:
            with self.subTest(value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_SUCCESS)
                _safe_delete_entity(self.schema_name, body["data"]["id"])

    def test_invalid_enum_values(self):
        invalid_cases = ["LOW", "critical", "", "low ", " medium"]
        for value in invalid_cases:
            with self.subTest(value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_INVALID_PARAM, f"value={value!r}")


class TestJSONFieldType(unittest.TestCase):
    """json 类型字段值校验"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("json")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "data", "type": "json", "required": True},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def _create(self, value: str) -> dict:
        return EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "data", "type": "json", "value": value}],
        })

    def test_valid_json(self):
        valid_cases = [
            ("object", '{"key":"value"}'),
            ("array", '[1,2,3]'),
            ("nested", '{"a":{"b":[1,2]}}'),
            ("string", '"hello"'),
            ("number", '42'),
            ("boolean", 'true'),
            ("null", 'null'),
        ]
        for label, value in valid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_SUCCESS, f"value={value!r}")
                _safe_delete_entity(self.schema_name, body["data"]["id"])

    def test_invalid_json(self):
        invalid_cases = [
            ("broken_obj", '{key: value}'),
            ("single_quote", "{'a': 1}"),
            ("trailing_comma", '{"a":1,}'),
            ("plain_text", 'hello world'),
        ]
        for label, value in invalid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_INVALID_PARAM, f"value={value!r}")


class TestArrayFieldType(unittest.TestCase):
    """array 类型字段值校验"""

    @classmethod
    def setUpClass(cls):
        cls.schema_name = _unique_name("arr")
        SchemaAPI.create({
            "name": cls.schema_name,
            "fields": [
                {"name": "items", "type": "array", "required": True},
            ],
        })

    @classmethod
    def tearDownClass(cls):
        _safe_delete_schema(cls.schema_name)

    def _create(self, value: str) -> dict:
        return EntityAPI.create({
            "schemaName": self.schema_name,
            "fields": [{"name": "items", "type": "array", "value": value}],
        })

    def test_valid_arrays(self):
        valid_cases = [
            ("strings", '["a","b","c"]'),
            ("numbers", '[1,2,3]'),
            ("mixed", '[1,"two",true,null]'),
            ("nested", '[[1,2],[3,4]]'),
            ("empty", '[]'),
            ("objects", '[{"a":1},{"b":2}]'),
        ]
        for label, value in valid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_SUCCESS, f"value={value!r}")
                _safe_delete_entity(self.schema_name, body["data"]["id"])

    def test_invalid_arrays(self):
        invalid_cases = [
            ("object_not_array", '{"key":"val"}'),
            ("string_not_array", '"hello"'),
            ("number_not_array", '42'),
            ("invalid_json", '[1,2,'),
            ("plain_text", 'not json'),
        ]
        for label, value in invalid_cases:
            with self.subTest(label=label, value=value):
                body = self._create(value)
                self.assertEqual(body["code"], CODE_INVALID_PARAM, f"value={value!r}")


if __name__ == "__main__":
    unittest.main()
