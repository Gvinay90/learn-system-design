import java.util.ArrayDeque;
import java.util.Deque;

/** One client's sliding-log state: timestamps (ms) still inside the window. */
public class Window {
    final Deque<Long> timestampsMillis = new ArrayDeque<>();
}
