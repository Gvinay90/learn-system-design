import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CountDownLatch;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out ElevatorTest` directly.
 */
public class ElevatorTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testHallCallDispatchedAndServiced();
        testCallToCurrentFloor();
        testAllCarsBusyStillAssigns();
        testNearestCarDispatchCorrectness();
        testNearestCarPrefersEnRouteOverIdle();
        testConcurrentHallCalls();
        System.out.println("All ElevatorTest cases passed.");
    }

    private static ElevatorSystem newSystem(int numFloors, int... startFloors) {
        List<Car> cars = new ArrayList<>();
        for (int i = 0; i < startFloors.length; i++) {
            cars.add(new Car(i, numFloors, startFloors[i]));
        }
        return new ElevatorSystem(numFloors, cars, new NearestCarStrategy());
    }

    private static void runUntilIdle(ElevatorSystem sys, int maxSteps) {
        for (int i = 0; i < maxSteps; i++) {
            sys.step();
            boolean allIdle = true;
            for (Car c : sys.getCars()) {
                if (c.getState() != CarState.IDLE || c.pendingCount() > 0) {
                    allIdle = false;
                }
            }
            if (allIdle) {
                return;
            }
        }
    }

    private static void testHallCallDispatchedAndServiced() {
        ElevatorSystem sys = newSystem(10, 0);
        Car car = sys.hallCall(5, Direction.UP);
        assertEquals(0, car.getId(), "expected only car to be dispatched");

        runUntilIdle(sys, 50);

        assertEquals(5, car.getCurrentFloor(), "expected car to reach floor 5");
        assertEquals(0, car.pendingCount(), "expected no pending targets");
        assertTrue(car.getState() == CarState.DOOR_OPEN || car.getState() == CarState.IDLE,
                "expected car to have serviced its stop");
    }

    // Covers the edge case where the hall call floor is the floor the only
    // car is already parked at: it should open its doors immediately.
    private static void testCallToCurrentFloor() {
        ElevatorSystem sys = newSystem(10, 3);
        Car car = sys.hallCall(3, Direction.UP);

        assertEquals(3, car.getCurrentFloor(), "expected car to remain at floor 3");
        assertEquals(CarState.DOOR_OPEN, car.getState(), "expected door to open immediately at current floor");
    }

    // Covers the edge case where every car is already moving: the strategy
    // must still pick some car rather than dropping the call.
    private static void testAllCarsBusyStillAssigns() {
        ElevatorSystem sys = newSystem(20, 0, 10);
        sys.hallCall(19, Direction.UP);
        sys.hallCall(1, Direction.DOWN);

        Car car = sys.hallCall(15, Direction.UP);
        assertTrue(car != null, "expected a car to be assigned even though both cars are busy");
        assertTrue(car.pendingCount() > 0, "expected the assigned car to have the new target registered");
    }

    private static void testNearestCarDispatchCorrectness() {
        ElevatorSystem sys = newSystem(20, 0, 18);
        Car car = sys.hallCall(15, Direction.UP);
        assertEquals(1, car.getId(), "expected nearest idle car (id 1) to be dispatched");
    }

    private static void testNearestCarPrefersEnRouteOverIdle() {
        ElevatorSystem sys = newSystem(20, 0, 0);
        sys.hallCall(10, Direction.UP);
        for (int i = 0; i < 4; i++) {
            sys.step();
        }
        Car car0 = sys.getCars().get(0);
        assertTrue(car0.getState() == CarState.MOVING_UP || car0.getState() == CarState.DOOR_OPEN,
                "expected car 0 to be moving up");

        Car car = sys.hallCall(6, Direction.UP);
        assertEquals(0, car.getId(), "expected en-route car 0 to be preferred over idle car 1");
    }

    private static void testConcurrentHallCalls() {
        final int numFloors = 50;
        ElevatorSystem sys = newSystem(numFloors, 0, 10, 25, 40);

        final int numCalls = 200;
        Car[] assigned = new Car[numCalls];
        CountDownLatch latch = new CountDownLatch(numCalls);

        for (int n = 0; n < numCalls; n++) {
            final int idx = n;
            Thread t = new Thread(() -> {
                int floor = (idx * 7) % numFloors;
                Direction dir = idx % 2 == 0 ? Direction.DOWN : Direction.UP;
                assigned[idx] = sys.hallCall(floor, dir);
                latch.countDown();
            });
            t.start();
        }
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        for (int i = 0; i < numCalls; i++) {
            assertTrue(assigned[i] != null, "call " + i + " was never assigned a car");
        }

        runUntilIdle(sys, numFloors * numCalls * 2);

        for (Car c : sys.getCars()) {
            assertEquals(0, c.pendingCount(), "car " + c.getId() + " still has pending targets after running to completion");
            assertEquals(CarState.IDLE, c.getState(), "car " + c.getId() + " expected to be idle after servicing all calls");
            assertTrue(c.getCurrentFloor() >= 0 && c.getCurrentFloor() < numFloors,
                    "car " + c.getId() + " ended up at out-of-range floor (state corruption)");
        }
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
}
