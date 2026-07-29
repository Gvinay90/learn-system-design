// Package decorator demonstrates the Decorator structural pattern:
// a base Coffee is wrapped by add-on decorators (milk, sugar, whip) that
// each layer in additional cost and description without subclass explosion.
package decorator

// Coffee is the component interface shared by the base drink and all decorators.
type Coffee interface {
	Cost() float64
	Description() string
}

type Espresso struct{}

func (Espresso) Cost() float64        { return 2.0 }
func (Espresso) Description() string  { return "Espresso" }

// MilkDecorator wraps a Coffee, adding its own cost/description on top.
type MilkDecorator struct {
	Wrapped Coffee
}

func (d MilkDecorator) Cost() float64 { return d.Wrapped.Cost() + 0.5 }
func (d MilkDecorator) Description() string {
	return d.Wrapped.Description() + " + Milk"
}

type SugarDecorator struct {
	Wrapped Coffee
}

func (d SugarDecorator) Cost() float64 { return d.Wrapped.Cost() + 0.25 }
func (d SugarDecorator) Description() string {
	return d.Wrapped.Description() + " + Sugar"
}

type WhipDecorator struct {
	Wrapped Coffee
}

func (d WhipDecorator) Cost() float64 { return d.Wrapped.Cost() + 0.75 }
func (d WhipDecorator) Description() string {
	return d.Wrapped.Description() + " + Whip"
}
