import java.io.PrintStream;

/**
 * Writes formatted records to a PrintStream (System.out by default).
 */
public class ConsoleAppender implements Appender {
    private final PrintStream out;

    public ConsoleAppender() {
        this(System.out);
    }

    public ConsoleAppender(PrintStream out) {
        this.out = out;
    }

    @Override
    public void append(Record record) {
        out.println(record.format());
    }
}
