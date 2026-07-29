import java.time.Instant;

public class Trip {
    public final String id;
    public final Rider rider;
    public Driver driver;
    public final Location pickup;
    public final Location dropoff;
    public TripStatus status;
    public final Instant requestedAt;
    public double fare;

    public Trip(String id, Rider rider, Location pickup, Location dropoff) {
        this.id = id;
        this.rider = rider;
        this.pickup = pickup;
        this.dropoff = dropoff;
        this.status = TripStatus.REQUESTED;
        this.requestedAt = Instant.now();
        this.fare = 0;
    }
}
