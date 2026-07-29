import java.util.Set;
import java.util.TreeSet;

/**
 * A single elevator cabin. Tracks its own current floor, motion state,
 * direction and the set of floors it still needs to visit. All mutating
 * methods are synchronized so hall-call dispatch and cabin requests can be
 * issued concurrently without corrupting the target set or state.
 */
public class Car {
    private final int id;
    private final int numFloors;
    private int currentFloor;
    private CarState state = CarState.IDLE;
    private Direction direction = Direction.IDLE;
    private final Set<Integer> targets = new TreeSet<>();

    public Car(int id, int numFloors, int startFloor) {
        this.id = id;
        this.numFloors = numFloors;
        this.currentFloor = startFloor;
    }

    public int getId() { return id; }
    public int getNumFloors() { return numFloors; }

    public synchronized int getCurrentFloor() { return currentFloor; }
    public synchronized CarState getState() { return state; }
    public synchronized Direction getDirection() { return direction; }
    public synchronized int pendingCount() { return targets.size(); }

    /** Registers a floor the car must visit; derives a direction if idle. */
    public synchronized void addTarget(int floor) {
        targets.add(floor);
        if (state == CarState.IDLE) {
            setDirectionTowards(floor);
        }
    }

    private void setDirectionTowards(int floor) {
        if (floor > currentFloor) {
            direction = Direction.UP;
            state = CarState.MOVING_UP;
        } else if (floor < currentFloor) {
            direction = Direction.DOWN;
            state = CarState.MOVING_DOWN;
        } else {
            // The target is the floor we're already sitting on: service it now,
            // otherwise it would never be removed and pickNextDirection would
            // keep re-selecting it as "nearest" forever, starving every other
            // pending target on this car.
            targets.remove(floor);
            state = CarState.DOOR_OPEN;
        }
    }

    /**
     * Advances simulated time by one tick: a car with doors open closes them
     * and picks a new direction (or goes idle); otherwise it moves one floor
     * towards its next target, servicing any target it lands on.
     */
    public synchronized void step() {
        if (state == CarState.DOOR_OPEN) {
            pickNextDirection();
            return;
        }
        if (targets.isEmpty()) {
            state = CarState.IDLE;
            direction = Direction.IDLE;
            return;
        }

        switch (direction) {
            case UP -> currentFloor++;
            case DOWN -> currentFloor--;
            default -> {
                pickNextDirection();
                return;
            }
        }

        if (targets.remove(currentFloor)) {
            state = CarState.DOOR_OPEN;
        }
    }

    /** Simplified SCAN: head towards the nearest remaining target. */
    private void pickNextDirection() {
        if (targets.isEmpty()) {
            state = CarState.IDLE;
            direction = Direction.IDLE;
            return;
        }
        int nearest = targets.iterator().next();
        int best = Math.abs(nearest - currentFloor);
        for (int f : targets) {
            int d = Math.abs(f - currentFloor);
            if (d < best) {
                best = d;
                nearest = f;
            }
        }
        setDirectionTowards(nearest);
    }
}
