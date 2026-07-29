// Package vendingmachine implements a Vending Machine using the State
// pattern: the machine delegates SelectItem/InsertMoney/Dispense/Cancel to
// its current VendingState object, which decides what happens (if anything)
// from that state. This mirrors the style used in
// lld/week-02-03-patterns/state: a small interface implemented by one
// zero-value struct per state, with the context forwarding every action.
package vendingmachine

import "errors"

var (
	// ErrInvalidSlot is returned when SelectItem names a slot that doesn't exist.
	ErrInvalidSlot = errors.New("vendingmachine: invalid slot")
	// ErrOutOfStock is returned when SelectItem names a slot with zero quantity.
	ErrOutOfStock = errors.New("vendingmachine: item out of stock")
	// ErrNoSelection is returned when money is inserted before an item is selected.
	ErrNoSelection = errors.New("vendingmachine: no item selected")
	// ErrInvalidAmount is returned when InsertMoney is called with a non-positive amount.
	ErrInvalidAmount = errors.New("vendingmachine: amount must be positive")
	// ErrNotEnoughMoney is returned when Dispense is attempted before enough money was inserted.
	ErrNotEnoughMoney = errors.New("vendingmachine: not enough money inserted")
	// ErrInvalidState is returned when an action isn't legal from the machine's current state.
	ErrInvalidState = errors.New("vendingmachine: action not allowed in current state")
)

// Slot holds the inventory for a single selectable position in the machine.
type Slot struct {
	Item     string
	Price    int // in cents
	Quantity int
}

// DispenseResult is returned by a successful Dispense call.
type DispenseResult struct {
	Item   string
	Change int // in cents
}

// VendingState is the interface every concrete state implements. Each
// method mutates the machine to the next state on success, or returns an
// error (without changing state) if the action isn't legal from that state.
type VendingState interface {
	Name() string
	SelectItem(vm *VendingMachine, slotID string) error
	InsertMoney(vm *VendingMachine, amount int) error
	Dispense(vm *VendingMachine) (DispenseResult, error)
	Cancel(vm *VendingMachine) (refund int, err error)
}

// VendingMachine is the context: it holds the current state, the inventory,
// the currently selected slot, and the money collected so far, and forwards
// every action to its current state.
type VendingMachine struct {
	slots    map[string]*Slot
	state    VendingState
	selected string
	balance  int // in cents
}

// NewVendingMachine builds a machine stocked with the given slots, keyed by
// slot ID (e.g. "A1"). It starts in the Idle state.
func NewVendingMachine(slots map[string]Slot) *VendingMachine {
	inventory := make(map[string]*Slot, len(slots))
	for id, s := range slots {
		copy := s
		inventory[id] = &copy
	}
	return &VendingMachine{slots: inventory, state: IdleState{}}
}

// State returns the machine's current state object.
func (vm *VendingMachine) State() VendingState { return vm.state }

// Balance returns the money collected so far for the in-progress transaction, in cents.
func (vm *VendingMachine) Balance() int { return vm.balance }

// Selected returns the currently selected slot ID, or "" if none.
func (vm *VendingMachine) Selected() string { return vm.selected }

// Inventory returns a copy of the slot at slotID and whether it exists.
func (vm *VendingMachine) Inventory(slotID string) (Slot, bool) {
	s, ok := vm.slots[slotID]
	if !ok {
		return Slot{}, false
	}
	return *s, true
}

// SelectItem picks a slot to purchase from.
func (vm *VendingMachine) SelectItem(slotID string) error {
	return vm.state.SelectItem(vm, slotID)
}

// InsertMoney adds amount (in cents) toward the selected item's price.
func (vm *VendingMachine) InsertMoney(amount int) error {
	return vm.state.InsertMoney(vm, amount)
}

// Dispense hands over the item once enough money has been inserted,
// decrementing inventory and returning any change owed.
func (vm *VendingMachine) Dispense() (DispenseResult, error) {
	return vm.state.Dispense(vm)
}

// Cancel aborts the in-progress transaction, refunding any money inserted.
func (vm *VendingMachine) Cancel() (int, error) {
	return vm.state.Cancel(vm)
}

func (vm *VendingMachine) setState(s VendingState) { vm.state = s }

func (vm *VendingMachine) slot(slotID string) (*Slot, error) {
	s, ok := vm.slots[slotID]
	if !ok {
		return nil, ErrInvalidSlot
	}
	return s, nil
}

// selectItem is the shared SelectItem logic used by both IdleState and
// OutOfStockState: validate the slot, and either settle on it (staying in
// Idle, awaiting money) or bounce to OutOfStockState if it's empty.
func selectItem(vm *VendingMachine, slotID string) error {
	s, err := vm.slot(slotID)
	if err != nil {
		return err
	}
	if s.Quantity <= 0 {
		vm.selected = slotID
		vm.setState(OutOfStockState{})
		return ErrOutOfStock
	}
	vm.selected = slotID
	vm.setState(IdleState{})
	return nil
}

// IdleState: no completed transaction in progress. An item may or may not
// be selected yet, and partial money may have been inserted; once inserted
// money covers the selected item's price the machine moves to HasMoney.
type IdleState struct{}

func (IdleState) Name() string { return "Idle" }

func (IdleState) SelectItem(vm *VendingMachine, slotID string) error {
	return selectItem(vm, slotID)
}

func (IdleState) InsertMoney(vm *VendingMachine, amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if vm.selected == "" {
		return ErrNoSelection
	}
	s, err := vm.slot(vm.selected)
	if err != nil {
		return err
	}
	vm.balance += amount
	if vm.balance >= s.Price {
		vm.setState(HasMoneyState{})
	}
	return nil
}

func (IdleState) Dispense(vm *VendingMachine) (DispenseResult, error) {
	return DispenseResult{}, ErrNotEnoughMoney
}

func (IdleState) Cancel(vm *VendingMachine) (int, error) {
	refund := vm.balance
	vm.balance = 0
	vm.selected = ""
	return refund, nil
}

// HasMoneyState: enough money has been inserted to cover the selected
// item's price. Further top-ups are still accepted; Dispense completes the
// purchase, and Cancel refunds everything inserted.
type HasMoneyState struct{}

func (HasMoneyState) Name() string { return "HasMoney" }

func (HasMoneyState) SelectItem(vm *VendingMachine, slotID string) error {
	return ErrInvalidState
}

func (HasMoneyState) InsertMoney(vm *VendingMachine, amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	vm.balance += amount
	return nil
}

func (HasMoneyState) Dispense(vm *VendingMachine) (DispenseResult, error) {
	vm.setState(DispensingState{})
	return vm.state.Dispense(vm)
}

func (HasMoneyState) Cancel(vm *VendingMachine) (int, error) {
	refund := vm.balance
	vm.balance = 0
	vm.selected = ""
	vm.setState(IdleState{})
	return refund, nil
}

// DispensingState: the machine is actively completing a purchase. This is
// entered only transiently from HasMoney.Dispense, which performs the
// inventory decrement and change calculation and then falls back to Idle;
// it rejects every other action as busy.
type DispensingState struct{}

func (DispensingState) Name() string { return "Dispensing" }

func (DispensingState) SelectItem(vm *VendingMachine, slotID string) error {
	return ErrInvalidState
}

func (DispensingState) InsertMoney(vm *VendingMachine, amount int) error {
	return ErrInvalidState
}

func (DispensingState) Dispense(vm *VendingMachine) (DispenseResult, error) {
	s, err := vm.slot(vm.selected)
	if err != nil {
		vm.setState(IdleState{})
		return DispenseResult{}, err
	}
	if vm.balance < s.Price {
		vm.setState(HasMoneyState{})
		return DispenseResult{}, ErrNotEnoughMoney
	}

	change := vm.balance - s.Price
	s.Quantity--
	item := s.Item

	vm.balance = 0
	vm.selected = ""
	vm.setState(IdleState{})

	return DispenseResult{Item: item, Change: change}, nil
}

func (DispensingState) Cancel(vm *VendingMachine) (int, error) {
	return 0, ErrInvalidState
}

// OutOfStockState: the last SelectItem named a slot with zero quantity.
// The machine holds no money in this state (nothing was accepted for an
// out-of-stock item), so Cancel simply returns to Idle; SelectItem may be
// retried with a different (in-stock) slot.
type OutOfStockState struct{}

func (OutOfStockState) Name() string { return "OutOfStock" }

func (OutOfStockState) SelectItem(vm *VendingMachine, slotID string) error {
	return selectItem(vm, slotID)
}

func (OutOfStockState) InsertMoney(vm *VendingMachine, amount int) error {
	return ErrInvalidState
}

func (OutOfStockState) Dispense(vm *VendingMachine) (DispenseResult, error) {
	return DispenseResult{}, ErrInvalidState
}

func (OutOfStockState) Cancel(vm *VendingMachine) (int, error) {
	refund := vm.balance
	vm.balance = 0
	vm.selected = ""
	vm.setState(IdleState{})
	return refund, nil
}
