import java.io.FileWriter;
import java.io.IOException;
import java.io.PrintWriter;

/**
 * Appends formatted records, one per line, to a file at a caller-supplied
 * path. The file is opened once (append mode) and kept open for the
 * lifetime of the appender.
 */
public class FileAppender implements Appender {
    private final PrintWriter writer;

    public FileAppender(String path) throws IOException {
        this.writer = new PrintWriter(new FileWriter(path, true));
    }

    @Override
    public synchronized void append(Record record) {
        writer.println(record.format());
        writer.flush();
    }

    public synchronized void close() {
        writer.close();
    }
}
