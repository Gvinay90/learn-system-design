public class FileNode implements Node {
    private final String name;
    private final StringBuilder content = new StringBuilder();

    public FileNode(String name, String initialContent) {
        this.name = name;
        this.content.append(initialContent);
    }

    @Override
    public String getName() { return name; }

    @Override
    public boolean isDir() { return false; }

    public synchronized String read() {
        return content.toString();
    }

    public synchronized void write(String data, boolean append) {
        if (!append) {
            content.setLength(0);
        }
        content.append(data);
    }
}
