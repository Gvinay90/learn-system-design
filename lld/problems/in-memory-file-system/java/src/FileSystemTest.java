import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

public class FileSystemTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testMkdirAndLs();
        testMkdirMissingParentFails();
        testWriteReadAndAppend();
        testCdAndRelativePaths();
        testRmNonEmptyDirFails();
        testConcurrentWritesToSameFile();
        System.out.println("All FileSystemTest cases passed.");
    }

    private static void testMkdirAndLs() {
        FileSystem fs = new FileSystem();
        fs.mkdir("/home");
        fs.mkdir("/home/docs");
        fs.createFile("/home/notes.txt", "hello");

        List<Entry> entries = fs.ls("/home");
        assertEquals(2, entries.size(), "entry count");
        assertEquals("docs", entries.get(0).getName(), "first entry name");
        assertEquals(true, entries.get(0).isDir(), "first entry is dir");
        assertEquals("notes.txt", entries.get(1).getName(), "second entry name");
        assertEquals(false, entries.get(1).isDir(), "second entry is dir");
    }

    private static void testMkdirMissingParentFails() {
        FileSystem fs = new FileSystem();
        try {
            fs.mkdir("/a/b");
            throw new AssertionError("expected NotFoundException");
        } catch (FileSystem.NotFoundException expected) {
        }
    }

    private static void testWriteReadAndAppend() {
        FileSystem fs = new FileSystem();
        fs.createFile("/a.txt", "hello");
        assertEquals("hello", fs.readFile("/a.txt"), "initial content");

        fs.writeFile("/a.txt", " world", true);
        assertEquals("hello world", fs.readFile("/a.txt"), "after append");

        fs.writeFile("/a.txt", "reset", false);
        assertEquals("reset", fs.readFile("/a.txt"), "after overwrite");
    }

    private static void testCdAndRelativePaths() {
        FileSystem fs = new FileSystem();
        fs.mkdir("/home");
        fs.mkdir("/home/docs");
        fs.createFile("/home/docs/readme.md", "hi");

        fs.cd("/home/docs");
        assertEquals("/home/docs", fs.pwd(), "pwd after cd");
        assertEquals("hi", fs.readFile("readme.md"), "relative read");

        fs.cd("..");
        assertEquals("/home", fs.pwd(), "pwd after cd ..");
    }

    private static void testRmNonEmptyDirFails() {
        FileSystem fs = new FileSystem();
        fs.mkdir("/home");
        fs.createFile("/home/a.txt", "x");

        try {
            fs.rm("/home", false);
            throw new AssertionError("expected DirectoryNotEmptyException");
        } catch (FileSystem.DirectoryNotEmptyException expected) {
        }
        fs.rm("/home", true);
        try {
            fs.ls("/home");
            throw new AssertionError("expected NotFoundException after recursive rm");
        } catch (FileSystem.NotFoundException expected) {
        }
    }

    private static void testConcurrentWritesToSameFile() {
        FileSystem fs = new FileSystem();
        fs.createFile("/log.txt", "");
        int n = 100;
        CountDownLatch latch = new CountDownLatch(n);
        AtomicInteger started = new AtomicInteger(0);

        for (int i = 0; i < n; i++) {
            new Thread(() -> {
                started.incrementAndGet();
                fs.writeFile("/log.txt", "x", true);
                latch.countDown();
            }).start();
        }
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        assertEquals(n, fs.readFile("/log.txt").length(), "no writes lost under concurrency");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
