import java.util.List;

/**
 * Coordinates a fleet of cars serving a building, dispatching hall calls via
 * a pluggable DispatchStrategy.
 */
public class ElevatorSystem {
    private final List<Car> cars;
    private final int numFloors;
    private final DispatchStrategy strategy;
    private final Object dispatchLock = new Object();

    public ElevatorSystem(int numFloors, List<Car> cars, DispatchStrategy strategy) {
        this.numFloors = numFloors;
        this.cars = cars;
        this.strategy = strategy;
    }

    public List<Car> getCars() { return cars; }
    public int getNumFloors() { return numFloors; }

    /**
     * Handles an external hall call, dispatching it to the strategy's chosen
     * car. The select-then-assign sequence is guarded by dispatchLock so two
     * concurrent hall calls can never both "win" a decision based on a stale
     * view of car state (analogous to Floor.findAndAssign in parking-lot).
     */
    public Car hallCall(int floor, Direction direction) {
        synchronized (dispatchLock) {
            Car car = strategy.selectCar(cars, floor, direction);
            car.addTarget(floor);
            return car;
        }
    }

    /** Handles an internal request made from inside a specific car. */
    public void cabinRequest(Car car, int destinationFloor) {
        car.addTarget(destinationFloor);
    }

    /** Advances every car in the fleet by one simulated tick. */
    public void step() {
        for (Car c : cars) {
            c.step();
        }
    }
}
