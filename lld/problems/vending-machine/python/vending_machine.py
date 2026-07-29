"""Vending Machine using the State pattern: the machine delegates
select_item/insert_money/dispense/cancel to its current state object, which
decides what happens (if anything) from that state. Mirrors the Go/Java
implementations in this same problem directory.
"""

from __future__ import annotations

from dataclasses import dataclass, replace
from typing import Dict, Optional


class VendingMachineError(Exception):
    pass


class InvalidSlotError(VendingMachineError):
    pass


class OutOfStockError(VendingMachineError):
    pass


class NoSelectionError(VendingMachineError):
    pass


class InvalidAmountError(VendingMachineError):
    pass


class NotEnoughMoneyError(VendingMachineError):
    pass


class InvalidStateError(VendingMachineError):
    pass


@dataclass
class Slot:
    item: str
    price: int  # in cents
    quantity: int


@dataclass
class DispenseResult:
    item: str
    change: int  # in cents


class VendingState:
    """Interface every concrete state implements."""

    name: str = "Base"

    def select_item(self, vm: "VendingMachine", slot_id: str) -> None:
        raise InvalidStateError("action not allowed in current state")

    def insert_money(self, vm: "VendingMachine", amount: int) -> None:
        raise InvalidStateError("action not allowed in current state")

    def dispense(self, vm: "VendingMachine") -> DispenseResult:
        raise InvalidStateError("action not allowed in current state")

    def cancel(self, vm: "VendingMachine") -> int:
        raise InvalidStateError("action not allowed in current state")


def _select_item(vm: "VendingMachine", slot_id: str) -> None:
    """Shared SelectItem logic used by both IdleState and OutOfStockState."""
    slot = vm._slot(slot_id)
    if slot.quantity <= 0:
        vm._selected = slot_id
        vm._set_state(OutOfStockState())
        raise OutOfStockError(f"item out of stock: {slot_id}")
    vm._selected = slot_id
    vm._set_state(IdleState())


class IdleState(VendingState):
    name = "Idle"

    def select_item(self, vm: "VendingMachine", slot_id: str) -> None:
        _select_item(vm, slot_id)

    def insert_money(self, vm: "VendingMachine", amount: int) -> None:
        if amount <= 0:
            raise InvalidAmountError("amount must be positive")
        if not vm._selected:
            raise NoSelectionError("no item selected")
        slot = vm._slot(vm._selected)
        vm._balance += amount
        if vm._balance >= slot.price:
            vm._set_state(HasMoneyState())

    def dispense(self, vm: "VendingMachine") -> DispenseResult:
        raise NotEnoughMoneyError("not enough money inserted")

    def cancel(self, vm: "VendingMachine") -> int:
        refund = vm._balance
        vm._balance = 0
        vm._selected = None
        return refund


class HasMoneyState(VendingState):
    name = "HasMoney"

    def insert_money(self, vm: "VendingMachine", amount: int) -> None:
        if amount <= 0:
            raise InvalidAmountError("amount must be positive")
        vm._balance += amount

    def dispense(self, vm: "VendingMachine") -> DispenseResult:
        vm._set_state(DispensingState())
        return vm._state.dispense(vm)

    def cancel(self, vm: "VendingMachine") -> int:
        refund = vm._balance
        vm._balance = 0
        vm._selected = None
        vm._set_state(IdleState())
        return refund


class DispensingState(VendingState):
    """Entered only transiently from HasMoneyState.dispense."""

    name = "Dispensing"

    def dispense(self, vm: "VendingMachine") -> DispenseResult:
        try:
            slot = vm._slot(vm._selected)
        except InvalidSlotError:
            vm._set_state(IdleState())
            raise

        if vm._balance < slot.price:
            vm._set_state(HasMoneyState())
            raise NotEnoughMoneyError("not enough money inserted")

        change = vm._balance - slot.price
        slot.quantity -= 1
        item = slot.item

        vm._balance = 0
        vm._selected = None
        vm._set_state(IdleState())

        return DispenseResult(item=item, change=change)


class OutOfStockState(VendingState):
    name = "OutOfStock"

    def select_item(self, vm: "VendingMachine", slot_id: str) -> None:
        _select_item(vm, slot_id)

    def cancel(self, vm: "VendingMachine") -> int:
        refund = vm._balance
        vm._balance = 0
        vm._selected = None
        vm._set_state(IdleState())
        return refund


class VendingMachine:
    """Context: holds the current state, inventory, selected slot, and
    money collected so far, and forwards every action to its current state.
    """

    def __init__(self, slots: Dict[str, Slot]):
        self._slots: Dict[str, Slot] = {sid: replace(s) for sid, s in slots.items()}
        self._state: VendingState = IdleState()
        self._selected: Optional[str] = None
        self._balance: int = 0

    @property
    def state(self) -> VendingState:
        return self._state

    @property
    def balance(self) -> int:
        return self._balance

    @property
    def selected(self) -> Optional[str]:
        return self._selected

    def inventory(self, slot_id: str) -> Optional[Slot]:
        slot = self._slots.get(slot_id)
        return replace(slot) if slot else None

    def select_item(self, slot_id: str) -> None:
        self._state.select_item(self, slot_id)

    def insert_money(self, amount: int) -> None:
        self._state.insert_money(self, amount)

    def dispense(self) -> DispenseResult:
        return self._state.dispense(self)

    def cancel(self) -> int:
        return self._state.cancel(self)

    def _set_state(self, state: VendingState) -> None:
        self._state = state

    def _slot(self, slot_id: str) -> Slot:
        slot = self._slots.get(slot_id)
        if slot is None:
            raise InvalidSlotError(f"invalid slot: {slot_id}")
        return slot


if __name__ == "__main__":
    vm = VendingMachine(
        {
            "A1": Slot(item="Soda", price=150, quantity=2),
            "A2": Slot(item="Chips", price=200, quantity=1),
            "B1": Slot(item="Water", price=100, quantity=0),
        }
    )
    vm.select_item("A1")
    vm.insert_money(200)
    result = vm.dispense()
    print(f"Dispensed {result.item}, change owed: {result.change} cents")
