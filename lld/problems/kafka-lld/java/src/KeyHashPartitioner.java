import java.nio.charset.StandardCharsets;
import java.util.concurrent.atomic.AtomicLong;

/**
 * Routes by hash(key) % numPartitions so all messages for the same key land
 * on the same partition (preserving per-key order). Keyless messages are
 * spread round robin across partitions.
 */
public class KeyHashPartitioner implements PartitionerStrategy {
    private final AtomicLong roundRobin = new AtomicLong(0);

    @Override
    public int partition(String key, int numPartitions) {
        if (key == null || key.isEmpty()) {
            long n = roundRobin.getAndIncrement();
            return (int) (n % numPartitions);
        }
        long hash = fnv1a32(key);
        return (int) (hash % numPartitions);
    }

    /** FNV-1a 32-bit hash, matching Go's hash/fnv New32a used by the reference implementation. */
    private static long fnv1a32(String s) {
        long hash = 0x811c9dc5L; // offset basis
        long prime = 0x01000193L; // FNV prime
        for (byte b : s.getBytes(StandardCharsets.UTF_8)) {
            hash ^= (b & 0xffL);
            hash = (hash * prime) & 0xffffffffL;
        }
        return hash;
    }
}
