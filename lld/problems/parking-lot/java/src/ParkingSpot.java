public class ParkingSpot {
    private final String id;
    private final VehicleType type;
    private boolean occupied;
    private Vehicle vehicle;

    public ParkingSpot(String id, VehicleType type) {
        this.id = id;
        this.type = type;
    }

    public String getId() { return id; }
    public VehicleType getType() { return type; }

    public synchronized boolean tryAssign(Vehicle v) {
        if (occupied || type != v.getType()) {
            return false;
        }
        occupied = true;
        vehicle = v;
        return true;
    }

    public synchronized void free() {
        occupied = false;
        vehicle = null;
    }
}
