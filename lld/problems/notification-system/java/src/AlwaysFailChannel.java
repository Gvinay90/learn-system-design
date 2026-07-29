/** Always fails; useful for exercising "retries exhausted" behavior. */
public class AlwaysFailChannel implements Notifier {
    private final Channel channelType;
    private int attempts = 0;

    public AlwaysFailChannel(Channel channelType) {
        this.channelType = channelType;
    }

    @Override
    public Channel channel() {
        return channelType;
    }

    @Override
    public synchronized void send(String recipient, String message) throws SendFailedException {
        attempts++;
        throw new SendFailedException("simulated permanent failure");
    }

    public synchronized int attempts() {
        return attempts;
    }
}
