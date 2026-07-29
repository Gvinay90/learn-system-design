/** Configures the retry-on-failure wrapper around a channel send. */
public class RetryPolicy {
    public final int maxAttempts; // total attempts, including the first; must be >= 1
    public final long delayMillis; // fixed delay between attempts

    public RetryPolicy(int maxAttempts, long delayMillis) {
        this.maxAttempts = maxAttempts;
        this.delayMillis = delayMillis;
    }

    /** Retries up to 3 times total with a 10ms delay, suitable for tests. */
    public static RetryPolicy defaultPolicy() {
        return new RetryPolicy(3, 10);
    }
}
