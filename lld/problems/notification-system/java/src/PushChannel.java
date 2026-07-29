import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/** Simulates sending a push notification by recording messages in memory. */
public class PushChannel implements Notifier {
    private final List<SentMessage> sent = new ArrayList<>();

    @Override
    public Channel channel() {
        return Channel.PUSH;
    }

    @Override
    public synchronized void send(String recipient, String message) {
        sent.add(new SentMessage(recipient, message));
    }

    public synchronized List<SentMessage> sent() {
        return Collections.unmodifiableList(new ArrayList<>(sent));
    }
}
