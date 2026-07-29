import java.time.Instant;

public class Ticket {
    private final String id;
    private final Vehicle vehicle;
    private final ParkingSpot spot;
    private final int floorLevel;
    private final Instant entryTime;

    public Ticket(String id, Vehicle vehicle, ParkingSpot spot, int floorLevel, Instant entryTime) {
        this.id = id;
        this.vehicle = vehicle;
        this.spot = spot;
        this.floorLevel = floorLevel;
        this.entryTime = entryTime;
    }

    public String getId() { return id; }
    public ParkingSpot getSpot() { return spot; }
    public int getFloorLevel() { return floorLevel; }
    public Instant getEntryTime() { return entryTime; }
}
