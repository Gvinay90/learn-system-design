import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Plain assert-based test runner (no JUnit dependency needed), mirroring the
 * Go test cases in taskscheduler_test.go.
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out TaskSchedulerTest` directly.
 */
public class TaskSchedulerTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testJobsRunAndSucceed();
        testPriorityOrdering();
        testDelayedJobDoesNotRunEarly();
        testRetriesOnFailureThenSucceeds();
        testExhaustsRetriesAndFails();
        testUnknownJobResultNotFound();
        testExponentialBackoffCapsAtMax();
        testConcurrentSubmitAndExecute();
        System.out.println("All TaskSchedulerTest cases passed.");
    }

    private static void testJobsRunAndSucceed() {
        Scheduler s = new Scheduler(2);
        s.start();
        try {
            AtomicInteger ran = new AtomicInteger();
            s.submit(new Job("J1", 1, System.currentTimeMillis(), () -> {
                ran.incrementAndGet();
            }));

            Optional<JobResult> res = s.waitForResult("J1", 1000);
            assertTrue(res.isPresent(), "expected job to complete within timeout");
            assertEquals(Status.SUCCEEDED, res.get().getStatus(), "expected SUCCEEDED status");
            assertEquals(1, res.get().getAttempts(), "expected 1 attempt");
            assertEquals(1, ran.get(), "expected task to run exactly once");
        } finally {
            s.stop();
        }
    }

    // Among jobs due at the same time, the higher-priority job should run
    // first. A single worker keeps execution order deterministic, and both
    // jobs share a RunAt in the near future so they're both "due" once the
    // scheduler starts.
    private static void testPriorityOrdering() {
        Scheduler s = new Scheduler(1);

        List<String> order = new ArrayList<>();
        Object mu = new Object();
        long due = System.currentTimeMillis() + 20;

        s.submit(new Job("low", 1, due, () -> {
            synchronized (mu) {
                order.add("low");
            }
        }));
        s.submit(new Job("high", 10, due, () -> {
            synchronized (mu) {
                order.add("high");
            }
        }));

        s.start();
        try {
            assertTrue(s.waitForResult("low", 1000).isPresent(), "expected low priority job to complete");
            assertTrue(s.waitForResult("high", 1000).isPresent(), "expected high priority job to complete");
        } finally {
            s.stop();
        }

        synchronized (mu) {
            assertEquals(2, order.size(), "expected exactly two executed jobs");
            assertEquals("high", order.get(0), "expected high-priority job to run first");
            assertEquals("low", order.get(1), "expected low-priority job to run second");
        }
    }

    private static void testDelayedJobDoesNotRunEarly() {
        Scheduler s = new Scheduler(1);
        s.start();
        try {
            long[] ranAt = new long[1];
            long submittedAt = System.currentTimeMillis();
            long delay = 60;

            s.submit(new Job("delayed", 1, submittedAt + delay, () -> {
                ranAt[0] = System.currentTimeMillis();
            }));

            Optional<JobResult> res = s.waitForResult("delayed", 1000);
            assertTrue(res.isPresent() && res.get().getStatus() == Status.SUCCEEDED,
                    "expected delayed job to eventually succeed");
            assertTrue(ranAt[0] - submittedAt >= delay,
                    "expected job to run no earlier than " + delay + "ms after submit");
        } finally {
            s.stop();
        }
    }

    private static void testRetriesOnFailureThenSucceeds() {
        Scheduler s = new Scheduler(1);
        s.start();
        try {
            AtomicInteger attempts = new AtomicInteger();
            s.submit(new Job("flaky", 1, System.currentTimeMillis(), () -> {
                int n = attempts.incrementAndGet();
                if (n < 3) {
                    throw new RuntimeException("transient failure");
                }
            }, 3, new ExponentialBackoff(2, 10)));

            Optional<JobResult> res = s.waitForResult("flaky", 1000);
            assertTrue(res.isPresent(), "expected flaky job to eventually reach a terminal state");
            assertEquals(Status.SUCCEEDED, res.get().getStatus(), "expected SUCCEEDED after retries");
            assertEquals(3, res.get().getAttempts(), "expected 3 attempts");
        } finally {
            s.stop();
        }
    }

    private static void testExhaustsRetriesAndFails() {
        Scheduler s = new Scheduler(1);
        s.start();
        try {
            AtomicInteger attempts = new AtomicInteger();
            s.submit(new Job("doomed", 1, System.currentTimeMillis(), () -> {
                attempts.incrementAndGet();
                throw new RuntimeException("permanent failure");
            }, 2, new ExponentialBackoff(2, 10)));

            Optional<JobResult> res = s.waitForResult("doomed", 1000);
            assertTrue(res.isPresent(), "expected doomed job to eventually reach a terminal state");
            assertEquals(Status.FAILED, res.get().getStatus(), "expected FAILED status");
            // MaxRetries=2 means 1 initial attempt + 2 retries = 3 attempts total.
            assertEquals(3, res.get().getAttempts(), "expected 3 total attempts (1 initial + 2 retries)");
            assertTrue(res.get().getError() != null
                            && "permanent failure".equals(res.get().getError().getMessage()),
                    "expected the last error preserved");
        } finally {
            s.stop();
        }
    }

    private static void testUnknownJobResultNotFound() {
        Scheduler s = new Scheduler(1);
        assertTrue(s.getResult("nonexistent").isEmpty(), "expected no result for unknown job id");
    }

    private static void testExponentialBackoffCapsAtMax() {
        ExponentialBackoff b = new ExponentialBackoff(10, 30);
        assertEquals(10L, b.nextDelayMillis(0), "expected 10ms for attempt 0");
        assertEquals(20L, b.nextDelayMillis(1), "expected 20ms for attempt 1");
        assertEquals(30L, b.nextDelayMillis(2), "expected backoff capped at 30ms for attempt 2");
        assertEquals(30L, b.nextDelayMillis(5), "expected backoff to stay capped at 30ms for larger attempts");
    }

    // Asserts many jobs submitted concurrently from multiple threads all get
    // executed exactly once by the worker pool.
    private static void testConcurrentSubmitAndExecute() {
        Scheduler s = new Scheduler(4);
        s.start();
        try {
            final int n = 50;
            CountDownLatch latch = new CountDownLatch(n);
            for (int i = 0; i < n; i++) {
                final int idx = i;
                Thread t = new Thread(() -> {
                    String id = jobId(idx);
                    s.submit(new Job(id, idx % 5, System.currentTimeMillis(), () -> {
                    }));
                    latch.countDown();
                });
                t.start();
            }
            try {
                latch.await();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }

            for (int i = 0; i < n; i++) {
                String id = jobId(i);
                Optional<JobResult> res = s.waitForResult(id, 2000);
                assertTrue(res.isPresent() && res.get().getStatus() == Status.SUCCEEDED,
                        "expected job " + id + " to succeed");
            }
        } finally {
            s.stop();
        }
    }

    private static String jobId(int i) {
        return "concurrent-" + (char) ('A' + i % 26) + (i / 26);
    }

    private static void assertTrue(boolean condition, String label) {
        if (!condition) {
            throw new AssertionError(label);
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
