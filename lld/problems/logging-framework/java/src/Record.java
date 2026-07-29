import java.time.Instant;

/**
 * A single formatted log entry handed to every Appender.
 */
public class Record {
    public final Instant timestamp;
    public final Level level;
    public final String message;

    public Record(Instant timestamp, Level level, String message) {
        this.timestamp = timestamp;
        this.level = level;
        this.message = message;
    }

    /** Renders as "<ISO-8601 timestamp> [LEVEL] message". */
    public String format() {
        return timestamp + " [" + level + "] " + message;
    }
}
