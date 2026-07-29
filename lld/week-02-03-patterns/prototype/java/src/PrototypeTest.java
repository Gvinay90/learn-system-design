import java.util.List;
import java.util.Map;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out PrototypeTest` directly.
 */
public class PrototypeTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testCloneProducesEqualButDistinctValues();
        testMutatingCloneTagsDoesNotAffectOriginal();
        testMutatingCloneSectionsDoesNotAffectOriginal();
        testMutatingClonePropsDoesNotAffectOriginal();
        testMutatingCloneAuthorDoesNotAffectOriginal();
        System.out.println("All PrototypeTest cases passed.");
    }

    private static Document newTestDocument() {
        return new Document(
                "Q3 Report",
                new Metadata("Vinay", List.of("finance", "draft")),
                List.of("Intro", "Numbers"),
                Map.of("status", "draft"));
    }

    private static void testCloneProducesEqualButDistinctValues() {
        Document original = newTestDocument();
        Document clone = original.deepClone();

        assertEquals(original.getTitle(), clone.getTitle(), "title");
        assertEquals(original.getMeta().getAuthor(), clone.getMeta().getAuthor(), "author");
        if (clone == original) {
            throw new AssertionError("clone should be a distinct object from original");
        }
        if (clone.getMeta() == original.getMeta()) {
            throw new AssertionError("clone.meta should be a distinct object from original.meta");
        }
    }

    private static void testMutatingCloneTagsDoesNotAffectOriginal() {
        Document original = newTestDocument();
        Document clone = original.deepClone();

        clone.getMeta().getTags().set(0, "mutated");
        clone.getMeta().getTags().add("extra");

        assertEquals("finance", original.getMeta().getTags().get(0), "original tag 0 unchanged");
        assertEquals(2, original.getMeta().getTags().size(), "original tags size unchanged");
    }

    private static void testMutatingCloneSectionsDoesNotAffectOriginal() {
        Document original = newTestDocument();
        Document clone = original.deepClone();

        clone.getSections().set(0, "Rewritten Intro");
        clone.getSections().add("Appendix");

        assertEquals("Intro", original.getSections().get(0), "original section 0 unchanged");
        assertEquals(2, original.getSections().size(), "original sections size unchanged");
    }

    private static void testMutatingClonePropsDoesNotAffectOriginal() {
        Document original = newTestDocument();
        Document clone = original.deepClone();

        clone.getProps().put("status", "final");
        clone.getProps().put("reviewer", "Asha");

        assertEquals("draft", original.getProps().get("status"), "original status unchanged");
        assertEquals(false, original.getProps().containsKey("reviewer"), "original should not gain new key");
    }

    private static void testMutatingCloneAuthorDoesNotAffectOriginal() {
        Document original = newTestDocument();
        Document clone = original.deepClone();

        clone.getMeta().setAuthor("Someone Else");

        assertEquals("Vinay", original.getMeta().getAuthor(), "original author unchanged");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        boolean equal = (expected == null) ? (actual == null) : expected.equals(actual);
        if (!equal) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
