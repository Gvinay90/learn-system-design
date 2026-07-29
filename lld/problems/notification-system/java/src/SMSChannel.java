import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/** Simulates sending SMS by recording messages in memory. */
public class SMSChannel implements Notifier {
    private final List<SentMessage> sent = new ArrayList<>();

    @Override
    public Channel channel() {
        return Channel.SMS;
    }

    @Override
    public synchronized void send(String recipient, String message) {
        sent.add(new SentMessage(recipient, message));
    }

    public synchronized List<SentMessage> sent() {
        return Collections.unmodifiableList(new ArrayList<>(sent));
    }
}
