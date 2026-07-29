import java.util.List;

/**
 * Prefers an idle car nearest to the call floor; failing that, a car already
 * moving in the requested direction that hasn't yet passed the call floor
 * (so it can pick it up en route); failing that, the overall nearest car
 * regardless of direction compatibility.
 */
public class NearestCarStrategy implements DispatchStrategy {
    @Override
    public Car selectCar(List<Car> cars, int floor, Direction direction) {
        Car best = null;
        int bestCost = Integer.MAX_VALUE;

        for (Car c : cars) {
            int f = c.getCurrentFloor();
            CarState state = c.getState();
            Direction d = c.getDirection();
            int cost = Math.abs(f - floor);

            boolean compatible = state == CarState.IDLE
                    || (d == direction && direction == Direction.UP && f <= floor)
                    || (d == direction && direction == Direction.DOWN && f >= floor);

            if (!compatible) {
                cost += c.getNumFloors();
            }

            if (best == null || cost < bestCost) {
                best = c;
                bestCost = cost;
            }
        }
        return best;
    }
}
