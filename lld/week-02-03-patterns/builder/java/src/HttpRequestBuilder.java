import java.util.LinkedHashMap;
import java.util.Map;

public class HttpRequestBuilder {
    private String method = "GET";
    private String url;
    private final Map<String, String> headers = new LinkedHashMap<>();
    private String body = "";

    public HttpRequestBuilder method(String method) {
        this.method = method;
        return this;
    }

    public HttpRequestBuilder url(String url) {
        this.url = url;
        return this;
    }

    public HttpRequestBuilder header(String key, String value) {
        this.headers.put(key, value);
        return this;
    }

    public HttpRequestBuilder body(String body) {
        this.body = body;
        return this;
    }

    public static class MissingUrlException extends RuntimeException {
        public MissingUrlException(String message) {
            super(message);
        }
    }

    // Build snapshots the current headers into a new map so a later reuse of
    // this builder (more chained calls, another build()) never mutates a
    // request that was already handed out.
    public HttpRequest build() {
        if (url == null || url.isEmpty()) {
            throw new MissingUrlException("URL is required");
        }
        return new HttpRequest(method, url, new LinkedHashMap<>(headers), body);
    }
}
