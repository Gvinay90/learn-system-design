"""State pattern — an Order delegates behavior to its current OrderState,
which decides what transition, if any, is legal from that state.

See ../README.md for the design writeup.
"""
from __future__ import annotations

from typing import Protocol


class IllegalTransitionError(Exception):
    pass


class OrderState(Protocol):
    name: str

    def pay(self, order: "Order") -> None: ...
    def ship(self, order: "Order") -> None: ...
    def deliver(self, order: "Order") -> None: ...
    def cancel(self, order: "Order") -> None: ...


class Order:
    """The context: holds a reference to its current state and forwards
    every action to it."""

    def __init__(self, order_id: str):
        self.id = order_id
        self.state: OrderState = CreatedState()

    def pay(self) -> None:
        self.state.pay(self)

    def ship(self) -> None:
        self.state.ship(self)

    def deliver(self) -> None:
        self.state.deliver(self)

    def cancel(self) -> None:
        self.state.cancel(self)


class CreatedState:
    name = "Created"

    def pay(self, order: Order) -> None:
        order.state = PaidState()

    def ship(self, order: Order) -> None:
        raise IllegalTransitionError("cannot ship an order that has not been paid")

    def deliver(self, order: Order) -> None:
        raise IllegalTransitionError("cannot deliver an order that has not shipped")

    def cancel(self, order: Order) -> None:
        order.state = CancelledState()


class PaidState:
    name = "Paid"

    def pay(self, order: Order) -> None:
        raise IllegalTransitionError("order already paid")

    def ship(self, order: Order) -> None:
        order.state = ShippedState()

    def deliver(self, order: Order) -> None:
        raise IllegalTransitionError("cannot deliver an order that has not shipped")

    def cancel(self, order: Order) -> None:
        order.state = CancelledState()


class ShippedState:
    name = "Shipped"

    def pay(self, order: Order) -> None:
        raise IllegalTransitionError("order already paid")

    def ship(self, order: Order) -> None:
        raise IllegalTransitionError("order already shipped")

    def deliver(self, order: Order) -> None:
        order.state = DeliveredState()

    def cancel(self, order: Order) -> None:
        raise IllegalTransitionError("cannot cancel an order that has already shipped")


class DeliveredState:
    name = "Delivered"

    def pay(self, order: Order) -> None:
        raise IllegalTransitionError("order already delivered")

    def ship(self, order: Order) -> None:
        raise IllegalTransitionError("order already delivered")

    def deliver(self, order: Order) -> None:
        raise IllegalTransitionError("order already delivered")

    def cancel(self, order: Order) -> None:
        raise IllegalTransitionError("cannot cancel a delivered order")


class CancelledState:
    name = "Cancelled"

    def pay(self, order: Order) -> None:
        raise IllegalTransitionError("order is cancelled")

    def ship(self, order: Order) -> None:
        raise IllegalTransitionError("order is cancelled")

    def deliver(self, order: Order) -> None:
        raise IllegalTransitionError("order is cancelled")

    def cancel(self, order: Order) -> None:
        raise IllegalTransitionError("order is already cancelled")


def _demo() -> None:
    order = Order("O-1")
    print(f"Initial state: {order.state.name}")
    order.pay()
    order.ship()
    order.deliver()
    print(f"Final state: {order.state.name}")


if __name__ == "__main__":
    _demo()
