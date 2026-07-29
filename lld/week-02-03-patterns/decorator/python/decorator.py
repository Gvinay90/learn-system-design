"""Decorator pattern — a base Coffee wrapped by stackable add-on decorators.
See ../README.md for the design writeup."""
from __future__ import annotations

from typing import Protocol


class Coffee(Protocol):
    def cost(self) -> float: ...
    def description(self) -> str: ...


class Espresso:
    def cost(self) -> float:
        return 2.0

    def description(self) -> str:
        return "Espresso"


class _CoffeeDecorator:
    def __init__(self, wrapped: Coffee):
        self._wrapped = wrapped


class MilkDecorator(_CoffeeDecorator):
    def cost(self) -> float:
        return self._wrapped.cost() + 0.5

    def description(self) -> str:
        return self._wrapped.description() + " + Milk"


class SugarDecorator(_CoffeeDecorator):
    def cost(self) -> float:
        return self._wrapped.cost() + 0.25

    def description(self) -> str:
        return self._wrapped.description() + " + Sugar"


class WhipDecorator(_CoffeeDecorator):
    def cost(self) -> float:
        return self._wrapped.cost() + 0.75

    def description(self) -> str:
        return self._wrapped.description() + " + Whip"


def _demo() -> None:
    coffee: Coffee = Espresso()
    coffee = MilkDecorator(coffee)
    coffee = SugarDecorator(coffee)
    coffee = WhipDecorator(coffee)
    print(f"{coffee.description()} costs {coffee.cost()}")


if __name__ == "__main__":
    _demo()
