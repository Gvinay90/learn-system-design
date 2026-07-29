import java.time.Instant;
import java.util.List;

public class Main {
    public static void main(String[] args) {
        Floor floor1 = new Floor(1, List.of(
                new ParkingSpot("F1-M1", VehicleType.MOTORCYCLE),
                new ParkingSpot("F1-C1", VehicleType.COMPACT),
                new ParkingSpot("F1-L1", VehicleType.LARGE)
        ));
        ParkingLot lot = new ParkingLot(List.of(floor1), new HourlyTieredPricing(5, 2));

        Vehicle car = new Vehicle("KA-01", VehicleType.COMPACT);
        Ticket ticket = lot.parkVehicle(car);
        System.out.println("Parked " + car.getLicensePlate() + " at spot " + ticket.getSpot().getId());

        double fee = lot.unparkVehicle(ticket.getId(), ticket.getEntryTime().plusSeconds(3 * 3600));
        System.out.println("Fee for 3 hours: " + fee);

        ParkingLotTest.runAll();
    }
}
