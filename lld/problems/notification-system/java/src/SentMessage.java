/**
 * Records one successful delivery, used by the in-memory channels so tests
 * can assert on what was actually sent.
 */
public class SentMessage {
    public final String recipient;
    public final String message;

    public SentMessage(String recipient, String message) {
        this.recipient = recipient;
        this.message = message;
    }
}
