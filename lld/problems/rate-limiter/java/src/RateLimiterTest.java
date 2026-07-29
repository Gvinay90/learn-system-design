import java.util.Optional;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

public class RateLimiterTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testTokenBucketAllowsBurstUpToCapacityThenBlocks();
        testTokenBucketRefillsOverTime();
        testSlidingWindowAllowsUpToMaxThenBlocksThenSlides();
        testPerClientIsolation();
        testRateLimiterRegistry();
        testConcurrentAllowDoesNotExceedCapacity();
        System.out.println("All RateLimiterTest cases passed.");
    }

    private static void testTokenBucketAllowsBurstUpToCapacityThenBlocks() {
        TokenBucketLimiter limiter = new TokenBucketLimiter(3, 1);
        FakeClock clock = new FakeClock(0);
        limiter.nowMillis = clock;

        int allowed = 0;
        for (int i = 0; i < 5; i++) {
            if (limiter.allow("client-a")) {
                allowed++;
            }
        }
        assertEquals(3, allowed, "burst allowed up to capacity 3");
    }

    private static void testTokenBucketRefillsOverTime() {
        TokenBucketLimiter limiter = new TokenBucketLimiter(2, 1); // 1 token/sec
        FakeClock clock = new FakeClock(0);
        limiter.nowMillis = clock;

        assertTrue(limiter.allow("client-a"), "1st request allowed");
        assertTrue(limiter.allow("client-a"), "2nd request allowed (burst up to capacity)");
        assertFalse(limiter.allow("client-a"), "3rd immediate request blocked");

        clock.advance(1000); // +1s -> refills exactly 1 token
        assertTrue(limiter.allow("client-a"), "request allowed after 1s refill");
        assertFalse(limiter.allow("client-a"), "next request blocked again (drained)");
    }

    private static void testSlidingWindowAllowsUpToMaxThenBlocksThenSlides() {
        SlidingWindowLimiter limiter = new SlidingWindowLimiter(3, 1000);
        FakeClock clock = new FakeClock(0);
        limiter.nowMillis = clock;

        for (int i = 0; i < 3; i++) {
            assertTrue(limiter.allow("client-a"), "request " + i + " allowed within limit");
        }
        assertFalse(limiter.allow("client-a"), "4th request within window blocked");

        clock.advance(1100); // slide window fully past recorded timestamps
        assertTrue(limiter.allow("client-a"), "request allowed after window slides");
    }

    private static void testPerClientIsolation() {
        TokenBucketLimiter tb = new TokenBucketLimiter(1, 1);
        FakeClock tbClock = new FakeClock(0);
        tb.nowMillis = tbClock;

        assertTrue(tb.allow("client-a"), "client-a first request allowed");
        assertFalse(tb.allow("client-a"), "client-a second immediate request blocked");
        assertTrue(tb.allow("client-b"), "client-b has its own independent bucket");

        SlidingWindowLimiter sw = new SlidingWindowLimiter(1, 1000);
        FakeClock swClock = new FakeClock(0);
        sw.nowMillis = swClock;

        assertTrue(sw.allow("client-a"), "client-a first request allowed");
        assertFalse(sw.allow("client-a"), "client-a second immediate request blocked");
        assertTrue(sw.allow("client-b"), "client-b has its own independent window");
    }

    private static void testRateLimiterRegistry() {
        RateLimiterRegistry registry = new RateLimiterRegistry();
        TokenBucketLimiter free = new TokenBucketLimiter(1, 1);
        SlidingWindowLimiter paid = new SlidingWindowLimiter(100, 1000);

        registry.register("free", free);
        registry.register("paid", paid);

        Optional<RateLimiter> got = registry.getLimiter("free");
        assertTrue(got.isPresent() && got.get() == free, "registry returns registered free-tier limiter");
        assertTrue(!registry.getLimiter("enterprise").isPresent(), "no limiter registered for unknown class");
    }

    private static void testConcurrentAllowDoesNotExceedCapacity() {
        final int capacity = 10;
        final int threads = 100;
        TokenBucketLimiter limiter = new TokenBucketLimiter(capacity, 0); // no refill during the burst
        CountDownLatch latch = new CountDownLatch(threads);
        AtomicInteger allowedCount = new AtomicInteger(0);

        for (int i = 0; i < threads; i++) {
            new Thread(() -> {
                try {
                    if (limiter.allow("client-a")) {
                        allowedCount.incrementAndGet();
                    }
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

        assertEquals(capacity, allowedCount.get(), "exactly capacity requests allowed under concurrency");
    }

    /** Deterministic clock for tests; avoids sleeping real time. */
    private static class FakeClock implements java.util.function.LongSupplier {
        private long millis;

        FakeClock(long millis) {
            this.millis = millis;
        }

        void advance(long deltaMillis) {
            millis += deltaMillis;
        }

        @Override
        public long getAsLong() {
            return millis;
        }
    }

    private static void assertTrue(boolean condition, String label) {
        if (!condition) {
            throw new AssertionError(label + ": expected true");
        }
    }

    private static void assertFalse(boolean condition, String label) {
        if (condition) {
            throw new AssertionError(label + ": expected false");
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
