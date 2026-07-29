import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;

/**
 * In-memory hierarchical file system: a Composite tree of DirectoryNode/FileNode
 * with path resolution ("/", "." and "..") and the usual mkdir/create/write/read/ls/rm/cd ops.
 */
public class FileSystem {

    public static class NotFoundException extends RuntimeException {
        public NotFoundException(String path) { super("no such file or directory: " + path); }
    }
    public static class AlreadyExistsException extends RuntimeException {
        public AlreadyExistsException(String path) { super("already exists: " + path); }
    }
    public static class NotDirectoryException extends RuntimeException {
        public NotDirectoryException(String path) { super("not a directory: " + path); }
    }
    public static class IsDirectoryException extends RuntimeException {
        public IsDirectoryException(String path) { super("is a directory: " + path); }
    }
    public static class DirectoryNotEmptyException extends RuntimeException {
        public DirectoryNotEmptyException(String path) { super("directory not empty: " + path); }
    }

    private final DirectoryNode root = new DirectoryNode("/");
    private volatile List<String> cwd = new ArrayList<>();

    private List<String> splitPath(String path, List<String> base) {
        if (path == null || path.isEmpty()) {
            throw new IllegalArgumentException("invalid path");
        }
        List<String> parts = new ArrayList<>();
        if (!path.startsWith("/")) {
            parts.addAll(base);
        }
        for (String seg : path.split("/")) {
            if (seg.isEmpty() || seg.equals(".")) continue;
            if (seg.equals("..")) {
                if (!parts.isEmpty()) parts.remove(parts.size() - 1);
            } else {
                parts.add(seg);
            }
        }
        return parts;
    }

    private synchronized DirectoryNode resolveParentDir(List<String> parts) {
        DirectoryNode dir = root;
        for (String seg : parts) {
            Node child = dir.children.get(seg);
            if (child == null) throw new NotFoundException(seg);
            if (!child.isDir()) throw new NotDirectoryException(seg);
            dir = (DirectoryNode) child;
        }
        return dir;
    }

    private synchronized Node resolve(List<String> parts) {
        if (parts.isEmpty()) return root;
        DirectoryNode parent = resolveParentDir(parts.subList(0, parts.size() - 1));
        String name = parts.get(parts.size() - 1);
        Node node = parent.children.get(name);
        if (node == null) throw new NotFoundException(name);
        return node;
    }

    public synchronized void mkdir(String path) {
        List<String> parts = splitPath(path, cwd);
        if (parts.isEmpty()) throw new AlreadyExistsException(path);
        DirectoryNode parent = resolveParentDir(parts.subList(0, parts.size() - 1));
        String name = parts.get(parts.size() - 1);
        if (parent.children.containsKey(name)) throw new AlreadyExistsException(path);
        parent.children.put(name, new DirectoryNode(name));
    }

    public synchronized void createFile(String path, String content) {
        List<String> parts = splitPath(path, cwd);
        if (parts.isEmpty()) throw new IsDirectoryException(path);
        DirectoryNode parent = resolveParentDir(parts.subList(0, parts.size() - 1));
        String name = parts.get(parts.size() - 1);
        if (parent.children.containsKey(name)) throw new AlreadyExistsException(path);
        parent.children.put(name, new FileNode(name, content));
    }

    public void writeFile(String path, String content, boolean append) {
        List<String> parts = splitPath(path, cwd);
        Node node = resolve(parts);
        if (!(node instanceof FileNode)) throw new IsDirectoryException(path);
        ((FileNode) node).write(content, append);
    }

    public String readFile(String path) {
        List<String> parts = splitPath(path, cwd);
        Node node = resolve(parts);
        if (!(node instanceof FileNode)) throw new IsDirectoryException(path);
        return ((FileNode) node).read();
    }

    public synchronized List<Entry> ls(String path) {
        List<String> parts = splitPath(path, cwd);
        Node node = resolve(parts);
        if (!(node instanceof DirectoryNode)) throw new NotDirectoryException(path);
        DirectoryNode dir = (DirectoryNode) node;
        List<Entry> entries = new ArrayList<>();
        for (Node child : dir.children.values()) {
            entries.add(new Entry(child.getName(), child.isDir()));
        }
        entries.sort(Comparator.comparing(Entry::getName));
        return entries;
    }

    public synchronized void rm(String path, boolean recursive) {
        List<String> parts = splitPath(path, cwd);
        if (parts.isEmpty()) throw new IllegalArgumentException("cannot remove root");
        DirectoryNode parent = resolveParentDir(parts.subList(0, parts.size() - 1));
        String name = parts.get(parts.size() - 1);
        Node node = parent.children.get(name);
        if (node == null) throw new NotFoundException(path);
        if (node instanceof DirectoryNode) {
            DirectoryNode dir = (DirectoryNode) node;
            if (!dir.children.isEmpty() && !recursive) throw new DirectoryNotEmptyException(path);
        }
        parent.children.remove(name);
    }

    public synchronized void cd(String path) {
        List<String> parts = splitPath(path, cwd);
        Node node = resolve(parts);
        if (!node.isDir()) throw new NotDirectoryException(path);
        cwd = parts;
    }

    public synchronized String pwd() {
        if (cwd.isEmpty()) return "/";
        return "/" + String.join("/", cwd);
    }
}
