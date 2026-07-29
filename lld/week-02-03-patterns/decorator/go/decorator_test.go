package decorator

import (
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestBaseCoffee(t *testing.T) {
	var c Coffee = Espresso{}
	if !almostEqual(c.Cost(), 2.0) {
		t.Fatalf("expected cost 2.0, got %v", c.Cost())
	}
	if c.Description() != "Espresso" {
		t.Fatalf("expected 'Espresso', got %q", c.Description())
	}
}

func TestSingleDecorator(t *testing.T) {
	var c Coffee = MilkDecorator{Wrapped: Espresso{}}
	if !almostEqual(c.Cost(), 2.5) {
		t.Fatalf("expected cost 2.5, got %v", c.Cost())
	}
	if c.Description() != "Espresso + Milk" {
		t.Fatalf("unexpected description: %q", c.Description())
	}
}

func TestStackedDecoratorsCumulativeCostAndDescription(t *testing.T) {
	var c Coffee = Espresso{}
	c = MilkDecorator{Wrapped: c}
	c = SugarDecorator{Wrapped: c}
	c = WhipDecorator{Wrapped: c}

	wantCost := 2.0 + 0.5 + 0.25 + 0.75
	if !almostEqual(c.Cost(), wantCost) {
		t.Fatalf("expected cost %v, got %v", wantCost, c.Cost())
	}
	wantDesc := "Espresso + Milk + Sugar + Whip"
	if c.Description() != wantDesc {
		t.Fatalf("expected %q, got %q", wantDesc, c.Description())
	}
}

func TestDecoratorOrderIndependenceOfCost(t *testing.T) {
	var a Coffee = WhipDecorator{Wrapped: MilkDecorator{Wrapped: Espresso{}}}
	var b Coffee = MilkDecorator{Wrapped: WhipDecorator{Wrapped: Espresso{}}}
	if !almostEqual(a.Cost(), b.Cost()) {
		t.Fatalf("cost should be order-independent: %v vs %v", a.Cost(), b.Cost())
	}
	if a.Description() == b.Description() {
		t.Fatalf("description should differ by wrap order, got same: %q", a.Description())
	}
}

func TestDemo(t *testing.T) {
	var c Coffee = Espresso{}
	c = MilkDecorator{Wrapped: c}
	c = SugarDecorator{Wrapped: c}
	t.Logf("%s costs %.2f", c.Description(), c.Cost())
}
