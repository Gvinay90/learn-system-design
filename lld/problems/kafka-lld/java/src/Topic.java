import java.util.ArrayList;
import java.util.List;

/** Owns a fixed set of partitions created at topic-creation time. */
public class Topic {
    public static class InvalidPartitionCountException extends RuntimeException {
        public InvalidPartitionCountException() {
            super("numPartitions must be > 0");
        }
    }

    private final String name;
    private final List<Partition> partitions;
    private final PartitionerStrategy partitioner;

    public Topic(String name, int numPartitions, PartitionerStrategy partitioner) {
        if (numPartitions <= 0) {
            throw new InvalidPartitionCountException();
        }
        this.name = name;
        this.partitioner = partitioner;
        this.partitions = new ArrayList<>(numPartitions);
        for (int i = 0; i < numPartitions; i++) {
            this.partitions.add(new Partition());
        }
    }

    public String getName() {
        return name;
    }

    public List<Partition> getPartitions() {
        return partitions;
    }

    public int selectPartition(String key) {
        return partitioner.partition(key, partitions.size());
    }
}
