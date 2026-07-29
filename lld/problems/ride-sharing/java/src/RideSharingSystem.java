import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Facade coordinating driver pool, trip map, matching, and pricing.
 *
 * All lifecycle methods synchronize on a single system lock so that, in
 * particular, the find-and-mark-unavailable step in matchDriver is atomic:
 * two concurrent calls can never both be assigned the same single driver.
 */
public class RideSharingSystem {
    private final DriverMatchingStrategy matching;
    private final PricingStrategy pricing;
    private final Object lock = new Object();
    private final List<Driver> drivers;
    private final Map<String, Trip> trips = new HashMap<>();
    private int seq = 0;

    public RideSharingSystem(List<Driver> drivers, DriverMatchingStrategy matching, PricingStrategy pricing) {
        this.drivers = drivers;
        this.matching = matching;
        this.pricing = pricing;
    }

    public Trip requestRide(Rider rider, Location pickup, Location dropoff) {
        synchronized (lock) {
            seq++;
            Trip trip = new Trip("TR-" + seq, rider, pickup, dropoff);
            trips.put(trip.id, trip);
            return trip;
        }
    }

    /**
     * Finds the nearest available driver and assigns it atomically: the
     * whole find-and-mark-unavailable step is guarded by the system lock so
     * two concurrent calls can never both be assigned the same single driver.
     */
    public Driver matchDriver(String tripId) throws TripNotFoundException, InvalidTransitionException, NoDriverAvailableException {
        synchronized (lock) {
            Trip trip = trips.get(tripId);
            if (trip == null) {
                throw new TripNotFoundException();
            }
            if (trip.status != TripStatus.REQUESTED) {
                throw new InvalidTransitionException();
            }

            Driver driver = matching.match(trip.pickup, drivers);
            if (driver == null) {
                throw new NoDriverAvailableException();
            }

            driver.setAvailable(false);
            trip.driver = driver;
            trip.status = TripStatus.ACCEPTED;
            return driver;
        }
    }

    public void startTrip(String tripId) throws TripNotFoundException, InvalidTransitionException {
        synchronized (lock) {
            Trip trip = trips.get(tripId);
            if (trip == null) {
                throw new TripNotFoundException();
            }
            if (trip.status != TripStatus.ACCEPTED) {
                throw new InvalidTransitionException();
            }
            trip.status = TripStatus.IN_PROGRESS;
        }
    }

    public double completeTrip(String tripId) throws TripNotFoundException, InvalidTransitionException {
        synchronized (lock) {
            Trip trip = trips.get(tripId);
            if (trip == null) {
                throw new TripNotFoundException();
            }
            if (trip.status != TripStatus.IN_PROGRESS) {
                throw new InvalidTransitionException();
            }
            trip.status = TripStatus.COMPLETED;
            trip.fare = pricing.calculateFare(trip);
            if (trip.driver != null) {
                trip.driver.setAvailable(true);
            }
            return trip.fare;
        }
    }

    public void cancelTrip(String tripId) throws TripNotFoundException, InvalidTransitionException {
        synchronized (lock) {
            Trip trip = trips.get(tripId);
            if (trip == null) {
                throw new TripNotFoundException();
            }
            if (trip.status != TripStatus.REQUESTED && trip.status != TripStatus.ACCEPTED) {
                throw new InvalidTransitionException();
            }
            trip.status = TripStatus.CANCELLED;
            if (trip.driver != null) {
                trip.driver.setAvailable(true);
            }
        }
    }
}
