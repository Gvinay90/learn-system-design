/**
 * Common strategy interface every rate-limiting algorithm implements, so
 * calling code never branches on which algorithm is active.
 */
public interface RateLimiter {
    /** Decides whether a request from clientID is permitted right now. */
    boolean allow(String clientId);
}
