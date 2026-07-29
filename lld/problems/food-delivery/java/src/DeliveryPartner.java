public class DeliveryPartner {
    private final String id;
    private final String name;
    private final Location location;
    private volatile boolean available;

    public DeliveryPartner(String id, String name, Location location, boolean available) {
        this.id = id;
        this.name = name;
        this.location = location;
        this.available = available;
    }

    public String getId() { return id; }
    public String getName() { return name; }
    public Location getLocation() { return location; }
    public boolean isAvailable() { return available; }
    public void setAvailable(boolean available) { this.available = available; }
}
