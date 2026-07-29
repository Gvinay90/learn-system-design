import java.util.HashMap;
import java.util.Map;
import java.util.Objects;

/**
 * Tracks a committed offset per (topic, partition), independent of any other
 * group reading the same topic.
 */
public class ConsumerGroup {
    private static final class PartitionKey {
        final String topic;
        final int partitionId;

        PartitionKey(String topic, int partitionId) {
            this.topic = topic;
            this.partitionId = partitionId;
        }

        @Override
        public boolean equals(Object o) {
            if (!(o instanceof PartitionKey)) return false;
            PartitionKey other = (PartitionKey) o;
            return partitionId == other.partitionId && topic.equals(other.topic);
        }

        @Override
        public int hashCode() {
            return Objects.hash(topic, partitionId);
        }
    }

    private final String id;
    private final Object lock = new Object();
    private final Map<PartitionKey, Long> offsets = new HashMap<>();

    public ConsumerGroup(String id) {
        this.id = id;
    }

    public String getId() {
        return id;
    }

    long committedOffset(String topic, int partitionId) {
        synchronized (lock) {
            return offsets.getOrDefault(new PartitionKey(topic, partitionId), 0L);
        }
    }

    void commit(String topic, int partitionId, long offset) {
        synchronized (lock) {
            offsets.put(new PartitionKey(topic, partitionId), offset);
        }
    }
}
