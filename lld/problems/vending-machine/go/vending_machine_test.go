package vendingmachine

import "testing"

func newTestMachine() *VendingMachine {
	return NewVendingMachine(map[string]Slot{
		"A1": {Item: "Soda", Price: 150, Quantity: 2},
		"A2": {Item: "Chips", Price: 200, Quantity: 1},
		"B1": {Item: "Water", Price: 100, Quantity: 0},
	})
}

func TestSelectAndInsertExactMoneyThenDispense(t *testing.T) {
	vm := newTestMachine()

	if err := vm.SelectItem("A1"); err != nil {
		t.Fatalf("SelectItem: unexpected err: %v", err)
	}
	if err := vm.InsertMoney(150); err != nil {
		t.Fatalf("InsertMoney: unexpected err: %v", err)
	}
	if vm.State().Name() != "HasMoney" {
		t.Fatalf("expected HasMoney state, got %s", vm.State().Name())
	}

	result, err := vm.Dispense()
	if err != nil {
		t.Fatalf("Dispense: unexpected err: %v", err)
	}
	if result.Item != "Soda" {
		t.Fatalf("expected Soda, got %s", result.Item)
	}
	if result.Change != 0 {
		t.Fatalf("expected no change, got %d", result.Change)
	}
	if vm.State().Name() != "Idle" {
		t.Fatalf("expected Idle state after dispense, got %s", vm.State().Name())
	}
}

func TestInsertMoneyGivesChangeOnOverpay(t *testing.T) {
	vm := newTestMachine()

	_ = vm.SelectItem("A1")
	_ = vm.InsertMoney(200)

	result, err := vm.Dispense()
	if err != nil {
		t.Fatalf("Dispense: unexpected err: %v", err)
	}
	if result.Change != 50 {
		t.Fatalf("expected 50 change, got %d", result.Change)
	}
}

func TestInsufficientMoneyThenTopUpThenDispense(t *testing.T) {
	vm := newTestMachine()

	_ = vm.SelectItem("A2") // Chips, price 200
	if err := vm.InsertMoney(100); err != nil {
		t.Fatalf("InsertMoney: unexpected err: %v", err)
	}
	if vm.State().Name() != "Idle" {
		t.Fatalf("expected still Idle after partial payment, got %s", vm.State().Name())
	}

	if _, err := vm.Dispense(); err != ErrNotEnoughMoney {
		t.Fatalf("expected ErrNotEnoughMoney, got %v", err)
	}

	if err := vm.InsertMoney(100); err != nil {
		t.Fatalf("InsertMoney (top-up): unexpected err: %v", err)
	}
	if vm.State().Name() != "HasMoney" {
		t.Fatalf("expected HasMoney after top-up, got %s", vm.State().Name())
	}

	result, err := vm.Dispense()
	if err != nil {
		t.Fatalf("Dispense: unexpected err: %v", err)
	}
	if result.Item != "Chips" || result.Change != 0 {
		t.Fatalf("unexpected dispense result: %+v", result)
	}
}

func TestCancelRefundsInsertedMoney(t *testing.T) {
	vm := newTestMachine()

	_ = vm.SelectItem("A1")
	_ = vm.InsertMoney(75)

	refund, err := vm.Cancel()
	if err != nil {
		t.Fatalf("Cancel: unexpected err: %v", err)
	}
	if refund != 75 {
		t.Fatalf("expected refund of 75, got %d", refund)
	}
	if vm.State().Name() != "Idle" {
		t.Fatalf("expected Idle after cancel, got %s", vm.State().Name())
	}
	if vm.Balance() != 0 {
		t.Fatalf("expected balance reset to 0, got %d", vm.Balance())
	}
	if vm.Selected() != "" {
		t.Fatalf("expected selection cleared, got %q", vm.Selected())
	}
}

func TestCancelFromHasMoneyRefundsFull(t *testing.T) {
	vm := newTestMachine()

	_ = vm.SelectItem("A1")
	_ = vm.InsertMoney(150) // exact price -> HasMoney

	if vm.State().Name() != "HasMoney" {
		t.Fatalf("expected HasMoney, got %s", vm.State().Name())
	}

	refund, err := vm.Cancel()
	if err != nil {
		t.Fatalf("Cancel: unexpected err: %v", err)
	}
	if refund != 150 {
		t.Fatalf("expected refund of 150, got %d", refund)
	}
	if vm.State().Name() != "Idle" {
		t.Fatalf("expected Idle after cancel, got %s", vm.State().Name())
	}
}

func TestSelectOutOfStockItemRejected(t *testing.T) {
	vm := newTestMachine()

	err := vm.SelectItem("B1") // Water, quantity 0
	if err != ErrOutOfStock {
		t.Fatalf("expected ErrOutOfStock, got %v", err)
	}
	if vm.State().Name() != "OutOfStock" {
		t.Fatalf("expected OutOfStock state, got %s", vm.State().Name())
	}

	// Money can't be inserted while out of stock.
	if err := vm.InsertMoney(100); err != ErrInvalidState {
		t.Fatalf("expected ErrInvalidState, got %v", err)
	}

	// Recover by selecting a different, in-stock item.
	if err := vm.SelectItem("A1"); err != nil {
		t.Fatalf("SelectItem after out-of-stock: unexpected err: %v", err)
	}
	if vm.State().Name() != "Idle" {
		t.Fatalf("expected Idle after re-selecting valid item, got %s", vm.State().Name())
	}
}

func TestSelectInvalidSlot(t *testing.T) {
	vm := newTestMachine()

	if err := vm.SelectItem("Z9"); err != ErrInvalidSlot {
		t.Fatalf("expected ErrInvalidSlot, got %v", err)
	}
}

func TestDispenseDecrementsInventory(t *testing.T) {
	vm := newTestMachine()

	before, ok := vm.Inventory("A1")
	if !ok {
		t.Fatalf("expected slot A1 to exist")
	}
	if before.Quantity != 2 {
		t.Fatalf("expected initial quantity 2, got %d", before.Quantity)
	}

	_ = vm.SelectItem("A1")
	_ = vm.InsertMoney(150)
	if _, err := vm.Dispense(); err != nil {
		t.Fatalf("Dispense: unexpected err: %v", err)
	}

	after, ok := vm.Inventory("A1")
	if !ok {
		t.Fatalf("expected slot A1 to exist")
	}
	if after.Quantity != 1 {
		t.Fatalf("expected quantity decremented to 1, got %d", after.Quantity)
	}
}

func TestInsertMoneyWithoutSelectionRejected(t *testing.T) {
	vm := newTestMachine()

	if err := vm.InsertMoney(100); err != ErrNoSelection {
		t.Fatalf("expected ErrNoSelection, got %v", err)
	}
}

func TestInsertNonPositiveAmountRejected(t *testing.T) {
	vm := newTestMachine()
	_ = vm.SelectItem("A1")

	tests := []int{0, -50}
	for _, amt := range tests {
		if err := vm.InsertMoney(amt); err != ErrInvalidAmount {
			t.Fatalf("InsertMoney(%d): expected ErrInvalidAmount, got %v", amt, err)
		}
	}
}

func TestSequentialPurchasesDrainInventory(t *testing.T) {
	vm := newTestMachine()

	// First soda.
	_ = vm.SelectItem("A1")
	_ = vm.InsertMoney(150)
	if _, err := vm.Dispense(); err != nil {
		t.Fatalf("first dispense: unexpected err: %v", err)
	}

	// Second (and last) soda.
	_ = vm.SelectItem("A1")
	_ = vm.InsertMoney(150)
	if _, err := vm.Dispense(); err != nil {
		t.Fatalf("second dispense: unexpected err: %v", err)
	}

	// Now out of stock.
	err := vm.SelectItem("A1")
	if err != ErrOutOfStock {
		t.Fatalf("expected ErrOutOfStock after inventory drained, got %v", err)
	}
}

func TestSelectItemDuringHasMoneyRejected(t *testing.T) {
	vm := newTestMachine()
	_ = vm.SelectItem("A1")
	_ = vm.InsertMoney(150)

	if err := vm.SelectItem("A2"); err != ErrInvalidState {
		t.Fatalf("expected ErrInvalidState, got %v", err)
	}
}
