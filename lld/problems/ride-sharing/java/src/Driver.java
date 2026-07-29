public class Driver {
    public final String id;
    public final String name;
    public final Location location;
    private boolean available;

    public Driver(String id, String name, Location location, boolean available) {
        this.id = id;
        this.name = name;
        this.location = location;
        this.available = available;
    }

    public boolean isAvailable() {
        return available;
    }

    public void setAvailable(boolean available) {
        this.available = available;
    }
}
