"""Builder pattern — fluent HttpRequestBuilder producing an immutable
HttpRequest. See ../README.md for the design writeup.
"""
from __future__ import annotations

from dataclasses import dataclass, field
from typing import Dict, Optional


class MissingUrlError(Exception):
    pass


@dataclass(frozen=True)
class HttpRequest:
    method: str
    url: str
    headers: Dict[str, str]
    body: str

    def header(self, key: str) -> Optional[str]:
        return self.headers.get(key)

    def __str__(self) -> str:
        lines = [f"{self.method} {self.url}"]
        for key in sorted(self.headers):
            lines.append(f"{key}: {self.headers[key]}")
        if self.body:
            lines.append("")
            lines.append(self.body)
        return "\n".join(lines)


class HttpRequestBuilder:
    def __init__(self) -> None:
        self._method = "GET"
        self._url: Optional[str] = None
        self._headers: Dict[str, str] = {}
        self._body = ""

    def method(self, method: str) -> "HttpRequestBuilder":
        self._method = method
        return self

    def url(self, url: str) -> "HttpRequestBuilder":
        self._url = url
        return self

    def header(self, key: str, value: str) -> "HttpRequestBuilder":
        self._headers[key] = value
        return self

    def body(self, body: str) -> "HttpRequestBuilder":
        self._body = body
        return self

    def build(self) -> HttpRequest:
        if not self._url:
            raise MissingUrlError("URL is required")
        # Snapshot the headers dict so a later reuse of this builder never
        # mutates a request that was already handed out.
        return HttpRequest(self._method, self._url, dict(self._headers), self._body)


def _demo() -> None:
    request = (
        HttpRequestBuilder()
        .method("POST")
        .url("https://api.example.com/orders")
        .header("Content-Type", "application/json")
        .header("Authorization", "Bearer token123")
        .body('{"item":"widget"}')
        .build()
    )
    print(request)


if __name__ == "__main__":
    _demo()
