import java.util.concurrent.atomic.AtomicInteger;

public class FakePaymentProcessor implements PaymentProcessor {
    private final AtomicInteger callCount = new AtomicInteger(0);
    private final int failTimes;
    private final boolean alwaysFail;

    public FakePaymentProcessor(int failTimes, boolean alwaysFail) {
        this.failTimes = failTimes;
        this.alwaysFail = alwaysFail;
    }

    public static FakePaymentProcessor succeedsImmediately() {
        return new FakePaymentProcessor(0, false);
    }

    public static FakePaymentProcessor failsNTimesThenSucceeds(int n) {
        return new FakePaymentProcessor(n, false);
    }

    public static FakePaymentProcessor alwaysFails() {
        return new FakePaymentProcessor(0, true);
    }

    @Override
    public void process(PaymentRequest request) throws PaymentProcessingException {
        int count = callCount.incrementAndGet();
        if (alwaysFail) {
            throw new PaymentProcessingException("simulated permanent failure");
        }
        if (count <= failTimes) {
            throw new PaymentProcessingException("simulated transient failure");
        }
    }

    public int getCallCount() { return callCount.get(); }
}
