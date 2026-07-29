import threading
from datetime import timedelta

import pytest

from parking_lot import (
    Floor,
    HourlyTieredPricing,
    LotFullError,
    ParkingLot,
    ParkingSpot,
    TicketNotFoundError,
    Vehicle,
    VehicleType,
)


def new_test_lot() -> ParkingLot:
    floor1 = Floor(1, [
        ParkingSpot("F1-M1", VehicleType.MOTORCYCLE),
        ParkingSpot("F1-C1", VehicleType.COMPACT),
        ParkingSpot("F1-L1", VehicleType.LARGE),
    ])
    return ParkingLot([floor1], HourlyTieredPricing(first_hour_rate=5, subsequent_rate=2))


def test_park_and_unpark():
    lot = new_test_lot()
    vehicle = Vehicle("KA-01", VehicleType.COMPACT)

    ticket = lot.park_vehicle(vehicle)
    assert ticket.spot.id == "F1-C1"

    fee = lot.unpark_vehicle(ticket.id, ticket.entry_time + timedelta(minutes=30))
    assert fee == 5


def test_lot_full_for_type():
    lot = new_test_lot()
    lot.park_vehicle(Vehicle("A1", VehicleType.COMPACT))
    with pytest.raises(LotFullError):
        lot.park_vehicle(Vehicle("A2", VehicleType.COMPACT))


def test_unpark_unknown_ticket():
    lot = new_test_lot()
    with pytest.raises(TicketNotFoundError):
        lot.unpark_vehicle("bogus", __import__("datetime").datetime.now())


def test_hourly_tiered_pricing():
    lot = new_test_lot()
    ticket = lot.park_vehicle(Vehicle("KA-02", VehicleType.LARGE))
    fee = lot.unpark_vehicle(ticket.id, ticket.entry_time + timedelta(hours=3))
    assert fee == 9  # 5 first hour + 2 hours * 2


def test_concurrent_parking_same_spot():
    lot = new_test_lot()
    successes = []
    lock = threading.Lock()

    def worker():
        try:
            lot.park_vehicle(Vehicle("race", VehicleType.MOTORCYCLE))
            with lock:
                successes.append(True)
        except LotFullError:
            pass

    threads = [threading.Thread(target=worker) for _ in range(2)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert len(successes) == 1
