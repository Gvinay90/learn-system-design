// Package fooddelivery implements a simplified Food Delivery LLD problem:
// order placement, nearest-partner assignment via the Strategy pattern, and
// an order status state machine.
package fooddelivery

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

type OrderStatus int

const (
	Placed OrderStatus = iota
	Accepted
	Preparing
	OutForDelivery
	Delivered
	Cancelled
)

// validNextStatus encodes the allowed forward transitions plus cancellation
// from Placed/Accepted. Anything not listed here is rejected.
var validNextStatus = map[OrderStatus]map[OrderStatus]bool{
	Placed:         {Accepted: true, Cancelled: true},
	Accepted:       {Preparing: true, Cancelled: true},
	Preparing:      {OutForDelivery: true},
	OutForDelivery: {Delivered: true},
	Delivered:      {},
	Cancelled:      {},
}

type Location struct {
	X int
	Y int
}

func (l Location) distanceTo(o Location) float64 {
	dx := float64(l.X - o.X)
	dy := float64(l.Y - o.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

type Customer struct {
	ID   string
	Name string
}

type MenuItem struct {
	ID    string
	Name  string
	Price float64
}

type Restaurant struct {
	ID       string
	Name     string
	Location Location
	Menu     []MenuItem
	IsOpen   bool
}

func (r *Restaurant) findItem(itemID string) (MenuItem, bool) {
	for _, item := range r.Menu {
		if item.ID == itemID {
			return item, true
		}
	}
	return MenuItem{}, false
}

type DeliveryPartner struct {
	ID        string
	Name      string
	Location  Location
	available bool
}

type Order struct {
	ID         string
	Customer   *Customer
	Restaurant *Restaurant
	Items      []MenuItem
	Status     OrderStatus
	Partner    *DeliveryPartner
	PlacedAt   time.Time
}

// AssignmentStrategy picks a delivery partner for an order among the available ones.
type AssignmentStrategy interface {
	Assign(restaurant *Restaurant, partners []*DeliveryPartner) *DeliveryPartner
}

// NearestAvailablePartnerStrategy picks the available partner closest to the restaurant.
type NearestAvailablePartnerStrategy struct{}

func (NearestAvailablePartnerStrategy) Assign(restaurant *Restaurant, partners []*DeliveryPartner) *DeliveryPartner {
	var nearest *DeliveryPartner
	best := math.MaxFloat64
	for _, p := range partners {
		if !p.available {
			continue
		}
		d := restaurant.Location.distanceTo(p.Location)
		if d < best {
			best = d
			nearest = p
		}
	}
	return nearest
}

var (
	ErrRestaurantClosed    = errors.New("restaurant is closed")
	ErrItemNotOnMenu       = errors.New("item not on restaurant menu")
	ErrOrderNotFound       = errors.New("order not found")
	ErrNoPartnerAvailable  = errors.New("no delivery partner available")
	ErrInvalidTransition   = errors.New("invalid order status transition")
	ErrOrderAlreadyAssigned = errors.New("order already has an assigned delivery partner")
)

type FoodDeliverySystem struct {
	Strategy AssignmentStrategy

	mu       sync.Mutex
	partners []*DeliveryPartner
	orders   map[string]*Order
	seq      int
}

func NewFoodDeliverySystem(partners []*DeliveryPartner, strategy AssignmentStrategy) *FoodDeliverySystem {
	return &FoodDeliverySystem{
		Strategy: strategy,
		partners: partners,
		orders:   make(map[string]*Order),
	}
}

// PlaceOrder validates the restaurant is open and every item is on its menu, then creates an order.
func (s *FoodDeliverySystem) PlaceOrder(customer *Customer, restaurant *Restaurant, itemIDs []string) (*Order, error) {
	if !restaurant.IsOpen {
		return nil, ErrRestaurantClosed
	}
	items := make([]MenuItem, 0, len(itemIDs))
	for _, id := range itemIDs {
		item, ok := restaurant.findItem(id)
		if !ok {
			return nil, ErrItemNotOnMenu
		}
		items = append(items, item)
	}

	s.mu.Lock()
	s.seq++
	order := &Order{
		ID:         fmt.Sprintf("O-%d", s.seq),
		Customer:   customer,
		Restaurant: restaurant,
		Items:      items,
		Status:     Placed,
		PlacedAt:   time.Now(),
	}
	s.orders[order.ID] = order
	s.mu.Unlock()
	return order, nil
}

// AssignDeliveryPartner atomically finds the nearest available partner and marks it unavailable,
// so two concurrent assignments can never race onto the same partner.
func (s *FoodDeliverySystem) AssignDeliveryPartner(orderID string) (*DeliveryPartner, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	if order.Partner != nil {
		return nil, ErrOrderAlreadyAssigned
	}

	partner := s.Strategy.Assign(order.Restaurant, s.partners)
	if partner == nil {
		return nil, ErrNoPartnerAvailable
	}
	partner.available = false
	order.Partner = partner
	return partner, nil
}

// UpdateOrderStatus enforces the order lifecycle state machine and frees the delivery
// partner once the order reaches a terminal state.
func (s *FoodDeliverySystem) UpdateOrderStatus(orderID string, next OrderStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return ErrOrderNotFound
	}
	if !validNextStatus[order.Status][next] {
		return ErrInvalidTransition
	}
	order.Status = next
	if (next == Delivered || next == Cancelled) && order.Partner != nil {
		order.Partner.available = true
		order.Partner = nil
	}
	return nil
}

func (s *FoodDeliverySystem) GetOrder(orderID string) (*Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	if !ok {
		return nil, ErrOrderNotFound
	}
	return order, nil
}
