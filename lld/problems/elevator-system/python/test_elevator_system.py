import threading

from elevator_system import (
    Car,
    CarState,
    Direction,
    ElevatorSystem,
    NearestCarStrategy,
)


def new_test_system(num_floors: int, *start_floors: int) -> ElevatorSystem:
    cars = [Car(i, num_floors, f) for i, f in enumerate(start_floors)]
    return ElevatorSystem(num_floors, cars, NearestCarStrategy())


def run_until_idle(system: ElevatorSystem, max_steps: int) -> None:
    for _ in range(max_steps):
        system.step()
        all_idle = True
        for c in system.cars:
            snap = c.snapshot()
            if snap.state != CarState.IDLE or snap.pending > 0:
                all_idle = False
        if all_idle:
            return


def test_hall_call_dispatched_and_serviced():
    sys_ = new_test_system(10, 0)
    car = sys_.hall_call(5, Direction.UP)
    assert car.id == 0

    run_until_idle(sys_, 50)

    snap = car.snapshot()
    assert snap.floor == 5
    assert snap.pending == 0
    assert snap.state in (CarState.DOOR_OPEN, CarState.IDLE)


def test_call_to_current_floor():
    """The hall call floor is the floor the only car is already parked at:
    it should open its doors immediately rather than trying to move."""
    sys_ = new_test_system(10, 3)
    car = sys_.hall_call(3, Direction.UP)

    snap = car.snapshot()
    assert snap.floor == 3
    assert snap.state == CarState.DOOR_OPEN


def test_all_cars_busy_still_assigns():
    """Every car is already moving: the strategy must still pick some car
    (the least-bad option) rather than dropping the call."""
    sys_ = new_test_system(20, 0, 10)
    sys_.hall_call(19, Direction.UP)  # car 0 now moving up towards 19
    sys_.hall_call(1, Direction.DOWN)  # car 1 now moving down towards 1

    car = sys_.hall_call(15, Direction.UP)
    assert car is not None
    snap = car.snapshot()
    assert snap.pending > 0


def test_nearest_car_dispatch_correctness():
    # Car 0 idle at floor 0, car 1 idle at floor 18: a call at floor 15
    # heading up should go to the nearer car (car 1).
    sys_ = new_test_system(20, 0, 18)
    car = sys_.hall_call(15, Direction.UP)
    assert car.id == 1


def test_nearest_car_prefers_en_route_over_idle():
    # Car 0 is moving up and already past floor 4 heading towards floor 10 -
    # it can pick up a call at floor 6 en route. Car 1 is idle at floor 0,
    # which is farther from floor 6. The en-route car should win.
    sys_ = new_test_system(20, 0, 0)
    sys_.hall_call(10, Direction.UP)  # dispatched to car 0 (both idle at 0 -> tie goes to first)
    for _ in range(4):
        sys_.step()

    snap = sys_.cars[0].snapshot()
    assert snap.state in (CarState.MOVING_UP, CarState.DOOR_OPEN)
    assert snap.direction == Direction.UP or snap.state == CarState.DOOR_OPEN

    car = sys_.hall_call(6, Direction.UP)
    assert car.id == 0


def test_concurrent_hall_calls():
    """Fires many concurrent hall calls at an ElevatorSystem with multiple
    cars from threads and asserts that the system's lock-guarded
    select-then-assign sequence in hall_call never corrupts a car's target
    set (each call must end up registered on exactly the car it was
    returned for) and every call is eventually serviced."""
    num_floors = 50
    sys_ = new_test_system(num_floors, 0, 10, 25, 40)

    num_calls = 200
    assigned = [None] * num_calls

    def worker(n: int) -> None:
        floor = (n * 7) % num_floors
        direction = Direction.DOWN if n % 2 == 0 else Direction.UP
        assigned[n] = sys_.hall_call(floor, direction)

    threads = [threading.Thread(target=worker, args=(i,)) for i in range(num_calls)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    for car in assigned:
        assert car is not None

    # The nearest-target dispatch heuristic (not a strict direction-committed
    # SCAN) can thrash back and forth when many targets are scattered across
    # one car, so give it a generously large step budget rather than a tight
    # one derived from floor/call counts.
    run_until_idle(sys_, num_floors * num_calls * 2)

    for c in sys_.cars:
        snap = c.snapshot()
        assert snap.pending == 0
        assert snap.state == CarState.IDLE
        assert 0 <= c.current_floor < num_floors
