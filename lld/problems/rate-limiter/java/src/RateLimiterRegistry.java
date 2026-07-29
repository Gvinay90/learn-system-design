import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

/**
 * Maps a client class (e.g. free/paid/enterprise) to the RateLimiter
 * instance configured for that tier, so different tiers can even run
 * different algorithms.
 */
public class RateLimiterRegistry {
    private final Map<String, RateLimiter> perClientClass = new HashMap<>();

    public synchronized void register(String clientClass, RateLimiter limiter) {
        perClientClass.put(clientClass, limiter);
    }

    public synchronized Optional<RateLimiter> getLimiter(String clientClass) {
        return Optional.ofNullable(perClientClass.get(clientClass));
    }
}
