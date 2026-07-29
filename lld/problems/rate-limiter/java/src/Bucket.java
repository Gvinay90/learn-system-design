/** One client's token-bucket state. */
public class Bucket {
    double tokens;
    long lastRefillMillis;

    Bucket(double tokens, long lastRefillMillis) {
        this.tokens = tokens;
        this.lastRefillMillis = lastRefillMillis;
    }
}
