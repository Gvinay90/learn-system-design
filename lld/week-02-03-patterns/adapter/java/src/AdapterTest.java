import java.util.List;
import java.util.Map;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out AdapterTest` directly.
 */
public class AdapterTest {

    public static void main(String[] args) throws Exception {
        runAll();
    }

    public static void runAll() throws Exception {
        testAdapterTranslatesLegacyCall();
        testAdapterPropagatesNotFound();
        testClientCodeIsProviderAgnostic();
        System.out.println("All AdapterTest cases passed.");
    }

    private static void testAdapterTranslatesLegacyCall() throws Exception {
        LegacyXmlDataProvider legacy = new LegacyXmlDataProvider(Map.of("u1", "42"));
        DataProvider adapted = new XmlToJsonAdapter(legacy);
        Map<String, String> record = adapted.fetchJson("u1");
        assertEquals("u1", record.get("id"), "adapted id");
        assertEquals("42", record.get("value"), "adapted value");
    }

    private static void testAdapterPropagatesNotFound() {
        LegacyXmlDataProvider legacy = new LegacyXmlDataProvider(Map.of());
        DataProvider adapted = new XmlToJsonAdapter(legacy);
        try {
            adapted.fetchJson("missing");
            throw new AssertionError("expected RecordNotFoundException");
        } catch (RecordNotFoundException e) {
            // expected
        }
    }

    private static void testClientCodeIsProviderAgnostic() throws Exception {
        DataProvider legacy = new XmlToJsonAdapter(new LegacyXmlDataProvider(Map.of("a", "10", "b", "20")));
        DataProvider modern = new ModernDataProvider(Map.of("a", "10", "b", "20"));

        int legacySum = Client.fetchAndSum(legacy, List.of("a", "b"));
        int modernSum = Client.fetchAndSum(modern, List.of("a", "b"));
        assertEquals(30, legacySum, "legacy-backed sum");
        assertEquals(30, modernSum, "modern-backed sum");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
