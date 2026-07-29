import java.util.List;

public class Floor {
    private final int level;
    private final List<ParkingSpot> spots;

    public Floor(int level, List<ParkingSpot> spots) {
        this.level = level;
        this.spots = spots;
    }

    public int getLevel() { return level; }

    /** Finds the first spot that fits the vehicle and atomically assigns it, or returns null. */
    public synchronized ParkingSpot findAndAssign(Vehicle v) {
        for (ParkingSpot s : spots) {
            if (s.tryAssign(v)) {
                return s;
            }
        }
        return null;
    }

    public void free(String spotId) {
        for (ParkingSpot s : spots) {
            if (s.getId().equals(spotId)) {
                s.free();
                return;
            }
        }
    }
}
