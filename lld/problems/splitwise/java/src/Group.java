import java.util.List;

public class Group {
    private final String id;
    private final String name;
    private final List<User> members;

    public Group(String id, String name, List<User> members) {
        this.id = id;
        this.name = name;
        this.members = members;
    }

    public String getId() { return id; }
    public String getName() { return name; }
    public List<User> getMembers() { return members; }
}
