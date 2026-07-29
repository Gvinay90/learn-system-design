import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicIntegerArray;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out SingletonTest` directly.
 */
public class SingletonTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testSameInstanceReturned();
        testSetAndGet();
        testConcurrentFirstAccess();
        System.out.println("All SingletonTest cases passed.");
    }

    private static void testSameInstanceReturned() {
        AppConfig a = AppConfig.getInstance();
        AppConfig b = AppConfig.getInstance();
        assertEquals(true, a == b, "getInstance must return the same reference");
        assertEquals(a.getId(), b.getId(), "instance ids must match");
    }

    private static void testSetAndGet() {
        AppConfig config = AppConfig.getInstance();
        config.set("region", "us-east-1");
        assertEquals("us-east-1", config.get("region"), "set/get roundtrip");
    }

    private static void testConcurrentFirstAccess() {
        final int threads = 50;
        AtomicIntegerArray ids = new AtomicIntegerArray(threads);
        CountDownLatch latch = new CountDownLatch(threads);

        for (int i = 0; i < threads; i++) {
            final int idx = i;
            new Thread(() -> {
                ids.set(idx, AppConfig.getInstance().getId());
                latch.countDown();
            }).start();
        }
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        int first = ids.get(0);
        for (int i = 0; i < threads; i++) {
            assertEquals(first, ids.get(i), "every thread must observe the same singleton id");
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
