public class Entry {
    private final String name;
    private final boolean dir;

    public Entry(String name, boolean dir) {
        this.name = name;
        this.dir = dir;
    }

    public String getName() { return name; }
    public boolean isDir() { return dir; }
}
