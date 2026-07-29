package builder

import (
	"strings"
	"testing"
)

func TestFluentChainBuildsExpectedRequest(t *testing.T) {
	req, err := NewHttpRequestBuilder().
		Method("POST").
		URL("https://api.example.com/orders").
		Header("Content-Type", "application/json").
		Header("Authorization", "Bearer token123").
		Body(`{"item":"widget"}`).
		Build()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if req.Method() != "POST" {
		t.Fatalf("expected POST, got %s", req.Method())
	}
	if req.URL() != "https://api.example.com/orders" {
		t.Fatalf("unexpected URL: %s", req.URL())
	}
	if v, ok := req.Header("Content-Type"); !ok || v != "application/json" {
		t.Fatalf("unexpected Content-Type header: %s, %v", v, ok)
	}
	if req.Body() != `{"item":"widget"}` {
		t.Fatalf("unexpected body: %s", req.Body())
	}
}

func TestDefaultMethodIsGet(t *testing.T) {
	req, err := NewHttpRequestBuilder().URL("https://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if req.Method() != "GET" {
		t.Fatalf("expected default GET, got %s", req.Method())
	}
}

func TestBuildFailsWithoutURL(t *testing.T) {
	_, err := NewHttpRequestBuilder().Method("GET").Build()
	if err != ErrMissingURL {
		t.Fatalf("expected ErrMissingURL, got %v", err)
	}
}

func TestBuiltRequestIsImmutableFromBuilderReuse(t *testing.T) {
	b := NewHttpRequestBuilder().URL("https://example.com").Header("X-A", "1")
	first, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	b.Header("X-A", "2").Header("X-B", "new")
	second, err := b.Build()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if v, _ := first.Header("X-A"); v != "1" {
		t.Fatalf("expected first request's X-A to stay 1, got %s", v)
	}
	if _, ok := first.Header("X-B"); ok {
		t.Fatalf("first request should not see header added after it was built")
	}
	if v, _ := second.Header("X-A"); v != "2" {
		t.Fatalf("expected second request's X-A to be 2, got %s", v)
	}
}

func TestStringRendersDeterministicHeaderOrder(t *testing.T) {
	req, _ := NewHttpRequestBuilder().
		Method("GET").
		URL("https://example.com").
		Header("Z-Header", "z").
		Header("A-Header", "a").
		Build()

	s := req.String()
	if strings.Index(s, "A-Header") > strings.Index(s, "Z-Header") {
		t.Fatalf("expected headers sorted alphabetically, got:\n%s", s)
	}
}
