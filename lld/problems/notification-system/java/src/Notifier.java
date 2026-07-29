/**
 * Strategy interface every delivery channel implements.
 */
public interface Notifier {
    Channel channel();

    /**
     * Attempts to deliver message to recipient.
     *
     * @throws SendFailedException to simulate a delivery failure.
     */
    void send(String recipient, String message) throws SendFailedException;
}
