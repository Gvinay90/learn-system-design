import java.util.ArrayList;
import java.util.List;

public class Main {
    public static void main(String[] args) {
        List<Car> cars = new ArrayList<>();
        cars.add(new Car(0, 10, 0));
        cars.add(new Car(1, 10, 9));
        ElevatorSystem system = new ElevatorSystem(10, cars, new NearestCarStrategy());

        Car assigned = system.hallCall(5, Direction.UP);
        System.out.println("Hall call at floor 5 (up) dispatched to car " + assigned.getId());

        for (int i = 0; i < 20 && anyBusy(system); i++) {
            system.step();
        }
        System.out.println("Car " + assigned.getId() + " now at floor " + assigned.getCurrentFloor()
                + ", state=" + assigned.getState());

        ElevatorTest.runAll();
    }

    private static boolean anyBusy(ElevatorSystem system) {
        for (Car c : system.getCars()) {
            if (c.getState() != CarState.IDLE || c.pendingCount() > 0) {
                return true;
            }
        }
        return false;
    }
}
