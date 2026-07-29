/**
 * Unit of work submitted to the {@link Scheduler}. Mirrors Go's
 * {@code func() error} — throwing signals failure, returning normally
 * signals success.
 */
@FunctionalInterface
public interface Task {
    void run() throws Exception;
}
