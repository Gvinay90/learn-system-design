// Package state demonstrates the State pattern: an Order delegates its
// behavior (Pay, Ship, Deliver, Cancel) to its current OrderState object,
// which decides what transition, if any, is legal from that state.
package state

import "errors"

var ErrIllegalTransition = errors.New("illegal transition for current order state")

// OrderState is the interface every concrete state implements. Each method
// returns the next state on success, or ErrIllegalTransition if the
// transition isn't allowed from that state.
type OrderState interface {
	Name() string
	Pay(o *Order) error
	Ship(o *Order) error
	Deliver(o *Order) error
	Cancel(o *Order) error
}

// Order is the context: it holds a reference to its current state and
// forwards every action to it.
type Order struct {
	ID    string
	state OrderState
}

func NewOrder(id string) *Order {
	return &Order{ID: id, state: CreatedState{}}
}

func (o *Order) State() OrderState { return o.state }

func (o *Order) Pay() error     { return o.transition(o.state.Pay) }
func (o *Order) Ship() error    { return o.transition(o.state.Ship) }
func (o *Order) Deliver() error { return o.transition(o.state.Deliver) }
func (o *Order) Cancel() error  { return o.transition(o.state.Cancel) }

func (o *Order) transition(action func(*Order) error) error {
	return action(o)
}

func (o *Order) setState(s OrderState) { o.state = s }

// CreatedState: order placed, awaiting payment. Can pay or cancel.
type CreatedState struct{}

func (CreatedState) Name() string { return "Created" }
func (CreatedState) Pay(o *Order) error {
	o.setState(PaidState{})
	return nil
}
func (CreatedState) Ship(o *Order) error    { return ErrIllegalTransition }
func (CreatedState) Deliver(o *Order) error { return ErrIllegalTransition }
func (CreatedState) Cancel(o *Order) error {
	o.setState(CancelledState{})
	return nil
}

// PaidState: payment captured, awaiting shipment. Can ship or cancel (refund).
type PaidState struct{}

func (PaidState) Name() string       { return "Paid" }
func (PaidState) Pay(o *Order) error { return ErrIllegalTransition }
func (PaidState) Ship(o *Order) error {
	o.setState(ShippedState{})
	return nil
}
func (PaidState) Deliver(o *Order) error { return ErrIllegalTransition }
func (PaidState) Cancel(o *Order) error {
	o.setState(CancelledState{})
	return nil
}

// ShippedState: in transit. Can only be delivered; cannot cancel once shipped.
type ShippedState struct{}

func (ShippedState) Name() string        { return "Shipped" }
func (ShippedState) Pay(o *Order) error  { return ErrIllegalTransition }
func (ShippedState) Ship(o *Order) error { return ErrIllegalTransition }
func (ShippedState) Deliver(o *Order) error {
	o.setState(DeliveredState{})
	return nil
}
func (ShippedState) Cancel(o *Order) error { return ErrIllegalTransition }

// DeliveredState: terminal success state. No further transitions.
type DeliveredState struct{}

func (DeliveredState) Name() string           { return "Delivered" }
func (DeliveredState) Pay(o *Order) error     { return ErrIllegalTransition }
func (DeliveredState) Ship(o *Order) error    { return ErrIllegalTransition }
func (DeliveredState) Deliver(o *Order) error { return ErrIllegalTransition }
func (DeliveredState) Cancel(o *Order) error  { return ErrIllegalTransition }

// CancelledState: terminal state. No further transitions.
type CancelledState struct{}

func (CancelledState) Name() string           { return "Cancelled" }
func (CancelledState) Pay(o *Order) error     { return ErrIllegalTransition }
func (CancelledState) Ship(o *Order) error    { return ErrIllegalTransition }
func (CancelledState) Deliver(o *Order) error { return ErrIllegalTransition }
func (CancelledState) Cancel(o *Order) error  { return ErrIllegalTransition }
