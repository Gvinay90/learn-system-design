import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.function.LongSupplier;

/**
 * Allows at most limit requests in any trailing window (a sliding log),
 * giving a smooth, hard cap with no boundary burst allowance.
 */
public class SlidingWindowLimiter implements RateLimiter {
    private final int limit;
    private final long windowMillis;
    private final Map<String, Window> windows = new HashMap<>();
    private final Object lock = new Object();

    // Package-private so tests can inject a fake clock; defaults to wall time.
    LongSupplier nowMillis = System::currentTimeMillis;

    public SlidingWindowLimiter(int limit, long windowMillis) {
        this.limit = limit;
        this.windowMillis = windowMillis;
    }

    @Override
    public boolean allow(String clientId) {
        synchronized (lock) {
            long now = nowMillis.getAsLong();
            Window w = windows.computeIfAbsent(clientId, id -> new Window());

            long cutoff = now - windowMillis;
            Iterator<Long> it = w.timestampsMillis.iterator();
            while (it.hasNext()) {
                if (it.next() <= cutoff) {
                    it.remove();
                } else {
                    break; // oldest-first deque; once we hit a fresh entry the rest are fresh too
                }
            }

            if (w.timestampsMillis.size() >= limit) {
                return false;
            }
            w.timestampsMillis.addLast(now);
            return true;
        }
    }
}
