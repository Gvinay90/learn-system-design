import pytest

from state import IllegalTransitionError, Order


def test_happy_path_transitions():
    order = Order("O-1")
    assert order.state.name == "Created"

    order.pay()
    assert order.state.name == "Paid"

    order.ship()
    assert order.state.name == "Shipped"

    order.deliver()
    assert order.state.name == "Delivered"


def test_cannot_ship_before_payment():
    order = Order("O-2")
    with pytest.raises(IllegalTransitionError):
        order.ship()
    assert order.state.name == "Created"


def test_cancel_allowed_before_shipping():
    order = Order("O-3")
    order.pay()
    order.cancel()
    assert order.state.name == "Cancelled"


def test_cannot_cancel_after_shipping():
    order = Order("O-4")
    order.pay()
    order.ship()
    with pytest.raises(IllegalTransitionError):
        order.cancel()
    assert order.state.name == "Shipped"


def test_delivered_is_terminal():
    order = Order("O-5")
    order.pay()
    order.ship()
    order.deliver()

    with pytest.raises(IllegalTransitionError):
        order.cancel()
    with pytest.raises(IllegalTransitionError):
        order.pay()
