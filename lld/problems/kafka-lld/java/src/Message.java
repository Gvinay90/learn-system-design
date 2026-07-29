/**
 * A single record stored in a partition: a key/value pair plus the
 * sequential offset it was assigned at append time.
 */
public class Message {
    private final String key;
    private final String value;
    private final long offset;

    public Message(String key, String value, long offset) {
        this.key = key;
        this.value = value;
        this.offset = offset;
    }

    public String getKey() {
        return key;
    }

    public String getValue() {
        return value;
    }

    public long getOffset() {
        return offset;
    }
}
