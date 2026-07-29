public class RetryPolicy {
    private final int maxAttempts;
    private final long delayMillis;

    public RetryPolicy(int maxAttempts, long delayMillis) {
        this.maxAttempts = maxAttempts;
        this.delayMillis = delayMillis;
    }

    public int getMaxAttempts() { return maxAttempts; }

    public long delayFor(int attempt) { return attempt * delayMillis; }
}
