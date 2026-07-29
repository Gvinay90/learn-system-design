import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Document {
    private String title;
    private Metadata meta;
    private final List<String> sections;
    private final Map<String, String> props;

    public Document(String title, Metadata meta, List<String> sections, Map<String, String> props) {
        this.title = title;
        this.meta = meta;
        this.sections = new ArrayList<>(sections);
        this.props = new HashMap<>(props);
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public Metadata getMeta() {
        return meta;
    }

    public List<String> getSections() {
        return sections;
    }

    public Map<String, String> getProps() {
        return props;
    }

    // Deep copy: nested Metadata, the sections list, and the props map are
    // all duplicated rather than shared with the original.
    public Document deepClone() {
        return new Document(title, meta.deepClone(), new ArrayList<>(sections), new HashMap<>(props));
    }
}
