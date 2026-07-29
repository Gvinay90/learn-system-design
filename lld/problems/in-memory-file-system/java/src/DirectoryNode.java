import java.util.LinkedHashMap;
import java.util.Map;

public class DirectoryNode implements Node {
    private final String name;
    final Map<String, Node> children = new LinkedHashMap<>();

    public DirectoryNode(String name) {
        this.name = name;
    }

    @Override
    public String getName() { return name; }

    @Override
    public boolean isDir() { return true; }
}
