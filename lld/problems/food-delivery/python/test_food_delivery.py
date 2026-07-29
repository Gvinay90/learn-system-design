import threading

import pytest

from food_delivery import (
    Customer,
    DeliveryPartner,
    FoodDeliverySystem,
    InvalidTransitionError,
    ItemNotOnMenuError,
    Location,
    MenuItem,
    NearestAvailablePartnerStrategy,
    NoPartnerAvailableError,
    OrderStatus,
    Restaurant,
)


def make_test_restaurant() -> Restaurant:
    return Restaurant(
        id="R1",
        name="Tasty Bites",
        location=Location(0, 0),
        menu=[MenuItem("I1", "Burger", 5), MenuItem("I2", "Fries", 2)],
        is_open=True,
    )


def new_test_system(*partners: DeliveryPartner) -> FoodDeliverySystem:
    return FoodDeliverySystem(list(partners), NearestAvailablePartnerStrategy())


def test_happy_path_place_assign_deliver():
    restaurant = make_test_restaurant()
    partner = DeliveryPartner("P1", "Alex", Location(1, 1))
    sys = new_test_system(partner)
    customer = Customer("C1", "Sam")

    order = sys.place_order(customer, restaurant, ["I1", "I2"])

    assigned = sys.assign_delivery_partner(order.id)
    assert assigned.id == "P1"
    assert partner.available is False

    for status in (OrderStatus.ACCEPTED, OrderStatus.PREPARING, OrderStatus.OUT_FOR_DELIVERY, OrderStatus.DELIVERED):
        sys.update_order_status(order.id, status)

    assert sys.get_order(order.id).status == OrderStatus.DELIVERED
    assert partner.available is True


def test_invalid_status_transition_rejected():
    restaurant = make_test_restaurant()
    partner = DeliveryPartner("P1", "Alex", Location(1, 1))
    sys = new_test_system(partner)
    customer = Customer("C1", "Sam")

    order = sys.place_order(customer, restaurant, ["I1"])
    with pytest.raises(InvalidTransitionError):
        sys.update_order_status(order.id, OrderStatus.DELIVERED)

    sys.update_order_status(order.id, OrderStatus.ACCEPTED)
    sys.update_order_status(order.id, OrderStatus.PREPARING)
    sys.update_order_status(order.id, OrderStatus.OUT_FOR_DELIVERY)
    sys.update_order_status(order.id, OrderStatus.DELIVERED)

    with pytest.raises(InvalidTransitionError):
        sys.update_order_status(order.id, OrderStatus.CANCELLED)


def test_item_not_on_menu_and_no_partner_available():
    restaurant = make_test_restaurant()
    sys = new_test_system()
    customer = Customer("C1", "Sam")

    with pytest.raises(ItemNotOnMenuError):
        sys.place_order(customer, restaurant, ["BOGUS"])

    order = sys.place_order(customer, restaurant, ["I1"])
    with pytest.raises(NoPartnerAvailableError):
        sys.assign_delivery_partner(order.id)


def test_concurrent_assignment():
    restaurant = make_test_restaurant()
    partner = DeliveryPartner("P1", "Alex", Location(1, 1))
    sys = new_test_system(partner)
    customer = Customer("C1", "Sam")

    order1 = sys.place_order(customer, restaurant, ["I1"])
    order2 = sys.place_order(customer, restaurant, ["I2"])

    successes = []
    lock = threading.Lock()

    def worker(order_id: str):
        try:
            sys.assign_delivery_partner(order_id)
            with lock:
                successes.append(True)
        except NoPartnerAvailableError:
            pass

    threads = [threading.Thread(target=worker, args=(oid,)) for oid in (order1.id, order2.id)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert len(successes) == 1
