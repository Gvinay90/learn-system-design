"""Food Delivery LLD — Python reference implementation.

Order placement, nearest-partner assignment via the Strategy pattern, and an
order status state machine. See ../README.md for the design writeup.
"""
from __future__ import annotations

import math
import threading
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum, auto
from typing import Dict, List, Optional, Protocol


class OrderStatus(Enum):
    PLACED = auto()
    ACCEPTED = auto()
    PREPARING = auto()
    OUT_FOR_DELIVERY = auto()
    DELIVERED = auto()
    CANCELLED = auto()


VALID_NEXT_STATUS: Dict[OrderStatus, set] = {
    OrderStatus.PLACED: {OrderStatus.ACCEPTED, OrderStatus.CANCELLED},
    OrderStatus.ACCEPTED: {OrderStatus.PREPARING, OrderStatus.CANCELLED},
    OrderStatus.PREPARING: {OrderStatus.OUT_FOR_DELIVERY},
    OrderStatus.OUT_FOR_DELIVERY: {OrderStatus.DELIVERED},
    OrderStatus.DELIVERED: set(),
    OrderStatus.CANCELLED: set(),
}


@dataclass(frozen=True)
class Location:
    x: int
    y: int

    def distance_to(self, other: "Location") -> float:
        return math.sqrt((self.x - other.x) ** 2 + (self.y - other.y) ** 2)


@dataclass
class Customer:
    id: str
    name: str


@dataclass
class MenuItem:
    id: str
    name: str
    price: float


@dataclass
class Restaurant:
    id: str
    name: str
    location: Location
    menu: List[MenuItem]
    is_open: bool

    def find_item(self, item_id: str) -> Optional[MenuItem]:
        for item in self.menu:
            if item.id == item_id:
                return item
        return None


class DeliveryPartner:
    def __init__(self, id: str, name: str, location: Location, available: bool = True):
        self.id = id
        self.name = name
        self.location = location
        self.available = available


@dataclass
class Order:
    id: str
    customer: Customer
    restaurant: Restaurant
    items: List[MenuItem]
    placed_at: datetime
    status: OrderStatus = OrderStatus.PLACED
    partner: Optional[DeliveryPartner] = None


class AssignmentStrategy(Protocol):
    def assign(self, restaurant: Restaurant, partners: List[DeliveryPartner]) -> Optional[DeliveryPartner]: ...


class NearestAvailablePartnerStrategy:
    def assign(self, restaurant: Restaurant, partners: List[DeliveryPartner]) -> Optional[DeliveryPartner]:
        nearest = None
        best = float("inf")
        for p in partners:
            if not p.available:
                continue
            d = restaurant.location.distance_to(p.location)
            if d < best:
                best = d
                nearest = p
        return nearest


class RestaurantClosedError(Exception):
    pass


class ItemNotOnMenuError(Exception):
    pass


class OrderNotFoundError(Exception):
    pass


class NoPartnerAvailableError(Exception):
    pass


class InvalidTransitionError(Exception):
    pass


class OrderAlreadyAssignedError(Exception):
    pass


class FoodDeliverySystem:
    def __init__(self, partners: List[DeliveryPartner], strategy: AssignmentStrategy):
        self.partners = partners
        self.strategy = strategy
        self._orders: Dict[str, Order] = {}
        self._lock = threading.Lock()
        self._seq = 0

    def place_order(self, customer: Customer, restaurant: Restaurant, item_ids: List[str]) -> Order:
        if not restaurant.is_open:
            raise RestaurantClosedError("restaurant is closed")
        items = []
        for item_id in item_ids:
            item = restaurant.find_item(item_id)
            if item is None:
                raise ItemNotOnMenuError("item not on restaurant menu")
            items.append(item)

        with self._lock:
            self._seq += 1
            order = Order(
                id=f"O-{self._seq}",
                customer=customer,
                restaurant=restaurant,
                items=items,
                placed_at=datetime.now(),
            )
            self._orders[order.id] = order
        return order

    def assign_delivery_partner(self, order_id: str) -> DeliveryPartner:
        with self._lock:
            order = self._orders.get(order_id)
            if order is None:
                raise OrderNotFoundError("order not found")
            if order.partner is not None:
                raise OrderAlreadyAssignedError("order already has an assigned delivery partner")

            partner = self.strategy.assign(order.restaurant, self.partners)
            if partner is None:
                raise NoPartnerAvailableError("no delivery partner available")
            partner.available = False
            order.partner = partner
            return partner

    def update_order_status(self, order_id: str, next_status: OrderStatus) -> None:
        with self._lock:
            order = self._orders.get(order_id)
            if order is None:
                raise OrderNotFoundError("order not found")
            if next_status not in VALID_NEXT_STATUS[order.status]:
                raise InvalidTransitionError("invalid order status transition")
            order.status = next_status
            if next_status in (OrderStatus.DELIVERED, OrderStatus.CANCELLED) and order.partner is not None:
                order.partner.available = True
                order.partner = None

    def get_order(self, order_id: str) -> Order:
        order = self._orders.get(order_id)
        if order is None:
            raise OrderNotFoundError("order not found")
        return order


def _demo() -> None:
    restaurant = Restaurant(
        id="R1",
        name="Tasty Bites",
        location=Location(0, 0),
        menu=[MenuItem("I1", "Burger", 5), MenuItem("I2", "Fries", 2)],
        is_open=True,
    )
    partner = DeliveryPartner("P1", "Alex", Location(1, 1))
    sys = FoodDeliverySystem([partner], NearestAvailablePartnerStrategy())
    customer = Customer("C1", "Sam")

    order = sys.place_order(customer, restaurant, ["I1", "I2"])
    print(f"Placed order {order.id}")

    assigned = sys.assign_delivery_partner(order.id)
    print(f"Assigned partner {assigned.id}")

    for status in (OrderStatus.ACCEPTED, OrderStatus.PREPARING, OrderStatus.OUT_FOR_DELIVERY, OrderStatus.DELIVERED):
        sys.update_order_status(order.id, status)
    print(f"Final status: {sys.get_order(order.id).status}")


if __name__ == "__main__":
    _demo()
