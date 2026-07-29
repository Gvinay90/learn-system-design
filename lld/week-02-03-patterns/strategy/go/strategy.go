// Package strategy demonstrates the Strategy pattern: a ShoppingCart delegates
// price computation to a pluggable PricingStrategy, selected at runtime.
package strategy

// PricingStrategy computes the final price for a given subtotal.
type PricingStrategy interface {
	ApplyDiscount(subtotal float64) float64
}

// RegularPricing applies no discount.
type RegularPricing struct{}

func (RegularPricing) ApplyDiscount(subtotal float64) float64 { return subtotal }

// PercentageDiscountPricing applies a flat percentage off the subtotal.
type PercentageDiscountPricing struct {
	PercentOff float64
}

func (p PercentageDiscountPricing) ApplyDiscount(subtotal float64) float64 {
	return subtotal * (1 - p.PercentOff/100)
}

// ClearancePricing applies a percentage discount plus a flat additional
// amount off, floored at zero.
type ClearancePricing struct {
	PercentOff float64
	FlatOff    float64
}

func (p ClearancePricing) ApplyDiscount(subtotal float64) float64 {
	discounted := subtotal*(1-p.PercentOff/100) - p.FlatOff
	if discounted < 0 {
		return 0
	}
	return discounted
}

// Item is a single line item in the cart.
type Item struct {
	Name     string
	Price    float64
	Quantity int
}

// ShoppingCart holds items and delegates final price computation to whatever
// PricingStrategy it's configured with, so pricing schemes can change at
// runtime without touching cart logic.
type ShoppingCart struct {
	Items    []Item
	Strategy PricingStrategy
}

func NewShoppingCart(strategy PricingStrategy) *ShoppingCart {
	return &ShoppingCart{Strategy: strategy}
}

func (c *ShoppingCart) AddItem(item Item) {
	c.Items = append(c.Items, item)
}

func (c *ShoppingCart) Subtotal() float64 {
	var total float64
	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func (c *ShoppingCart) Checkout() float64 {
	return c.Strategy.ApplyDiscount(c.Subtotal())
}
