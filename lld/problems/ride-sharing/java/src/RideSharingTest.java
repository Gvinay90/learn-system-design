import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out RideSharingTest` directly.
 */
public class RideSharingTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testHappyPathLifecycle();
        testInvalidTransitionRejected();
        testNoDriverAvailable();
        testConcurrentMatching();
        System.out.println("All RideSharingTest cases passed.");
    }

    private static RideSharingSystem newTestSystem() {
        List<Driver> drivers = new ArrayList<>();
        drivers.add(new Driver("D1", "Alice", new Location(0, 0), true));
        drivers.add(new Driver("D2", "Bob", new Location(10, 10), true));
        return new RideSharingSystem(drivers, new NearestAvailableDriverStrategy(),
                new DistanceBasedPricing(2, 1.5));
    }

    private static void testHappyPathLifecycle() {
        try {
            RideSharingSystem sys = newTestSystem();
            Rider rider = new Rider("R1", "Riya");

            Trip trip = sys.requestRide(rider, new Location(0, 0), new Location(3, 4));
            assertEquals(TripStatus.REQUESTED, trip.status, "expected REQUESTED");

            Driver driver = sys.matchDriver(trip.id);
            assertEquals("D1", driver.id, "expected nearest driver D1");
            assertEquals(TripStatus.ACCEPTED, trip.status, "expected ACCEPTED");

            sys.startTrip(trip.id);
            assertEquals(TripStatus.IN_PROGRESS, trip.status, "expected IN_PROGRESS");

            double fare = sys.completeTrip(trip.id);
            // base 2 + 5 (distance from (0,0) to (3,4)) * 1.5 = 9.5
            assertEquals(9.5, fare, "expected fare 9.5");
            assertEquals(TripStatus.COMPLETED, trip.status, "expected COMPLETED");
            assertTrue(driver.isAvailable(), "expected driver freed after trip completion");
        } catch (Exception e) {
            throw new AssertionError("unexpected exception: " + e, e);
        }
    }

    private static void testInvalidTransitionRejected() {
        RideSharingSystem sys = newTestSystem();
        Rider rider = new Rider("R1", "Riya");
        Trip trip = sys.requestRide(rider, new Location(0, 0), new Location(1, 1));

        try {
            sys.completeTrip(trip.id);
            throw new AssertionError("expected InvalidTransitionException completing a REQUESTED trip");
        } catch (InvalidTransitionException expected) {
            // expected
        } catch (Exception e) {
            throw new AssertionError("unexpected exception: " + e, e);
        }

        try {
            sys.matchDriver(trip.id);
            sys.startTrip(trip.id);
            sys.completeTrip(trip.id);
        } catch (Exception e) {
            throw new AssertionError("unexpected exception: " + e, e);
        }

        try {
            sys.cancelTrip(trip.id);
            throw new AssertionError("expected InvalidTransitionException cancelling a COMPLETED trip");
        } catch (InvalidTransitionException expected) {
            // expected
        } catch (Exception e) {
            throw new AssertionError("unexpected exception: " + e, e);
        }
    }

    private static void testNoDriverAvailable() {
        RideSharingSystem sys = newTestSystem();
        Rider rider = new Rider("R1", "Riya");

        Trip t1 = sys.requestRide(rider, new Location(0, 0), new Location(1, 1));
        Trip t2 = sys.requestRide(rider, new Location(0, 0), new Location(1, 1));
        Trip t3 = sys.requestRide(rider, new Location(0, 0), new Location(1, 1));

        try {
            sys.matchDriver(t1.id);
            sys.matchDriver(t2.id);
        } catch (Exception e) {
            throw new AssertionError("unexpected exception: " + e, e);
        }

        try {
            sys.matchDriver(t3.id);
            throw new AssertionError("expected NoDriverAvailableException");
        } catch (NoDriverAvailableException expected) {
            // expected
        } catch (Exception e) {
            throw new AssertionError("unexpected exception: " + e, e);
        }
    }

    // Asserts two threads racing for the same single available driver never
    // both succeed — the lock in matchDriver must serialize them.
    private static void testConcurrentMatching() {
        List<Driver> drivers = new ArrayList<>();
        drivers.add(new Driver("D1", "Alice", new Location(0, 0), true));
        RideSharingSystem sys = new RideSharingSystem(drivers, new NearestAvailableDriverStrategy(),
                new DistanceBasedPricing(2, 1.5));
        Rider rider = new Rider("R1", "Riya");

        Trip trip1 = sys.requestRide(rider, new Location(0, 0), new Location(1, 1));
        Trip trip2 = sys.requestRide(rider, new Location(0, 0), new Location(1, 1));

        String[] tripIds = new String[] { trip1.id, trip2.id };
        AtomicInteger successCount = new AtomicInteger(0);
        CountDownLatch latch = new CountDownLatch(tripIds.length);

        for (String tripId : tripIds) {
            Thread t = new Thread(() -> {
                try {
                    sys.matchDriver(tripId);
                    successCount.incrementAndGet();
                } catch (Exception ignored) {
                    // expected for the loser
                } finally {
                    latch.countDown();
                }
            });
            t.start();
        }

        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        assertEquals(1, successCount.get(), "expected exactly 1 success for single available driver");
    }

    private static void assertTrue(boolean condition, String label) {
        if (!condition) {
            throw new AssertionError(label);
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }

    private static void assertEquals(double expected, double actual, String label) {
        if (Math.abs(expected - actual) > 1e-9) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
