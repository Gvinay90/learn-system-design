import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out ConcurrencyTest` directly.
 */
public class ConcurrencyTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testNoLostOrDuplicatedItems();
        testSingleProducerSingleWorker();
        System.out.println("All ConcurrencyTest cases passed.");
    }

    private static void testNoLostOrDuplicatedItems() {
        final int numProducers = 5;
        final int itemsPerProducer = 200;
        final int numWorkers = 8;
        final int bufferCapacity = 16;

        ConcurrentHashMap<String, Integer> counts = new ConcurrentHashMap<>();
        BoundedPipeline pipeline = new BoundedPipeline(numWorkers, bufferCapacity);
        pipeline.run(numProducers, itemsPerProducer, item ->
                counts.merge(item.id(), 1, Integer::sum));

        int wantTotal = numProducers * itemsPerProducer;
        assertEquals(wantTotal, counts.size(), "unique items consumed");

        for (int p = 0; p < numProducers; p++) {
            for (int seq = 0; seq < itemsPerProducer; seq++) {
                String id = "p" + p + "-" + seq;
                Integer count = counts.get(id);
                if (count == null) {
                    throw new AssertionError("item " + id + " was never consumed");
                }
                assertEquals(1, count, "delivery count for " + id);
            }
        }
    }

    private static void testSingleProducerSingleWorker() {
        Set<String> consumed = ConcurrentHashMap.newKeySet();
        BoundedPipeline pipeline = new BoundedPipeline(1, 1);
        pipeline.run(1, 50, item -> consumed.add(item.id()));
        assertEquals(50, consumed.size(), "items consumed with single producer/worker");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
