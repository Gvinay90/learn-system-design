"""Parking Lot LLD — Python reference implementation.

Multiple floors, multiple spot sizes, nearest-spot allocation, pluggable
pricing via the Strategy pattern. See ../README.md for the design writeup.
"""
from __future__ import annotations

import threading
from dataclasses import dataclass, field
from datetime import datetime, timedelta
from enum import Enum, auto
from typing import Dict, List, Optional, Protocol


class VehicleType(Enum):
    MOTORCYCLE = auto()
    COMPACT = auto()
    LARGE = auto()


@dataclass
class Vehicle:
    license_plate: str
    type: VehicleType


class ParkingSpot:
    def __init__(self, spot_id: str, spot_type: VehicleType):
        self.id = spot_id
        self.type = spot_type
        self.occupied = False
        self.vehicle: Optional[Vehicle] = None
        self._lock = threading.Lock()

    def try_assign(self, vehicle: Vehicle) -> bool:
        with self._lock:
            if self.occupied or self.type != vehicle.type:
                return False
            self.occupied = True
            self.vehicle = vehicle
            return True

    def free(self) -> None:
        with self._lock:
            self.occupied = False
            self.vehicle = None


class Floor:
    def __init__(self, level: int, spots: List[ParkingSpot]):
        self.level = level
        self.spots = spots
        self._lock = threading.Lock()

    def find_and_assign(self, vehicle: Vehicle) -> Optional[ParkingSpot]:
        with self._lock:
            for spot in self.spots:
                if spot.try_assign(vehicle):
                    return spot
        return None

    def free(self, spot_id: str) -> None:
        for spot in self.spots:
            if spot.id == spot_id:
                spot.free()
                return


@dataclass
class Ticket:
    id: str
    vehicle: Vehicle
    spot: ParkingSpot
    floor_level: int
    entry_time: datetime


class PricingStrategy(Protocol):
    def calculate_fee(self, ticket: Ticket, exit_time: datetime) -> float: ...


@dataclass
class HourlyTieredPricing:
    first_hour_rate: float
    subsequent_rate: float

    def calculate_fee(self, ticket: Ticket, exit_time: datetime) -> float:
        hours = (exit_time - ticket.entry_time).total_seconds() / 3600.0
        if hours <= 1:
            return self.first_hour_rate
        return self.first_hour_rate + (hours - 1) * self.subsequent_rate


class LotFullError(Exception):
    pass


class TicketNotFoundError(Exception):
    pass


class ParkingLot:
    def __init__(self, floors: List[Floor], pricing: PricingStrategy):
        self.floors = floors
        self.pricing = pricing
        self._tickets: Dict[str, Ticket] = {}
        self._lock = threading.Lock()
        self._seq = 0

    def park_vehicle(self, vehicle: Vehicle) -> Ticket:
        for floor in self.floors:
            spot = floor.find_and_assign(vehicle)
            if spot is not None:
                with self._lock:
                    self._seq += 1
                    ticket = Ticket(
                        id=f"T-{self._seq}",
                        vehicle=vehicle,
                        spot=spot,
                        floor_level=floor.level,
                        entry_time=datetime.now(),
                    )
                    self._tickets[ticket.id] = ticket
                return ticket
        raise LotFullError("no available spot for vehicle type")

    def unpark_vehicle(self, ticket_id: str, exit_time: datetime) -> float:
        with self._lock:
            ticket = self._tickets.pop(ticket_id, None)
        if ticket is None:
            raise TicketNotFoundError("ticket not recognized")

        for floor in self.floors:
            if floor.level == ticket.floor_level:
                floor.free(ticket.spot.id)
                break
        return self.pricing.calculate_fee(ticket, exit_time)


def _demo() -> None:
    floor1 = Floor(1, [
        ParkingSpot("F1-M1", VehicleType.MOTORCYCLE),
        ParkingSpot("F1-C1", VehicleType.COMPACT),
        ParkingSpot("F1-L1", VehicleType.LARGE),
    ])
    lot = ParkingLot([floor1], HourlyTieredPricing(first_hour_rate=5, subsequent_rate=2))

    car = Vehicle("KA-01", VehicleType.COMPACT)
    ticket = lot.park_vehicle(car)
    print(f"Parked {car.license_plate} at spot {ticket.spot.id}")

    fee = lot.unpark_vehicle(ticket.id, ticket.entry_time + timedelta(hours=3))
    print(f"Fee for 3 hours: {fee}")


if __name__ == "__main__":
    _demo()
