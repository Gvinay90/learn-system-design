"""Strategy pattern — ShoppingCart delegates pricing to a pluggable strategy.

See ../README.md for the design writeup.
"""
from __future__ import annotations

from dataclasses import dataclass
from typing import List, Protocol


class PricingStrategy(Protocol):
    def apply_discount(self, subtotal: float) -> float: ...


@dataclass
class RegularPricing:
    def apply_discount(self, subtotal: float) -> float:
        return subtotal


@dataclass
class PercentageDiscountPricing:
    percent_off: float

    def apply_discount(self, subtotal: float) -> float:
        return subtotal * (1 - self.percent_off / 100)


@dataclass
class ClearancePricing:
    percent_off: float
    flat_off: float

    def apply_discount(self, subtotal: float) -> float:
        discounted = subtotal * (1 - self.percent_off / 100) - self.flat_off
        return max(discounted, 0)


@dataclass
class Item:
    name: str
    price: float
    quantity: int


class ShoppingCart:
    """Delegates final price computation to whatever PricingStrategy it's
    configured with, so pricing schemes can change at runtime."""

    def __init__(self, strategy: PricingStrategy):
        self.items: List[Item] = []
        self.strategy = strategy

    def add_item(self, item: Item) -> None:
        self.items.append(item)

    def subtotal(self) -> float:
        return sum(item.price * item.quantity for item in self.items)

    def checkout(self) -> float:
        return self.strategy.apply_discount(self.subtotal())


def _demo() -> None:
    cart = ShoppingCart(RegularPricing())
    cart.add_item(Item("book", 20, 2))
    print(f"Regular checkout: {cart.checkout()}")

    cart.strategy = PercentageDiscountPricing(percent_off=10)
    print(f"10% off checkout: {cart.checkout()}")


if __name__ == "__main__":
    _demo()
