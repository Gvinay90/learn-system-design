import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Thread-safe singleton using the initialization-on-demand holder idiom:
 * the JVM class-loading guarantee makes Holder.INSTANCE construction
 * atomic without any explicit locking.
 */
public class AppConfig {

    private static final AtomicInteger SEQ = new AtomicInteger(0);

    private final int id;
    private final Map<String, String> settings = new HashMap<>();

    private AppConfig() {
        this.id = SEQ.incrementAndGet();
    }

    private static class Holder {
        private static final AppConfig INSTANCE = new AppConfig();
    }

    public static AppConfig getInstance() {
        return Holder.INSTANCE;
    }

    public int getId() {
        return id;
    }

    public synchronized void set(String key, String value) {
        settings.put(key, value);
    }

    public synchronized String get(String key) {
        return settings.get(key);
    }
}
