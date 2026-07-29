import java.util.List;
import java.util.concurrent.CountDownLatch;

public class PaymentGatewayTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testHappyPathChargeRecordsLedgerEntry();
        testSameIdempotencyKeyDoesNotReprocess();
        testRetryPolicySucceedsAfterTransientFailures();
        testRetryPolicyExhaustsToFailed();
        testConcurrentChargesSameKeyProcessOnce();
        System.out.println("All PaymentGatewayTest cases passed.");
    }

    private static PaymentRequest request(String key) {
        return new PaymentRequest(key, "payer-1", "payee-1", 100.0, "INR");
    }

    private static void testHappyPathChargeRecordsLedgerEntry() {
        PaymentGateway gateway = new PaymentGateway(FakePaymentProcessor.succeedsImmediately(), new RetryPolicy(3, 1));
        PaymentResult result = gateway.charge(request("key-1"));
        assertEquals(PaymentStatus.SUCCESS, result.getStatus(), "status");

        List<LedgerEntry> entries = gateway.getLedger().getEntries();
        assertEquals(1, entries.size(), "ledger entry count");
        assertEquals(result.getId(), entries.get(0).getPaymentId(), "ledger payment id");
        assertEquals(100.0, entries.get(0).getAmount(), "ledger amount");
    }

    private static void testSameIdempotencyKeyDoesNotReprocess() {
        FakePaymentProcessor processor = FakePaymentProcessor.succeedsImmediately();
        PaymentGateway gateway = new PaymentGateway(processor, new RetryPolicy(3, 1));

        PaymentResult first = gateway.charge(request("key-2"));
        PaymentResult second = gateway.charge(request("key-2"));

        assertEquals(1, processor.getCallCount(), "processor call count");
        assertEquals(first.getId(), second.getId(), "cached result id");
        assertEquals(first.getStatus(), second.getStatus(), "cached result status");
    }

    private static void testRetryPolicySucceedsAfterTransientFailures() {
        PaymentGateway gateway = new PaymentGateway(FakePaymentProcessor.failsNTimesThenSucceeds(2), new RetryPolicy(3, 1));
        PaymentResult result = gateway.charge(request("key-3"));

        assertEquals(PaymentStatus.SUCCESS, result.getStatus(), "status after retries");
        assertEquals(3, result.getAttempts().size(), "attempt count");
    }

    private static void testRetryPolicyExhaustsToFailed() {
        FakePaymentProcessor processor = FakePaymentProcessor.alwaysFails();
        PaymentGateway gateway = new PaymentGateway(processor, new RetryPolicy(3, 1));

        PaymentResult result = gateway.charge(request("key-4"));
        assertEquals(PaymentStatus.FAILED, result.getStatus(), "status after exhausted retries");
        assertEquals(3, result.getAttempts().size(), "attempt count");

        PaymentResult again = gateway.charge(request("key-4"));
        assertEquals(PaymentStatus.FAILED, again.getStatus(), "cached failed result");
        assertEquals(3, processor.getCallCount(), "no reprocessing on repeat charge");
    }

    private static void testConcurrentChargesSameKeyProcessOnce() {
        FakePaymentProcessor processor = FakePaymentProcessor.succeedsImmediately();
        PaymentGateway gateway = new PaymentGateway(processor, new RetryPolicy(3, 1));

        int n = 20;
        CountDownLatch latch = new CountDownLatch(n);
        PaymentResult[] results = new PaymentResult[n];

        for (int i = 0; i < n; i++) {
            int idx = i;
            new Thread(() -> {
                try {
                    results[idx] = gateway.charge(request("shared-key"));
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

        assertEquals(1, processor.getCallCount(), "processor invoked exactly once");
        String firstId = results[0].getId();
        for (int i = 0; i < n; i++) {
            assertEquals(firstId, results[i].getId(), "goroutine " + i + " result id");
            assertEquals(PaymentStatus.SUCCESS, results[i].getStatus(), "goroutine " + i + " result status");
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
