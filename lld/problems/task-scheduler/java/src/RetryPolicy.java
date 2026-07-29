/** Computes the delay before the (attempt+1)-th retry, in milliseconds. */
public interface RetryPolicy {
    long nextDelayMillis(int attempt);
}
