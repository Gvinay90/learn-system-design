import java.util.List;

/** Picks a driver for a trip from the available pool. */
public interface DriverMatchingStrategy {
    Driver match(Location pickup, List<Driver> drivers);
}
