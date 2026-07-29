import java.util.ArrayList;
import java.util.List;

public class Main {
    public static void main(String[] args) throws Exception {
        List<Driver> drivers = new ArrayList<>();
        drivers.add(new Driver("D1", "Alice", new Location(0, 0), true));
        drivers.add(new Driver("D2", "Bob", new Location(10, 10), true));
        RideSharingSystem system = new RideSharingSystem(drivers, new NearestAvailableDriverStrategy(),
                new DistanceBasedPricing(2, 1.5));

        Rider rider = new Rider("R1", "Riya");
        Trip trip = system.requestRide(rider, new Location(0, 0), new Location(3, 4));
        System.out.println("Requested trip " + trip.id + " status=" + trip.status);

        Driver driver = system.matchDriver(trip.id);
        System.out.println("Matched driver " + driver.id + " status=" + trip.status);

        system.startTrip(trip.id);
        System.out.println("Started trip, status=" + trip.status);

        double fare = system.completeTrip(trip.id);
        System.out.println("Completed trip, fare=" + fare + " status=" + trip.status);

        RideSharingTest.runAll();
    }
}
