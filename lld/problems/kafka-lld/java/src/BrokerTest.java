import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicLong;

public class BrokerTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testProduceConsumeInOrder();
        testConsumerGroupsTrackOffsetsIndependently();
        testEdgeCases();
        testConcurrentProduceIntoSamePartition();
        System.out.println("All BrokerTest cases passed.");
    }

    private static void testProduceConsumeInOrder() {
        Broker b = new Broker();
        b.createTopic("orders", 1);

        String[] values = {"a", "b", "c"};
        for (int i = 0; i < values.length; i++) {
            Broker.ProduceResult r = b.produce("orders", "k1", values[i]);
            assertEquals(0, r.partitionId, "expected partition 0");
            assertEquals((long) i, r.offset, "expected offset " + i);
        }

        List<Message> messages = b.consume("g1", "orders", 0, 10);
        assertEquals(3, messages.size(), "expected 3 messages");
        for (int i = 0; i < messages.size(); i++) {
            assertEquals((long) i, messages.get(i).getOffset(), "expected offset " + i);
        }

        List<Message> more = b.consume("g1", "orders", 0, 10);
        assertEquals(0, more.size(), "expected no new messages after auto-commit");
    }

    private static void testConsumerGroupsTrackOffsetsIndependently() {
        Broker b = new Broker();
        b.createTopic("orders", 1);
        for (String v : new String[] {"a", "b", "c"}) {
            b.produce("orders", "", v);
        }

        List<Message> g1Messages = b.consume("group-1", "orders", 0, 2);
        assertEquals(2, g1Messages.size(), "expected 2 messages for group-1");

        List<Message> g2Messages = b.consume("group-2", "orders", 0, 10);
        assertEquals(3, g2Messages.size(), "expected group-2 to see all 3 messages independently");

        List<Message> g1Rest = b.consume("group-1", "orders", 0, 10);
        assertEquals(1, g1Rest.size(), "expected group-1 to resume at its own offset and see 1 message");
        assertEquals(2L, g1Rest.get(0).getOffset(), "expected offset 2");
    }

    private static void testEdgeCases() {
        Broker b = new Broker();
        b.createTopic("orders", 1);
        b.produce("orders", "", "only-message");

        List<Message> messages = b.consume("g1", "orders", 0, 10);
        assertEquals(1, messages.size(), "expected 1 message");

        List<Message> past = b.consume("g1", "orders", 0, 10);
        assertEquals(0, past.size(), "expected no messages past end");

        try {
            b.produce("unknown-topic", "k", "v");
            throw new AssertionError("expected TopicNotFoundException");
        } catch (Broker.TopicNotFoundException expected) {
        }

        try {
            b.consume("g1", "unknown-topic", 0, 10);
            throw new AssertionError("expected TopicNotFoundException");
        } catch (Broker.TopicNotFoundException expected) {
        }

        try {
            b.consume("g1", "orders", 5, 10);
            throw new AssertionError("expected PartitionNotFoundException");
        } catch (Broker.PartitionNotFoundException expected) {
        }
    }

    /**
     * Asserts many threads racing to append to the same partition never lose
     * a message or assign a duplicate offset — the lock in Partition.append
     * must serialize them.
     */
    private static void testConcurrentProduceIntoSamePartition() {
        Broker b = new Broker();
        b.createTopic("orders", 1);

        final int n = 500;
        CountDownLatch latch = new CountDownLatch(n);
        long[] offsets = new long[n];
        AtomicLong idx = new AtomicLong(0);

        for (int i = 0; i < n; i++) {
            new Thread(() -> {
                Broker.ProduceResult r = b.produce("orders", "", "v");
                offsets[(int) idx.getAndIncrement()] = r.offset;
                latch.countDown();
            }).start();
        }
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        Set<Long> seen = new HashSet<>();
        for (long o : offsets) {
            if (!seen.add(o)) {
                throw new AssertionError("duplicate offset " + o);
            }
        }
        assertEquals(n, seen.size(), "expected " + n + " unique offsets");
        for (long i = 0; i < n; i++) {
            if (!seen.contains(i)) {
                throw new AssertionError("gap in offsets: missing " + i);
            }
        }
    }

    private static void assertEquals(long expected, long actual, String label) {
        if (expected != actual) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }

    private static void assertEquals(int expected, int actual, String label) {
        if (expected != actual) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
