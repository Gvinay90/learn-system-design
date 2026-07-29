import java.util.Optional;
import java.util.concurrent.CountDownLatch;

public class LRUCacheTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testGetPutUpdateInPlace();
        testEvictsActualLRU();
        testMissingKeyAndCapacityOne();
        testConcurrentAccess();
        System.out.println("All LRUCacheTest cases passed.");
    }

    private static void testGetPutUpdateInPlace() {
        LRUCache<String, Integer> cache = new LRUCache<>(2);
        cache.put("a", 1);
        cache.put("b", 2);
        assertEquals(Optional.of(1), cache.get("a"), "get a");

        cache.put("a", 100);
        assertEquals(Optional.of(100), cache.get("a"), "get a after update");
        assertEquals(2, cache.size(), "size after update");
    }

    private static void testEvictsActualLRU() {
        LRUCache<String, Integer> cache = new LRUCache<>(2);
        cache.put("a", 1);
        cache.put("b", 2);
        cache.get("a");
        cache.put("c", 3);

        assertEquals(Optional.empty(), cache.get("b"), "b should be evicted");
        assertEquals(Optional.of(1), cache.get("a"), "a should survive");
        assertEquals(Optional.of(3), cache.get("c"), "c should survive");
        assertEquals(2, cache.size(), "size after eviction");
    }

    private static void testMissingKeyAndCapacityOne() {
        LRUCache<String, Integer> cache = new LRUCache<>(1);
        assertEquals(Optional.empty(), cache.get("missing"), "miss on empty cache");
        assertEquals(0, cache.size(), "size stays 0 after a miss");

        cache.put("a", 1);
        cache.put("b", 2);
        assertEquals(Optional.empty(), cache.get("a"), "a evicted at capacity 1");
        assertEquals(Optional.of(2), cache.get("b"), "b present at capacity 1");
    }

    private static void testConcurrentAccess() {
        final int capacity = 50;
        final int threads = 100;
        LRUCache<Integer, Integer> cache = new LRUCache<>(capacity);
        CountDownLatch latch = new CountDownLatch(threads);

        for (int i = 0; i < threads; i++) {
            final int n = i;
            new Thread(() -> {
                try {
                    cache.put(n, n * n);
                    cache.get(n);
                    cache.get(n - 1);
                } finally {
                    latch.countDown();
                }
            }).start();
        }
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        if (cache.size() > capacity) {
            throw new AssertionError("size " + cache.size() + " exceeds capacity " + capacity);
        }
        assertEquals(cache.size(), cache.listLength(), "map size and list length must agree");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
