import threading

import pytest

from ride_sharing import (
    DistanceBasedPricing,
    Driver,
    InvalidTransitionError,
    Location,
    NearestAvailableDriverStrategy,
    NoDriverAvailableError,
    RideSharingSystem,
    Rider,
    TripStatus,
)


def new_test_system() -> RideSharingSystem:
    drivers = [
        Driver("D1", "Alice", Location(0, 0), available=True),
        Driver("D2", "Bob", Location(10, 10), available=True),
    ]
    return RideSharingSystem(drivers, NearestAvailableDriverStrategy(), DistanceBasedPricing(2, 1.5))


def test_happy_path_lifecycle():
    sys = new_test_system()
    rider = Rider("R1", "Riya")

    trip = sys.request_ride(rider, Location(0, 0), Location(3, 4))
    assert trip.status == TripStatus.REQUESTED

    driver = sys.match_driver(trip.id)
    assert driver.id == "D1"
    assert trip.status == TripStatus.ACCEPTED

    sys.start_trip(trip.id)
    assert trip.status == TripStatus.IN_PROGRESS

    fare = sys.complete_trip(trip.id)
    # base 2 + 5 (distance from (0,0) to (3,4)) * 1.5 = 9.5
    assert fare == 9.5
    assert trip.status == TripStatus.COMPLETED
    assert driver.available


def test_invalid_transition_rejected():
    sys = new_test_system()
    rider = Rider("R1", "Riya")
    trip = sys.request_ride(rider, Location(0, 0), Location(1, 1))

    with pytest.raises(InvalidTransitionError):
        sys.complete_trip(trip.id)

    sys.match_driver(trip.id)
    sys.start_trip(trip.id)
    sys.complete_trip(trip.id)

    with pytest.raises(InvalidTransitionError):
        sys.cancel_trip(trip.id)


def test_no_driver_available():
    sys = new_test_system()
    rider = Rider("R1", "Riya")

    t1 = sys.request_ride(rider, Location(0, 0), Location(1, 1))
    t2 = sys.request_ride(rider, Location(0, 0), Location(1, 1))
    t3 = sys.request_ride(rider, Location(0, 0), Location(1, 1))

    sys.match_driver(t1.id)
    sys.match_driver(t2.id)

    with pytest.raises(NoDriverAvailableError):
        sys.match_driver(t3.id)


def test_concurrent_matching():
    """Two threads racing for the same single available driver must never
    both succeed — the lock in match_driver must serialize them."""
    drivers = [Driver("D1", "Alice", Location(0, 0), available=True)]
    sys = RideSharingSystem(drivers, NearestAvailableDriverStrategy(), DistanceBasedPricing(2, 1.5))
    rider = Rider("R1", "Riya")

    trip1 = sys.request_ride(rider, Location(0, 0), Location(1, 1))
    trip2 = sys.request_ride(rider, Location(0, 0), Location(1, 1))

    results = []
    results_lock = threading.Lock()

    def worker(trip_id: str) -> None:
        try:
            sys.match_driver(trip_id)
            with results_lock:
                results.append(None)
        except Exception as e:  # noqa: BLE001
            with results_lock:
                results.append(e)

    threads = [threading.Thread(target=worker, args=(tid,)) for tid in (trip1.id, trip2.id)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    success_count = sum(1 for r in results if r is None)
    assert success_count == 1
