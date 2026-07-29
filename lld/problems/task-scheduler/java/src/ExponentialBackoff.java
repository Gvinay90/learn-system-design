public class ExponentialBackoff implements RetryPolicy {
    private final long baseMillis;
    private final long maxMillis;

    public ExponentialBackoff(long baseMillis, long maxMillis) {
        this.baseMillis = baseMillis;
        this.maxMillis = maxMillis;
    }

    @Override
    public long nextDelayMillis(int attempt) {
        long d = baseMillis;
        for (int i = 0; i < attempt; i++) {
            d *= 2;
            if (d > maxMillis) {
                return maxMillis;
            }
        }
        return d;
    }
}
