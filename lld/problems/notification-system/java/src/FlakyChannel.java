/**
 * Wraps a delegate Notifier and fails the first N sends before delegating
 * for real. Exists purely so tests can exercise the retry mechanism
 * deterministically.
 */
public class FlakyChannel implements Notifier {
    private final Notifier delegate;
    private int failuresRemaining;
    private int attempts = 0;

    public FlakyChannel(Notifier delegate, int failCount) {
        this.delegate = delegate;
        this.failuresRemaining = failCount;
    }

    @Override
    public Channel channel() {
        return delegate.channel();
    }

    @Override
    public synchronized void send(String recipient, String message) throws SendFailedException {
        attempts++;
        if (failuresRemaining > 0) {
            failuresRemaining--;
            throw new SendFailedException("simulated flaky failure");
        }
        delegate.send(recipient, message);
    }

    public synchronized int attempts() {
        return attempts;
    }
}
