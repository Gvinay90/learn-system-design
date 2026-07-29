import java.util.ArrayList;
import java.util.List;

public class Metadata {
    private String author;
    private final List<String> tags;

    public Metadata(String author, List<String> tags) {
        this.author = author;
        this.tags = new ArrayList<>(tags);
    }

    public String getAuthor() {
        return author;
    }

    public void setAuthor(String author) {
        this.author = author;
    }

    public List<String> getTags() {
        return tags;
    }

    // Deep-copies the tag list so mutating the clone's tags never affects
    // the original.
    Metadata deepClone() {
        return new Metadata(author, new ArrayList<>(tags));
    }
}
