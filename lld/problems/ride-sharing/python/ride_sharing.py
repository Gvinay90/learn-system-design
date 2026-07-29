"""Ride Sharing LLD — Python reference implementation.

Rider/driver matching (Strategy), trip lifecycle (state machine), and
pluggable pricing. See ../README.md for the design writeup.
"""
from __future__ import annotations

import math
import threading
from abc import ABC, abstractmethod
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum, auto
from typing import List, Optional


@dataclass(frozen=True)
class Location:
    x: float
    y: float

    def distance_to(self, other: "Location") -> float:
        return math.sqrt((self.x - other.x) ** 2 + (self.y - other.y) ** 2)


@dataclass
class Rider:
    id: str
    name: str


@dataclass
class Driver:
    id: str
    name: str
    location: Location
    available: bool = True


class TripStatus(Enum):
    REQUESTED = auto()
    ACCEPTED = auto()
    IN_PROGRESS = auto()
    COMPLETED = auto()
    CANCELLED = auto()


@dataclass
class Trip:
    id: str
    rider: Rider
    pickup: Location
    dropoff: Location
    driver: Optional[Driver] = None
    status: TripStatus = TripStatus.REQUESTED
    requested_at: datetime = field(default_factory=datetime.now)
    fare: float = 0.0


class TripNotFoundError(Exception):
    """Raised when a trip ID is not recognized."""


class NoDriverAvailableError(Exception):
    """Raised when no available driver can be matched."""


class InvalidTransitionError(Exception):
    """Raised when a trip status transition is not legal."""


class DriverMatchingStrategy(ABC):
    """Picks a driver for a trip from the available pool."""

    @abstractmethod
    def match(self, pickup: Location, drivers: List[Driver]) -> Optional[Driver]:
        ...


class NearestAvailableDriverStrategy(DriverMatchingStrategy):
    """Picks the closest available driver to the pickup point."""

    def match(self, pickup: Location, drivers: List[Driver]) -> Optional[Driver]:
        best: Optional[Driver] = None
        best_dist = math.inf
        for d in drivers:
            if not d.available:
                continue
            dist = pickup.distance_to(d.location)
            if dist < best_dist:
                best_dist = dist
                best = d
        return best


class PricingStrategy(ABC):
    """Computes the fare for a completed trip."""

    @abstractmethod
    def calculate_fare(self, trip: Trip) -> float:
        ...


@dataclass
class DistanceBasedPricing(PricingStrategy):
    """Charges a base fare plus a per-unit-distance rate."""

    base_fare: float
    per_distance: float

    def calculate_fare(self, trip: Trip) -> float:
        return self.base_fare + self.per_distance * trip.pickup.distance_to(trip.dropoff)


class RideSharingSystem:
    """Facade coordinating the driver pool, trip map, matching, and pricing.

    All lifecycle methods acquire a single system lock so that, in
    particular, the find-and-mark-unavailable step in match_driver is
    atomic: two concurrent calls can never both be assigned the same
    single driver.
    """

    def __init__(self, drivers: List[Driver], matching: DriverMatchingStrategy, pricing: PricingStrategy):
        self.matching = matching
        self.pricing = pricing
        self._lock = threading.Lock()
        self._drivers = drivers
        self._trips: dict[str, Trip] = {}
        self._seq = 0

    def request_ride(self, rider: Rider, pickup: Location, dropoff: Location) -> Trip:
        with self._lock:
            self._seq += 1
            trip = Trip(id=f"TR-{self._seq}", rider=rider, pickup=pickup, dropoff=dropoff)
            self._trips[trip.id] = trip
            return trip

    def match_driver(self, trip_id: str) -> Driver:
        """Finds the nearest available driver and assigns it atomically."""
        with self._lock:
            trip = self._trips.get(trip_id)
            if trip is None:
                raise TripNotFoundError(trip_id)
            if trip.status != TripStatus.REQUESTED:
                raise InvalidTransitionError(trip.status)

            driver = self.matching.match(trip.pickup, self._drivers)
            if driver is None:
                raise NoDriverAvailableError()

            driver.available = False
            trip.driver = driver
            trip.status = TripStatus.ACCEPTED
            return driver

    def start_trip(self, trip_id: str) -> None:
        with self._lock:
            trip = self._trips.get(trip_id)
            if trip is None:
                raise TripNotFoundError(trip_id)
            if trip.status != TripStatus.ACCEPTED:
                raise InvalidTransitionError(trip.status)
            trip.status = TripStatus.IN_PROGRESS

    def complete_trip(self, trip_id: str) -> float:
        with self._lock:
            trip = self._trips.get(trip_id)
            if trip is None:
                raise TripNotFoundError(trip_id)
            if trip.status != TripStatus.IN_PROGRESS:
                raise InvalidTransitionError(trip.status)
            trip.status = TripStatus.COMPLETED
            trip.fare = self.pricing.calculate_fare(trip)
            if trip.driver is not None:
                trip.driver.available = True
            return trip.fare

    def cancel_trip(self, trip_id: str) -> None:
        with self._lock:
            trip = self._trips.get(trip_id)
            if trip is None:
                raise TripNotFoundError(trip_id)
            if trip.status not in (TripStatus.REQUESTED, TripStatus.ACCEPTED):
                raise InvalidTransitionError(trip.status)
            trip.status = TripStatus.CANCELLED
            if trip.driver is not None:
                trip.driver.available = True


def _demo() -> None:
    drivers = [
        Driver("D1", "Alice", Location(0, 0)),
        Driver("D2", "Bob", Location(10, 10)),
    ]
    system = RideSharingSystem(drivers, NearestAvailableDriverStrategy(), DistanceBasedPricing(2, 1.5))

    rider = Rider("R1", "Riya")
    trip = system.request_ride(rider, Location(0, 0), Location(3, 4))
    print(f"Requested trip {trip.id} status={trip.status}")

    driver = system.match_driver(trip.id)
    print(f"Matched driver {driver.id} status={trip.status}")

    system.start_trip(trip.id)
    print(f"Started trip, status={trip.status}")

    fare = system.complete_trip(trip.id)
    print(f"Completed trip, fare={fare} status={trip.status}")


if __name__ == "__main__":
    _demo()
