import java.util.ArrayList;
import java.util.List;

/**
 * Test double that records every Record it receives. Synchronized so it is
 * safe to share across threads in the concurrency test.
 */
public class MockAppender implements Appender {
    private final List<Record> records = new ArrayList<>();

    @Override
    public synchronized void append(Record record) {
        records.add(record);
    }

    public synchronized int count() {
        return records.size();
    }

    public synchronized Record get(int index) {
        return records.get(index);
    }
}
