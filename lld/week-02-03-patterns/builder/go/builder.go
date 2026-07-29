// Package builder implements the Builder pattern via a fluent
// HttpRequestBuilder that assembles an immutable HttpRequest step by step.
package builder

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type HttpRequest struct {
	method  string
	url     string
	headers map[string]string
	body    string
}

func (r *HttpRequest) Method() string { return r.method }
func (r *HttpRequest) URL() string    { return r.url }
func (r *HttpRequest) Body() string   { return r.body }

// Header returns the value for the given (case-sensitive) header key.
func (r *HttpRequest) Header(key string) (string, bool) {
	v, ok := r.headers[key]
	return v, ok
}

// String renders the request in a raw-HTTP-like form, with headers sorted
// for deterministic output.
func (r *HttpRequest) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s %s\n", r.method, r.url)

	keys := make([]string, 0, len(r.headers))
	for k := range r.headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s: %s\n", k, r.headers[k])
	}
	if r.body != "" {
		fmt.Fprintf(&sb, "\n%s", r.body)
	}
	return sb.String()
}

var ErrMissingURL = errors.New("builder: URL is required")

// HttpRequestBuilder accumulates request fields via chained calls and
// produces an immutable HttpRequest on Build. It is not safe for concurrent
// use by multiple goroutines.
type HttpRequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    string
}

func NewHttpRequestBuilder() *HttpRequestBuilder {
	return &HttpRequestBuilder{method: "GET", headers: make(map[string]string)}
}

func (b *HttpRequestBuilder) Method(method string) *HttpRequestBuilder {
	b.method = method
	return b
}

func (b *HttpRequestBuilder) URL(url string) *HttpRequestBuilder {
	b.url = url
	return b
}

func (b *HttpRequestBuilder) Header(key, value string) *HttpRequestBuilder {
	b.headers[key] = value
	return b
}

func (b *HttpRequestBuilder) Body(body string) *HttpRequestBuilder {
	b.body = body
	return b
}

// Build validates required fields and returns a new HttpRequest with its own
// copy of the headers map, so later mutation of the builder never leaks
// into a previously built request.
func (b *HttpRequestBuilder) Build() (*HttpRequest, error) {
	if b.url == "" {
		return nil, ErrMissingURL
	}
	headers := make(map[string]string, len(b.headers))
	for k, v := range b.headers {
		headers[k] = v
	}
	return &HttpRequest{method: b.method, url: b.url, headers: headers, body: b.body}, nil
}
