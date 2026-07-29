package strategy

import "testing"

func TestRegularPricingNoDiscount(t *testing.T) {
	cart := NewShoppingCart(RegularPricing{})
	cart.AddItem(Item{Name: "book", Price: 20, Quantity: 2})

	if got := cart.Checkout(); got != 40 {
		t.Fatalf("expected 40, got %v", got)
	}
}

func TestPercentageDiscountPricing(t *testing.T) {
	cart := NewShoppingCart(PercentageDiscountPricing{PercentOff: 10})
	cart.AddItem(Item{Name: "shoes", Price: 100, Quantity: 1})

	if got := cart.Checkout(); got != 90 {
		t.Fatalf("expected 90, got %v", got)
	}
}

func TestClearancePricing(t *testing.T) {
	cart := NewShoppingCart(ClearancePricing{PercentOff: 20, FlatOff: 5})
	cart.AddItem(Item{Name: "jacket", Price: 100, Quantity: 1})

	// 100 * 0.8 = 80, minus flat 5 = 75
	if got := cart.Checkout(); got != 75 {
		t.Fatalf("expected 75, got %v", got)
	}
}

func TestClearancePricingFloorsAtZero(t *testing.T) {
	cart := NewShoppingCart(ClearancePricing{PercentOff: 50, FlatOff: 100})
	cart.AddItem(Item{Name: "sticker", Price: 5, Quantity: 1})

	if got := cart.Checkout(); got != 0 {
		t.Fatalf("expected floored at 0, got %v", got)
	}
}

func TestSwitchingStrategyAtRuntime(t *testing.T) {
	cart := NewShoppingCart(RegularPricing{})
	cart.AddItem(Item{Name: "widget", Price: 50, Quantity: 2})

	if got := cart.Checkout(); got != 100 {
		t.Fatalf("expected 100 with regular pricing, got %v", got)
	}

	cart.Strategy = PercentageDiscountPricing{PercentOff: 25}
	if got := cart.Checkout(); got != 75 {
		t.Fatalf("expected 75 after switching strategy, got %v", got)
	}
}
