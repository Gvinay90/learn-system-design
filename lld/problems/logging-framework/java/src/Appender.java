/**
 * A pluggable output destination for log records (Strategy pattern).
 */
public interface Appender {
    void append(Record record);
}
