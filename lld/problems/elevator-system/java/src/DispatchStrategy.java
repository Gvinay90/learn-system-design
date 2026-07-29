import java.util.List;

public interface DispatchStrategy {
    Car selectCar(List<Car> cars, int floor, Direction direction);
}
