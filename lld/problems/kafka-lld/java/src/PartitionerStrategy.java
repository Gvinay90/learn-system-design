/** Selects a partition index for a produced message. */
public interface PartitionerStrategy {
    int partition(String key, int numPartitions);
}
