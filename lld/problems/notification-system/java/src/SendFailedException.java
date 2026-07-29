/**
 * Thrown by a Notifier implementation to simulate a delivery failure
 * (e.g. a flaky network call).
 */
public class SendFailedException extends Exception {
    public SendFailedException(String message) {
        super(message);
    }
}
