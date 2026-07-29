package state

import "testing"

func TestHappyPathTransitions(t *testing.T) {
	o := NewOrder("O-1")
	if o.State().Name() != "Created" {
		t.Fatalf("expected Created, got %s", o.State().Name())
	}

	if err := o.Pay(); err != nil {
		t.Fatalf("unexpected err on Pay: %v", err)
	}
	if o.State().Name() != "Paid" {
		t.Fatalf("expected Paid, got %s", o.State().Name())
	}

	if err := o.Ship(); err != nil {
		t.Fatalf("unexpected err on Ship: %v", err)
	}
	if o.State().Name() != "Shipped" {
		t.Fatalf("expected Shipped, got %s", o.State().Name())
	}

	if err := o.Deliver(); err != nil {
		t.Fatalf("unexpected err on Deliver: %v", err)
	}
	if o.State().Name() != "Delivered" {
		t.Fatalf("expected Delivered, got %s", o.State().Name())
	}
}

func TestCannotShipBeforePayment(t *testing.T) {
	o := NewOrder("O-2")
	if err := o.Ship(); err != ErrIllegalTransition {
		t.Fatalf("expected ErrIllegalTransition, got %v", err)
	}
	if o.State().Name() != "Created" {
		t.Fatalf("expected state unchanged at Created, got %s", o.State().Name())
	}
}

func TestCancelAllowedBeforeShipping(t *testing.T) {
	o := NewOrder("O-3")
	_ = o.Pay()
	if err := o.Cancel(); err != nil {
		t.Fatalf("unexpected err on Cancel: %v", err)
	}
	if o.State().Name() != "Cancelled" {
		t.Fatalf("expected Cancelled, got %s", o.State().Name())
	}
}

func TestCannotCancelAfterShipping(t *testing.T) {
	o := NewOrder("O-4")
	_ = o.Pay()
	_ = o.Ship()
	if err := o.Cancel(); err != ErrIllegalTransition {
		t.Fatalf("expected ErrIllegalTransition, got %v", err)
	}
	if o.State().Name() != "Shipped" {
		t.Fatalf("expected state unchanged at Shipped, got %s", o.State().Name())
	}
}

func TestDeliveredIsTerminal(t *testing.T) {
	o := NewOrder("O-5")
	_ = o.Pay()
	_ = o.Ship()
	_ = o.Deliver()

	if err := o.Cancel(); err != ErrIllegalTransition {
		t.Fatalf("expected ErrIllegalTransition on Cancel after Delivered, got %v", err)
	}
	if err := o.Pay(); err != ErrIllegalTransition {
		t.Fatalf("expected ErrIllegalTransition on Pay after Delivered, got %v", err)
	}
}
