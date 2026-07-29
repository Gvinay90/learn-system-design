import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReadWriteLock;
import java.util.concurrent.locks.ReentrantReadWriteLock;

/** Top-level orchestrator: it owns topics and consumer groups. */
public class Broker {
    public static class TopicExistsException extends RuntimeException {
        public TopicExistsException() {
            super("topic already exists");
        }
    }

    public static class TopicNotFoundException extends RuntimeException {
        public TopicNotFoundException() {
            super("topic not found");
        }
    }

    public static class PartitionNotFoundException extends RuntimeException {
        public PartitionNotFoundException() {
            super("partition not found");
        }
    }

    /** Result of a produce call: the partition chosen and the offset assigned. */
    public static final class ProduceResult {
        public final int partitionId;
        public final long offset;

        public ProduceResult(int partitionId, long offset) {
            this.partitionId = partitionId;
            this.offset = offset;
        }
    }

    private final ReadWriteLock mu = new ReentrantReadWriteLock();
    private final Map<String, Topic> topics = new ConcurrentHashMap<>();
    private final Map<String, ConsumerGroup> consumerGroups = new ConcurrentHashMap<>();

    public void createTopic(String name, int numPartitions) {
        mu.writeLock().lock();
        try {
            if (topics.containsKey(name)) {
                throw new TopicExistsException();
            }
            Topic topic = new Topic(name, numPartitions, new KeyHashPartitioner());
            topics.put(name, topic);
        } finally {
            mu.writeLock().unlock();
        }
    }

    public ProduceResult produce(String topicName, String key, String value) {
        Topic topic;
        mu.readLock().lock();
        try {
            topic = topics.get(topicName);
        } finally {
            mu.readLock().unlock();
        }
        if (topic == null) {
            throw new TopicNotFoundException();
        }
        int partitionId = topic.selectPartition(key);
        long offset = topic.getPartitions().get(partitionId).append(key, value);
        return new ProduceResult(partitionId, offset);
    }

    private ConsumerGroup getOrCreateConsumerGroup(String groupId) {
        return consumerGroups.computeIfAbsent(groupId, ConsumerGroup::new);
    }

    /**
     * Returns up to maxMessages messages after the group's last committed
     * offset for (topic, partitionId), then auto-commits past the returned
     * batch. Auto-commit-on-read keeps this exercise's API small and
     * deterministic; the alternative is an explicit commitOffset call with
     * at-least-once redelivery on crash.
     */
    public List<Message> consume(String groupId, String topicName, int partitionId, int maxMessages) {
        Topic topic;
        mu.readLock().lock();
        try {
            topic = topics.get(topicName);
        } finally {
            mu.readLock().unlock();
        }
        if (topic == null) {
            throw new TopicNotFoundException();
        }
        if (partitionId < 0 || partitionId >= topic.getPartitions().size()) {
            throw new PartitionNotFoundException();
        }

        ConsumerGroup group = getOrCreateConsumerGroup(groupId);
        long fromOffset = group.committedOffset(topicName, partitionId);
        List<Message> messages = topic.getPartitions().get(partitionId).read(fromOffset, maxMessages);
        if (!messages.isEmpty()) {
            group.commit(topicName, partitionId, messages.get(messages.size() - 1).getOffset() + 1);
        }
        return messages;
    }

    public long committedOffset(String groupId, String topicName, int partitionId) {
        ConsumerGroup group = getOrCreateConsumerGroup(groupId);
        return group.committedOffset(topicName, partitionId);
    }
}
