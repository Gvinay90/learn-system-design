import java.util.List;

/** Picks the closest available driver to the pickup point. */
public class NearestAvailableDriverStrategy implements DriverMatchingStrategy {
    @Override
    public Driver match(Location pickup, List<Driver> drivers) {
        Driver best = null;
        double bestDist = Double.MAX_VALUE;
        for (Driver d : drivers) {
            if (!d.isAvailable()) {
                continue;
            }
            double dist = pickup.distanceTo(d.location);
            if (dist < bestDist) {
                bestDist = dist;
                best = d;
            }
        }
        return best;
    }
}
