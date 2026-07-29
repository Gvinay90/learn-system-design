import pytest

from builder import HttpRequestBuilder, MissingUrlError


def test_fluent_chain_builds_expected_request():
    req = (
        HttpRequestBuilder()
        .method("POST")
        .url("https://api.example.com/orders")
        .header("Content-Type", "application/json")
        .header("Authorization", "Bearer token123")
        .body('{"item":"widget"}')
        .build()
    )

    assert req.method == "POST"
    assert req.url == "https://api.example.com/orders"
    assert req.header("Content-Type") == "application/json"
    assert req.body == '{"item":"widget"}'


def test_default_method_is_get():
    req = HttpRequestBuilder().url("https://example.com").build()
    assert req.method == "GET"


def test_build_fails_without_url():
    with pytest.raises(MissingUrlError):
        HttpRequestBuilder().method("GET").build()


def test_built_request_is_immutable_from_builder_reuse():
    b = HttpRequestBuilder().url("https://example.com").header("X-A", "1")
    first = b.build()

    b.header("X-A", "2").header("X-B", "new")
    second = b.build()

    assert first.header("X-A") == "1"
    assert first.header("X-B") is None
    assert second.header("X-A") == "2"


def test_built_request_is_frozen():
    req = HttpRequestBuilder().url("https://example.com").build()
    with pytest.raises(Exception):
        req.method = "POST"
