package fooddelivery

import (
	"sync"
	"testing"
)

func newTestSystem(partners ...*DeliveryPartner) *FoodDeliverySystem {
	return NewFoodDeliverySystem(partners, NearestAvailablePartnerStrategy{})
}

func testRestaurant() *Restaurant {
	return &Restaurant{
		ID:       "R1",
		Name:     "Tasty Bites",
		Location: Location{X: 0, Y: 0},
		IsOpen:   true,
		Menu: []MenuItem{
			{ID: "I1", Name: "Burger", Price: 5},
			{ID: "I2", Name: "Fries", Price: 2},
		},
	}
}

func TestHappyPathPlaceAssignDeliver(t *testing.T) {
	restaurant := testRestaurant()
	partner := &DeliveryPartner{ID: "P1", Name: "Alex", Location: Location{X: 1, Y: 1}, available: true}
	sys := newTestSystem(partner)
	customer := &Customer{ID: "C1", Name: "Sam"}

	order, err := sys.PlaceOrder(customer, restaurant, []string{"I1", "I2"})
	if err != nil {
		t.Fatalf("expected order placed, got err: %v", err)
	}

	assigned, err := sys.AssignDeliveryPartner(order.ID)
	if err != nil {
		t.Fatalf("expected assignment, got err: %v", err)
	}
	if assigned.ID != "P1" {
		t.Fatalf("expected P1 assigned, got %s", assigned.ID)
	}
	if partner.available {
		t.Fatalf("expected partner marked unavailable after assignment")
	}

	for _, next := range []OrderStatus{Accepted, Preparing, OutForDelivery, Delivered} {
		if err := sys.UpdateOrderStatus(order.ID, next); err != nil {
			t.Fatalf("expected transition to %v to succeed, got err: %v", next, err)
		}
	}

	got, _ := sys.GetOrder(order.ID)
	if got.Status != Delivered {
		t.Fatalf("expected order delivered, got %v", got.Status)
	}
	if !partner.available {
		t.Fatalf("expected partner freed after delivery")
	}
}

func TestInvalidStatusTransitionRejected(t *testing.T) {
	restaurant := testRestaurant()
	partner := &DeliveryPartner{ID: "P1", Name: "Alex", Location: Location{X: 1, Y: 1}, available: true}
	sys := newTestSystem(partner)
	customer := &Customer{ID: "C1", Name: "Sam"}

	order, _ := sys.PlaceOrder(customer, restaurant, []string{"I1"})

	if err := sys.UpdateOrderStatus(order.ID, Delivered); err != ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition jumping straight to Delivered, got %v", err)
	}

	_ = sys.UpdateOrderStatus(order.ID, Accepted)
	_ = sys.UpdateOrderStatus(order.ID, Preparing)
	_ = sys.UpdateOrderStatus(order.ID, OutForDelivery)
	_ = sys.UpdateOrderStatus(order.ID, Delivered)

	if err := sys.UpdateOrderStatus(order.ID, Cancelled); err != ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition cancelling a delivered order, got %v", err)
	}
	if err := sys.UpdateOrderStatus(order.ID, Preparing); err != ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition going back to Preparing from Delivered, got %v", err)
	}
}

func TestItemNotOnMenuAndNoPartnerAvailable(t *testing.T) {
	restaurant := testRestaurant()
	sys := newTestSystem()
	customer := &Customer{ID: "C1", Name: "Sam"}

	_, err := sys.PlaceOrder(customer, restaurant, []string{"BOGUS"})
	if err != ErrItemNotOnMenu {
		t.Fatalf("expected ErrItemNotOnMenu, got %v", err)
	}

	order, err := sys.PlaceOrder(customer, restaurant, []string{"I1"})
	if err != nil {
		t.Fatalf("expected order placed, got err: %v", err)
	}

	_, err = sys.AssignDeliveryPartner(order.ID)
	if err != ErrNoPartnerAvailable {
		t.Fatalf("expected ErrNoPartnerAvailable, got %v", err)
	}
}

// TestConcurrentAssignment asserts two goroutines racing for the single available
// delivery partner never both succeed — the mutex in AssignDeliveryPartner must
// serialize the find-and-mark-unavailable critical section.
func TestConcurrentAssignment(t *testing.T) {
	restaurant := testRestaurant()
	partner := &DeliveryPartner{ID: "P1", Name: "Alex", Location: Location{X: 1, Y: 1}, available: true}
	sys := newTestSystem(partner)
	customer := &Customer{ID: "C1", Name: "Sam"}

	order1, _ := sys.PlaceOrder(customer, restaurant, []string{"I1"})
	order2, _ := sys.PlaceOrder(customer, restaurant, []string{"I2"})

	var wg sync.WaitGroup
	results := make(chan *DeliveryPartner, 2)

	for _, oid := range []string{order1.ID, order2.ID} {
		wg.Add(1)
		go func(orderID string) {
			defer wg.Done()
			p, err := sys.AssignDeliveryPartner(orderID)
			if err == nil {
				results <- p
			} else {
				results <- nil
			}
		}(oid)
	}
	wg.Wait()
	close(results)

	successCount := 0
	for p := range results {
		if p != nil {
			if p.ID != "P1" {
				t.Fatalf("expected the only success to be P1, got %s", p.ID)
			}
			successCount++
		}
	}
	if successCount != 1 {
		t.Fatalf("expected exactly 1 success for single available partner, got %d", successCount)
	}
}
