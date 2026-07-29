import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

/**
 * Dispatches records to its appenders when the record's level meets or
 * exceeds the logger's minimum threshold. Thread-safe: appender dispatch is
 * synchronized so concurrent log() calls don't interleave writes.
 */
public class Logger {
    private Level level;
    private final List<Appender> appenders = new ArrayList<>();

    public Logger(Level level, Appender... appenders) {
        this.level = level;
        for (Appender a : appenders) {
            this.appenders.add(a);
        }
    }

    public synchronized void addAppender(Appender appender) {
        appenders.add(appender);
    }

    public synchronized void setLevel(Level level) {
        this.level = level;
    }

    public synchronized void log(Level level, String message) {
        if (level.compareTo(this.level) < 0) {
            return;
        }
        Record record = new Record(Instant.now(), level, message);
        for (Appender a : appenders) {
            a.append(record);
        }
    }

    public void debug(String message) {
        log(Level.DEBUG, message);
    }

    public void info(String message) {
        log(Level.INFO, message);
    }

    public void warn(String message) {
        log(Level.WARN, message);
    }

    public void error(String message) {
        log(Level.ERROR, message);
    }
}
