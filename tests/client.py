"""Forma API Client — 封装所有 HTTP 调用，测试文件不直接拼 URL"""

import requests
from tests.config import BASE_URL, TIMEOUT


class SchemaAPI:
    """Schema 接口封装"""

    @staticmethod
    def create(payload: dict) -> dict:
        resp = requests.post(f"{BASE_URL}/schema/create", json=payload, timeout=TIMEOUT)
        resp.raise_for_status()
        return resp.json()

    @staticmethod
    def update(payload: dict) -> dict:
        resp = requests.post(f"{BASE_URL}/schema/update", json=payload, timeout=TIMEOUT)
        resp.raise_for_status()
        return resp.json()

    @staticmethod
    def delete(name: str) -> dict:
        resp = requests.post(f"{BASE_URL}/schema/delete", json={"name": name}, timeout=TIMEOUT)
        resp.raise_for_status()
        return resp.json()

    @staticmethod
    def detail(name: str) -> dict:
        resp = requests.get(f"{BASE_URL}/schema/detail", params={"name": name}, timeout=TIMEOUT)
        resp.raise_for_status()
        return resp.json()

    @staticmethod
    def list() -> dict:
        resp = requests.get(f"{BASE_URL}/schema/list", timeout=TIMEOUT)
        resp.raise_for_status()
        return resp.json()


class EntityAPI:
    """Entity 接口封装"""

    @staticmethod
    def create(payload: dict) -> dict:
        resp = requests.post(f"{BASE_URL}/entity/create", json=payload, timeout=TIMEOUT)
        resp.raise_for_status()
        return resp.json()

    @staticmethod
    def update(payload: dict) -> dict:
        resp = requests.post(f"{BASE_URL}/entity/update", json=payload, timeout=TIMEOUT)
        resp.raise_for_status()
        return resp.json()

    @staticmethod
    def delete(schema_name: str, entity_id: str) -> dict:
        resp = requests.post(
            f"{BASE_URL}/entity/delete",
            json={"schemaName": schema_name, "id": entity_id},
            timeout=TIMEOUT,
        )
        resp.raise_for_status()
        return resp.json()

    @staticmethod
    def detail(schema_name: str, entity_id: str) -> dict:
        resp = requests.get(
            f"{BASE_URL}/entity/detail",
            params={"schemaName": schema_name, "id": entity_id},
            timeout=TIMEOUT,
        )
        resp.raise_for_status()
        return resp.json()

    @staticmethod
    def list(schema_name: str, page: int = 1, page_size: int = 20) -> dict:
        resp = requests.get(
            f"{BASE_URL}/entity/list",
            params={"schemaName": schema_name, "page": page, "pageSize": page_size},
            timeout=TIMEOUT,
        )
        resp.raise_for_status()
        return resp.json()
