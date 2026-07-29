import java.time.Instant;
import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out ParkingLotTest` directly.
 */
public class ParkingLotTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testParkAndUnpark();
        testLotFullForType();
        testUnparkUnknownTicket();
        testConcurrentParkingSameSpot();
        System.out.println("All ParkingLotTest cases passed.");
    }

    private static ParkingLot newTestLot() {
        Floor floor1 = new Floor(1, List.of(
                new ParkingSpot("F1-M1", VehicleType.MOTORCYCLE),
                new ParkingSpot("F1-C1", VehicleType.COMPACT),
                new ParkingSpot("F1-L1", VehicleType.LARGE)
        ));
        return new ParkingLot(List.of(floor1), new HourlyTieredPricing(5, 2));
    }

    private static void testParkAndUnpark() {
        ParkingLot lot = newTestLot();
        Vehicle v = new Vehicle("KA-01", VehicleType.COMPACT);
        Ticket ticket = lot.parkVehicle(v);
        assertEquals("F1-C1", ticket.getSpot().getId(), "spot assignment");

        double fee = lot.unparkVehicle(ticket.getId(), ticket.getEntryTime().plusSeconds(1800));
        assertEquals(5.0, fee, "flat first-hour fee");
    }

    private static void testLotFullForType() {
        ParkingLot lot = newTestLot();
        lot.parkVehicle(new Vehicle("A1", VehicleType.COMPACT));
        try {
            lot.parkVehicle(new Vehicle("A2", VehicleType.COMPACT));
            throw new AssertionError("expected LotFullException");
        } catch (ParkingLot.LotFullException e) {
            // expected
        }
    }

    private static void testUnparkUnknownTicket() {
        ParkingLot lot = newTestLot();
        try {
            lot.unparkVehicle("bogus", Instant.now());
            throw new AssertionError("expected TicketNotFoundException");
        } catch (ParkingLot.TicketNotFoundException e) {
            // expected
        }
    }

    private static void testConcurrentParkingSameSpot() {
        ParkingLot lot = newTestLot();
        AtomicInteger successCount = new AtomicInteger(0);
        CountDownLatch latch = new CountDownLatch(2);

        Runnable task = () -> {
            try {
                lot.parkVehicle(new Vehicle("race", VehicleType.MOTORCYCLE));
                successCount.incrementAndGet();
            } catch (ParkingLot.LotFullException ignored) {
            } finally {
                latch.countDown();
            }
        };
        new Thread(task).start();
        new Thread(task).start();
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        assertEquals(1, successCount.get(), "exactly one thread should win the single motorcycle spot");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
