"""Elevator System LLD — Python reference implementation.

N floors, M elevator cars, hall calls + cabin requests, and a pluggable
dispatch/scheduling strategy that picks the best car for a hall call.
See ../go/elevatorsystem.go for the original design writeup.
"""
from __future__ import annotations

import threading
from enum import Enum, auto
from typing import Dict, List, Optional, Protocol


class Direction(Enum):
    IDLE = auto()
    UP = auto()
    DOWN = auto()


class CarState(Enum):
    IDLE = auto()
    MOVING_UP = auto()
    MOVING_DOWN = auto()
    DOOR_OPEN = auto()


class CarSnapshot:
    """Immutable view of a car's state at a point in time."""

    __slots__ = ("floor", "state", "direction", "pending")

    def __init__(self, floor: int, state: CarState, direction: Direction, pending: int):
        self.floor = floor
        self.state = state
        self.direction = direction
        self.pending = pending


class Car:
    """A single elevator cabin.

    Tracks its own current floor, motion state, direction and the set of
    floors it still needs to visit. A per-car lock guards all of this so
    that hall-call dispatch (from the system) and cabin requests (from
    inside the car) can be issued concurrently without corrupting the
    target set or state.
    """

    def __init__(self, car_id: int, num_floors: int, start_floor: int):
        self.id = car_id
        self.num_floors = num_floors
        self.current_floor = start_floor

        self._lock = threading.Lock()
        self._state = CarState.IDLE
        self._direction = Direction.IDLE
        self._targets: set[int] = set()

    def snapshot(self) -> CarSnapshot:
        with self._lock:
            return CarSnapshot(self.current_floor, self._state, self._direction, len(self._targets))

    def add_target(self, floor: int) -> None:
        """Registers a floor the car must visit (from a hall call assigned
        to this car, or a cabin request made from inside it). If the car
        is currently idle, its direction is derived from the new target.
        """
        with self._lock:
            self._targets.add(floor)
            if self._state == CarState.IDLE:
                self._set_direction_towards(floor)

    def _set_direction_towards(self, floor: int) -> None:
        if floor > self.current_floor:
            self._direction = Direction.UP
            self._state = CarState.MOVING_UP
        elif floor < self.current_floor:
            self._direction = Direction.DOWN
            self._state = CarState.MOVING_DOWN
        else:
            # The target is the floor we're already sitting on: service it
            # now, otherwise it would never be removed and
            # _pick_next_direction would keep re-selecting it as "nearest"
            # forever, starving every other pending target on this car.
            self._targets.discard(floor)
            self._state = CarState.DOOR_OPEN

    def step(self) -> None:
        """Advances simulated time by one tick: a car in DOOR_OPEN closes
        its doors and picks a new direction (or goes idle); otherwise it
        moves one floor towards its next target, servicing (removing) any
        target it lands on.
        """
        with self._lock:
            if self._state == CarState.DOOR_OPEN:
                self._pick_next_direction()
                return
            if not self._targets:
                self._state = CarState.IDLE
                self._direction = Direction.IDLE
                return

            if self._direction == Direction.UP:
                self.current_floor += 1
            elif self._direction == Direction.DOWN:
                self.current_floor -= 1
            else:
                self._pick_next_direction()
                return

            if self.current_floor in self._targets:
                self._targets.discard(self.current_floor)
                self._state = CarState.DOOR_OPEN

    def _pick_next_direction(self) -> None:
        """Chooses the direction of the nearest remaining target,
        simplified SCAN behaviour: continue in the current direction if a
        target still lies ahead, otherwise reverse, otherwise go idle.
        """
        if not self._targets:
            self._state = CarState.IDLE
            self._direction = Direction.IDLE
            return

        # Sort so that ties in distance are broken deterministically in
        # favor of the lower floor number (matches the reference Go
        # implementation's linear scan over a sorted slice).
        floors = sorted(self._targets)
        nearest = floors[0]
        best = abs(nearest - self.current_floor)
        for f in floors:
            d = abs(f - self.current_floor)
            if d < best:
                best = d
                nearest = f
        self._set_direction_towards(nearest)


class DispatchStrategy(Protocol):
    """Picks the best car in the fleet to service a hall call."""

    def select_car(self, cars: List[Car], floor: int, direction: Direction) -> Car:
        ...


class NearestCarStrategy:
    """Prefers an idle car nearest to the call floor; failing that, a car
    already moving in the requested direction that hasn't yet passed the
    call floor (so it can pick it up en route); failing that, the overall
    nearest car regardless of direction compatibility.
    """

    def select_car(self, cars: List[Car], floor: int, direction: Direction) -> Optional[Car]:
        best: Optional[Car] = None
        best_cost = -1

        for c in cars:
            snap = c.snapshot()
            cost = abs(snap.floor - floor)

            compatible = (
                snap.state == CarState.IDLE
                or (snap.direction == direction and direction == Direction.UP and snap.floor <= floor)
                or (snap.direction == direction and direction == Direction.DOWN and snap.floor >= floor)
            )

            if not compatible:
                cost += c.num_floors  # penalize incompatible cars so idle/en-route cars win when available

            if best is None or cost < best_cost:
                best = c
                best_cost = cost

        return best


class ElevatorSystem:
    """Coordinates a fleet of cars serving a building with num_floors
    floors, dispatching hall calls via a pluggable strategy.
    """

    def __init__(self, num_floors: int, cars: List[Car], strategy: DispatchStrategy):
        self.cars = cars
        self.num_floors = num_floors
        self.strategy = strategy
        # Guards the select-then-assign sequence for hall calls so that two
        # concurrent hall calls can never both "win" the same decision
        # based on a stale view of car state.
        self._lock = threading.Lock()

    def hall_call(self, floor: int, direction: Direction) -> Car:
        """Handles an external request for a car at `floor` heading
        `direction`, dispatching it to the strategy's chosen car and
        returning that car.
        """
        with self._lock:
            car = self.strategy.select_car(self.cars, floor, direction)
            car.add_target(floor)
            return car

    def cabin_request(self, car: Car, destination_floor: int) -> None:
        """Handles an internal request made from inside a specific car."""
        car.add_target(destination_floor)

    def step(self) -> None:
        """Advances every car in the fleet by one simulated tick."""
        for c in self.cars:
            c.step()


def _demo() -> None:
    cars = [Car(0, 10, 0), Car(1, 10, 5)]
    system = ElevatorSystem(10, cars, NearestCarStrategy())

    car = system.hall_call(7, Direction.UP)
    print(f"Hall call at floor 7 dispatched to car {car.id}")

    for _ in range(20):
        system.step()

    for c in system.cars:
        snap = c.snapshot()
        print(f"Car {c.id}: floor={snap.floor} state={snap.state.name} pending={snap.pending}")


if __name__ == "__main__":
    _demo()
