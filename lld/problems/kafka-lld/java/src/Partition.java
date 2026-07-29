import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * An append-only log of messages. Offsets are assigned sequentially starting
 * at 0, in append order.
 */
public class Partition {
    private final Object lock = new Object();
    private final List<Message> messages = new ArrayList<>();

    /**
     * Assigns the next sequential offset to (key, value) while holding the
     * partition's lock, so concurrent producers can never race on offset
     * assignment (no lost updates, no duplicate offsets).
     */
    public long append(String key, String value) {
        synchronized (lock) {
            long offset = messages.size();
            messages.add(new Message(key, value, offset));
            return offset;
        }
    }

    public List<Message> read(long fromOffset, int maxMessages) {
        synchronized (lock) {
            int size = messages.size();
            if (fromOffset < 0 || fromOffset >= size) {
                return Collections.emptyList();
            }
            long end = size;
            if (maxMessages > 0 && fromOffset + maxMessages < end) {
                end = fromOffset + maxMessages;
            }
            return new ArrayList<>(messages.subList((int) fromOffset, (int) end));
        }
    }

    public long len() {
        synchronized (lock) {
            return messages.size();
        }
    }
}
