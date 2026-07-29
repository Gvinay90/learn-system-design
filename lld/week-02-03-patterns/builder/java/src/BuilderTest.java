/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out BuilderTest` directly.
 */
public class BuilderTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testFluentChainBuildsExpectedRequest();
        testDefaultMethodIsGet();
        testBuildFailsWithoutUrl();
        testBuiltRequestIsImmutableFromBuilderReuse();
        System.out.println("All BuilderTest cases passed.");
    }

    private static void testFluentChainBuildsExpectedRequest() {
        HttpRequest req = new HttpRequestBuilder()
                .method("POST")
                .url("https://api.example.com/orders")
                .header("Content-Type", "application/json")
                .body("{\"item\":\"widget\"}")
                .build();

        assertEquals("POST", req.getMethod(), "method");
        assertEquals("https://api.example.com/orders", req.getUrl(), "url");
        assertEquals("application/json", req.getHeader("Content-Type"), "content-type header");
        assertEquals("{\"item\":\"widget\"}", req.getBody(), "body");
    }

    private static void testDefaultMethodIsGet() {
        HttpRequest req = new HttpRequestBuilder().url("https://example.com").build();
        assertEquals("GET", req.getMethod(), "default method");
    }

    private static void testBuildFailsWithoutUrl() {
        try {
            new HttpRequestBuilder().method("GET").build();
            throw new AssertionError("expected MissingUrlException");
        } catch (HttpRequestBuilder.MissingUrlException e) {
            // expected
        }
    }

    private static void testBuiltRequestIsImmutableFromBuilderReuse() {
        HttpRequestBuilder b = new HttpRequestBuilder().url("https://example.com").header("X-A", "1");
        HttpRequest first = b.build();

        b.header("X-A", "2").header("X-B", "new");
        HttpRequest second = b.build();

        assertEquals("1", first.getHeader("X-A"), "first request X-A unchanged");
        assertEquals(null, first.getHeader("X-B"), "first request should not see later header");
        assertEquals("2", second.getHeader("X-A"), "second request X-A updated");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        boolean equal = (expected == null) ? (actual == null) : expected.equals(actual);
        if (!equal) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
