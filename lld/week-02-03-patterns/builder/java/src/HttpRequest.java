import java.util.Collections;
import java.util.Map;
import java.util.TreeMap;

public final class HttpRequest {
    private final String method;
    private final String url;
    private final Map<String, String> headers;
    private final String body;

    HttpRequest(String method, String url, Map<String, String> headers, String body) {
        this.method = method;
        this.url = url;
        this.headers = Collections.unmodifiableMap(new TreeMap<>(headers));
        this.body = body;
    }

    public String getMethod() {
        return method;
    }

    public String getUrl() {
        return url;
    }

    public String getHeader(String key) {
        return headers.get(key);
    }

    public Map<String, String> getHeaders() {
        return headers;
    }

    public String getBody() {
        return body;
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        sb.append(method).append(' ').append(url).append('\n');
        for (Map.Entry<String, String> e : headers.entrySet()) {
            sb.append(e.getKey()).append(": ").append(e.getValue()).append('\n');
        }
        if (body != null && !body.isEmpty()) {
            sb.append('\n').append(body);
        }
        return sb.toString();
    }
}
