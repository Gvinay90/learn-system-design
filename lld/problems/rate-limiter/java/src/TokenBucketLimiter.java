import java.util.HashMap;
import java.util.Map;
import java.util.function.LongSupplier;

/**
 * Allows bursts up to capacity tokens, refilling at refillRate tokens/second.
 * Refill is computed lazily on each allow() call rather than via a
 * background thread per client.
 */
public class TokenBucketLimiter implements RateLimiter {
    private final double capacity;
    private final double refillRatePerSecond;
    private final Map<String, Bucket> buckets = new HashMap<>();
    private final Object lock = new Object();

    // Package-private so tests can inject a fake clock; defaults to wall time.
    LongSupplier nowMillis = System::currentTimeMillis;

    public TokenBucketLimiter(double capacity, double refillRatePerSecond) {
        this.capacity = capacity;
        this.refillRatePerSecond = refillRatePerSecond;
    }

    private void refill(Bucket b, long now) {
        double elapsedSeconds = (now - b.lastRefillMillis) / 1000.0;
        if (elapsedSeconds <= 0) {
            return;
        }
        b.tokens = Math.min(capacity, b.tokens + elapsedSeconds * refillRatePerSecond);
        b.lastRefillMillis = now;
    }

    @Override
    public boolean allow(String clientId) {
        synchronized (lock) {
            long now = nowMillis.getAsLong();
            Bucket b = buckets.get(clientId);
            if (b == null) {
                b = new Bucket(capacity, now);
                buckets.put(clientId, b);
            } else {
                refill(b, now);
            }

            if (b.tokens < 1) {
                return false;
            }
            b.tokens -= 1;
            return true;
        }
    }
}
