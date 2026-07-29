"""Adapter pattern — legacy XML data provider adapted to a modern
JSON-style DataProvider interface. See ../README.md for the design writeup."""
from __future__ import annotations

import re
from typing import Dict, List, Protocol


class RecordNotFoundError(Exception):
    pass


class DataProvider(Protocol):
    def fetch_json(self, record_id: str) -> Dict[str, str]: ...


class LegacyXmlDataProvider:
    def __init__(self, store: Dict[str, str]):
        self._store = store

    def fetch_xml(self, record_id: str) -> str:
        value = self._store.get(record_id)
        if value is None:
            raise RecordNotFoundError(record_id)
        return f'<record id="{record_id}">{value}</record>'


_RECORD_RE = re.compile(r'<record id="(.*?)">(.*?)</record>')


class XmlToJsonAdapter:
    def __init__(self, legacy: LegacyXmlDataProvider):
        self._legacy = legacy

    def fetch_json(self, record_id: str) -> Dict[str, str]:
        raw = self._legacy.fetch_xml(record_id)
        match = _RECORD_RE.match(raw)
        if not match:
            raise RecordNotFoundError(record_id)
        return {"id": match.group(1), "value": match.group(2)}


class ModernDataProvider:
    def __init__(self, store: Dict[str, str]):
        self._store = store

    def fetch_json(self, record_id: str) -> Dict[str, str]:
        value = self._store.get(record_id)
        if value is None:
            raise RecordNotFoundError(record_id)
        return {"id": record_id, "value": value}


def fetch_and_sum(provider: DataProvider, ids: List[str]) -> int:
    return sum(int(provider.fetch_json(i)["value"]) for i in ids)


def _demo() -> None:
    legacy = LegacyXmlDataProvider({"u1": "42"})
    adapted = XmlToJsonAdapter(legacy)
    print("Adapted record:", adapted.fetch_json("u1"))
    print("Sum via adapter:", fetch_and_sum(adapted, ["u1"]))


if __name__ == "__main__":
    _demo()
