import pytest

from vending_machine import (
    DispensingState,
    HasMoneyState,
    IdleState,
    InvalidAmountError,
    InvalidSlotError,
    InvalidStateError,
    NoSelectionError,
    NotEnoughMoneyError,
    OutOfStockError,
    OutOfStockState,
    Slot,
    VendingMachine,
)


def new_test_machine() -> VendingMachine:
    return VendingMachine(
        {
            "A1": Slot(item="Soda", price=150, quantity=2),
            "A2": Slot(item="Chips", price=200, quantity=1),
            "B1": Slot(item="Water", price=100, quantity=0),
        }
    )


def test_select_and_insert_exact_money_then_dispense():
    vm = new_test_machine()

    vm.select_item("A1")
    vm.insert_money(150)
    assert vm.state.name == "HasMoney"

    result = vm.dispense()
    assert result.item == "Soda"
    assert result.change == 0
    assert vm.state.name == "Idle"


def test_insert_money_gives_change_on_overpay():
    vm = new_test_machine()

    vm.select_item("A1")
    vm.insert_money(200)

    result = vm.dispense()
    assert result.change == 50


def test_insufficient_money_then_top_up_then_dispense():
    vm = new_test_machine()

    vm.select_item("A2")  # Chips, price 200
    vm.insert_money(100)
    assert vm.state.name == "Idle"

    with pytest.raises(NotEnoughMoneyError):
        vm.dispense()

    vm.insert_money(100)
    assert vm.state.name == "HasMoney"

    result = vm.dispense()
    assert result.item == "Chips"
    assert result.change == 0


def test_cancel_refunds_inserted_money():
    vm = new_test_machine()

    vm.select_item("A1")
    vm.insert_money(75)

    refund = vm.cancel()
    assert refund == 75
    assert vm.state.name == "Idle"
    assert vm.selected is None


def test_select_out_of_stock_slot():
    vm = new_test_machine()

    with pytest.raises(OutOfStockError):
        vm.select_item("B1")
    assert vm.state.name == "OutOfStock"

    # Re-selecting an in-stock slot recovers back to Idle.
    vm.select_item("A1")
    assert vm.state.name == "Idle"


def test_select_invalid_slot():
    vm = new_test_machine()
    with pytest.raises(InvalidSlotError):
        vm.select_item("Z9")


def test_insert_money_without_selection_raises():
    vm = new_test_machine()
    with pytest.raises(NoSelectionError):
        vm.insert_money(100)


def test_insert_non_positive_amount_raises():
    vm = new_test_machine()
    vm.select_item("A1")
    with pytest.raises(InvalidAmountError):
        vm.insert_money(0)
    with pytest.raises(InvalidAmountError):
        vm.insert_money(-10)


def test_select_item_rejected_in_has_money_state():
    vm = new_test_machine()
    vm.select_item("A1")
    vm.insert_money(150)
    assert vm.state.name == "HasMoney"

    with pytest.raises(InvalidStateError):
        vm.select_item("A2")


def test_dispense_decrements_inventory():
    vm = new_test_machine()
    vm.select_item("A1")
    vm.insert_money(150)
    vm.dispense()

    slot = vm.inventory("A1")
    assert slot.quantity == 1


def test_dispense_without_money_raises():
    vm = new_test_machine()
    vm.select_item("A1")
    with pytest.raises(NotEnoughMoneyError):
        vm.dispense()


if __name__ == "__main__":
    import sys

    sys.exit(pytest.main([__file__, "-v"]))
