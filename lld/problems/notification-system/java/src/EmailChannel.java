import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/** Simulates sending email by recording messages in memory. */
public class EmailChannel implements Notifier {
    private final List<SentMessage> sent = new ArrayList<>();

    @Override
    public Channel channel() {
        return Channel.EMAIL;
    }

    @Override
    public synchronized void send(String recipient, String message) {
        sent.add(new SentMessage(recipient, message));
    }

    public synchronized List<SentMessage> sent() {
        return Collections.unmodifiableList(new ArrayList<>(sent));
    }
}
