/** Captures the outcome of dispatching to a single channel. */
public class SendResult {
    public final Channel channel;
    public final int attempts;
    public final Exception error; // null on success

    public SendResult(Channel channel, int attempts, Exception error) {
        this.channel = channel;
        this.attempts = attempts;
        this.error = error;
    }

    public boolean success() {
        return error == null;
    }
}
