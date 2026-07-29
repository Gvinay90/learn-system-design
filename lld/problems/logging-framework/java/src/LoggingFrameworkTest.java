import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.nio.file.Files;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CountDownLatch;

public class LoggingFrameworkTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testLevelFiltering();
        testMultipleAppendersAllReceiveRecord();
        testAddAppenderAfterConstruction();
        testSetLevelChangesThreshold();
        testRecordFormatContainsLevelAndMessage();
        testFileAppenderWritesToTempFile();
        testConcurrentLogCalls();
        System.out.println("All LoggingFrameworkTest cases passed.");
    }

    private static void testLevelFiltering() {
        assertFilterResult(Level.DEBUG, Level.DEBUG, 1, "debug threshold logs debug");
        assertFilterResult(Level.INFO, Level.DEBUG, 0, "info threshold drops debug");
        assertFilterResult(Level.INFO, Level.INFO, 1, "info threshold logs info");
        assertFilterResult(Level.WARN, Level.INFO, 0, "warn threshold drops info");
        assertFilterResult(Level.WARN, Level.WARN, 1, "warn threshold logs warn");
        assertFilterResult(Level.ERROR, Level.WARN, 0, "error threshold drops warn");
        assertFilterResult(Level.ERROR, Level.ERROR, 1, "error threshold logs error");
    }

    private static void assertFilterResult(Level threshold, Level logAt, int wantCount, String label) {
        MockAppender appender = new MockAppender();
        Logger logger = new Logger(threshold, appender);
        logger.log(logAt, "message");
        assertEquals(wantCount, appender.count(), label);
    }

    private static void testMultipleAppendersAllReceiveRecord() {
        MockAppender a1 = new MockAppender();
        MockAppender a2 = new MockAppender();
        MockAppender a3 = new MockAppender();
        Logger logger = new Logger(Level.INFO, a1, a2, a3);

        logger.info("hello");
        logger.debug("should be filtered");

        assertEquals(1, a1.count(), "appender 1 record count");
        assertEquals(1, a2.count(), "appender 2 record count");
        assertEquals(1, a3.count(), "appender 3 record count");
    }

    private static void testAddAppenderAfterConstruction() {
        MockAppender a1 = new MockAppender();
        Logger logger = new Logger(Level.DEBUG, a1);

        MockAppender a2 = new MockAppender();
        logger.addAppender(a2);

        logger.info("hi");
        assertEquals(1, a1.count(), "a1 count");
        assertEquals(1, a2.count(), "a2 count");
    }

    private static void testSetLevelChangesThreshold() {
        MockAppender appender = new MockAppender();
        Logger logger = new Logger(Level.ERROR, appender);

        logger.warn("filtered");
        assertEquals(0, appender.count(), "warn filtered at ERROR threshold");

        logger.setLevel(Level.WARN);
        logger.warn("passes now");
        assertEquals(1, appender.count(), "warn passes after lowering threshold");
    }

    private static void testRecordFormatContainsLevelAndMessage() {
        MockAppender appender = new MockAppender();
        Logger logger = new Logger(Level.DEBUG, appender);

        logger.error("disk on fire");

        String formatted = appender.get(0).format();
        if (!formatted.contains("[ERROR]")) {
            throw new AssertionError("expected formatted record to contain level tag, got: " + formatted);
        }
        if (!formatted.contains("disk on fire")) {
            throw new AssertionError("expected formatted record to contain message, got: " + formatted);
        }
    }

    private static void testFileAppenderWritesToTempFile() {
        try {
            File tempDir = Files.createTempDirectory("logging-framework-test").toFile();
            tempDir.deleteOnExit();
            File logFile = new File(tempDir, "app.log");
            logFile.deleteOnExit();

            FileAppender appender = new FileAppender(logFile.getAbsolutePath());
            Logger logger = new Logger(Level.INFO, appender);
            logger.info("first line");
            logger.warn("second line");
            logger.debug("filtered, should not appear");
            appender.close();

            List<String> lines = readLines(logFile);
            assertEquals(2, lines.size(), "number of lines written to file");
            if (!lines.get(0).contains("first line")) {
                throw new AssertionError("expected first line to contain 'first line', got: " + lines.get(0));
            }
            if (!lines.get(1).contains("second line")) {
                throw new AssertionError("expected second line to contain 'second line', got: " + lines.get(1));
            }
            for (String line : lines) {
                if (line.contains("filtered, should not appear")) {
                    throw new AssertionError("did not expect filtered debug message in file");
                }
            }
        } catch (IOException e) {
            throw new AssertionError("unexpected IOException in file appender test", e);
        }
    }

    private static List<String> readLines(File file) throws IOException {
        List<String> lines = new ArrayList<>();
        try (BufferedReader reader = new BufferedReader(new FileReader(file))) {
            String line;
            while ((line = reader.readLine()) != null) {
                lines.add(line);
            }
        }
        return lines;
    }

    private static void testConcurrentLogCalls() {
        final int threads = 100;
        final int perThread = 20;

        MockAppender appender = new MockAppender();
        Logger logger = new Logger(Level.DEBUG, appender);
        CountDownLatch latch = new CountDownLatch(threads);

        for (int i = 0; i < threads; i++) {
            new Thread(() -> {
                try {
                    for (int j = 0; j < perThread; j++) {
                        logger.info("concurrent message");
                    }
                } finally {
                    latch.countDown();
                }
            }).start();
        }
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        assertEquals(threads * perThread, appender.count(), "concurrent log call count");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
