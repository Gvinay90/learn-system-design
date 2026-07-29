import pytest

from adapter import (
    LegacyXmlDataProvider,
    ModernDataProvider,
    RecordNotFoundError,
    XmlToJsonAdapter,
    fetch_and_sum,
)


def test_adapter_translates_legacy_call():
    legacy = LegacyXmlDataProvider({"u1": "42"})
    adapted = XmlToJsonAdapter(legacy)
    record = adapted.fetch_json("u1")
    assert record == {"id": "u1", "value": "42"}


def test_adapter_propagates_not_found():
    legacy = LegacyXmlDataProvider({})
    adapted = XmlToJsonAdapter(legacy)
    with pytest.raises(RecordNotFoundError):
        adapted.fetch_json("missing")


def test_client_code_is_provider_agnostic():
    legacy = XmlToJsonAdapter(LegacyXmlDataProvider({"a": "10", "b": "20"}))
    modern = ModernDataProvider({"a": "10", "b": "20"})

    assert fetch_and_sum(legacy, ["a", "b"]) == 30
    assert fetch_and_sum(modern, ["a", "b"]) == 30


def test_legacy_provider_raw_xml():
    legacy = LegacyXmlDataProvider({"u1": "42"})
    raw = legacy.fetch_xml("u1")
    assert "u1" in raw and "42" in raw
