import java.util.List;
import java.util.Optional;

public class Restaurant {
    private final String id;
    private final String name;
    private final Location location;
    private final List<MenuItem> menu;
    private final boolean open;

    public Restaurant(String id, String name, Location location, List<MenuItem> menu, boolean open) {
        this.id = id;
        this.name = name;
        this.location = location;
        this.menu = menu;
        this.open = open;
    }

    public String getId() { return id; }
    public String getName() { return name; }
    public Location getLocation() { return location; }
    public boolean isOpen() { return open; }

    public Optional<MenuItem> findItem(String itemId) {
        return menu.stream().filter(i -> i.getId().equals(itemId)).findFirst();
    }
}
