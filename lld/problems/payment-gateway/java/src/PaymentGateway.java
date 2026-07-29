import java.time.Instant;
import java.util.concurrent.atomic.AtomicInteger;

public class PaymentGateway {
    private final PaymentProcessor processor;
    private final RetryPolicy retryPolicy;
    private final IdempotencyStore store = new IdempotencyStore();
    private final Ledger ledger = new Ledger();
    private final AtomicInteger seq = new AtomicInteger(0);

    public PaymentGateway(PaymentProcessor processor, RetryPolicy retryPolicy) {
        this.processor = processor;
        this.retryPolicy = retryPolicy;
    }

    public Ledger getLedger() { return ledger; }

    // Same idempotencyKey always returns the same terminal result, SUCCESS
    // or FAILED, without reprocessing (some real gateways allow retrying a
    // failed key; this exercise keeps the simpler "terminal result is final"
    // semantics).
    public PaymentResult charge(PaymentRequest request) {
        IdempotencyStore.Reservation reservation = store.reserveOrWait(request.getIdempotencyKey());
        if (!reservation.isOwner) {
            return reservation.existingResult;
        }

        PaymentResult result = processWithRetry(request);
        store.complete(request.getIdempotencyKey(), result);
        return result;
    }

    private PaymentResult processWithRetry(PaymentRequest request) {
        String id = "PAY-" + seq.incrementAndGet();
        PaymentResult result = new PaymentResult(id, request);

        for (int attempt = 1; attempt <= retryPolicy.getMaxAttempts(); attempt++) {
            String error = null;
            try {
                processor.process(request);
            } catch (PaymentProcessor.PaymentProcessingException e) {
                error = e.getMessage();
            }
            result.addAttempt(new Attempt(attempt, error == null, error, Instant.now()));

            if (error == null) {
                result.setStatus(PaymentStatus.SUCCESS);
                ledger.record(new LedgerEntry(id, request.getPayerId(), request.getPayeeId(),
                        request.getAmount(), Instant.now()));
                return result;
            }

            if (attempt < retryPolicy.getMaxAttempts()) {
                sleep(retryPolicy.delayFor(attempt));
            }
        }

        result.setStatus(PaymentStatus.FAILED);
        return result;
    }

    private static void sleep(long millis) {
        try {
            Thread.sleep(millis);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
