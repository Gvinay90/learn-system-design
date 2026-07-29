// Package parkinglot implements the classic Parking Lot LLD problem:
// multiple floors, multiple spot sizes, nearest-spot allocation, and
// pluggable pricing via the Strategy pattern.
package parkinglot

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type VehicleType int

const (
	Motorcycle VehicleType = iota
	Compact
	Large
)

type Vehicle struct {
	LicensePlate string
	Type         VehicleType
}

type ParkingSpot struct {
	ID       string
	Type     VehicleType
	occupied bool
	vehicle  *Vehicle
}

// fits reports whether this spot can hold the given vehicle type.
// A vehicle can only park in a spot of its own size class in this simplified model.
func (s *ParkingSpot) fits(vt VehicleType) bool {
	return s.Type == vt
}

type Floor struct {
	Level int
	Spots []*ParkingSpot
	mu    sync.Mutex
}

// findAndAssign atomically finds the first available spot fitting vt and marks it occupied.
func (f *Floor) findAndAssign(v *Vehicle) *ParkingSpot {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, s := range f.Spots {
		if !s.occupied && s.fits(v.Type) {
			s.occupied = true
			s.vehicle = v
			return s
		}
	}
	return nil
}

func (f *Floor) free(spotID string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, s := range f.Spots {
		if s.ID == spotID {
			s.occupied = false
			s.vehicle = nil
			return
		}
	}
}

type Ticket struct {
	ID        string
	Vehicle   *Vehicle
	Spot      *ParkingSpot
	FloorNum  int
	EntryTime time.Time
}

// PricingStrategy computes a fee for a completed parking session.
type PricingStrategy interface {
	CalculateFee(t *Ticket, exitTime time.Time) float64
}

// HourlyTieredPricing charges a flat rate for the first hour, a lower rate after.
type HourlyTieredPricing struct {
	FirstHourRate float64
	SubsequentRate float64
}

func (p HourlyTieredPricing) CalculateFee(t *Ticket, exitTime time.Time) float64 {
	duration := exitTime.Sub(t.EntryTime)
	hours := duration.Hours()
	if hours <= 1 {
		return p.FirstHourRate
	}
	return p.FirstHourRate + (hours-1)*p.SubsequentRate
}

var ErrLotFull = errors.New("no available spot for vehicle type")
var ErrTicketNotFound = errors.New("ticket not recognized")

type ParkingLot struct {
	Floors  []*Floor
	Pricing PricingStrategy

	mu      sync.Mutex
	tickets map[string]*Ticket
	seq     int
}

func NewParkingLot(floors []*Floor, pricing PricingStrategy) *ParkingLot {
	return &ParkingLot{Floors: floors, Pricing: pricing, tickets: make(map[string]*Ticket)}
}

// ParkVehicle finds the nearest available spot (lowest floor first) and issues a ticket.
func (l *ParkingLot) ParkVehicle(v *Vehicle) (*Ticket, error) {
	for _, floor := range l.Floors {
		if spot := floor.findAndAssign(v); spot != nil {
			l.mu.Lock()
			l.seq++
			ticket := &Ticket{
				ID:        fmt.Sprintf("T-%d", l.seq),
				Vehicle:   v,
				Spot:      spot,
				FloorNum:  floor.Level,
				EntryTime: time.Now(),
			}
			l.tickets[ticket.ID] = ticket
			l.mu.Unlock()
			return ticket, nil
		}
	}
	return nil, ErrLotFull
}

// UnparkVehicle frees the spot tied to ticketID and returns the computed fee.
func (l *ParkingLot) UnparkVehicle(ticketID string, exitTime time.Time) (float64, error) {
	l.mu.Lock()
	ticket, ok := l.tickets[ticketID]
	if ok {
		delete(l.tickets, ticketID)
	}
	l.mu.Unlock()
	if !ok {
		return 0, ErrTicketNotFound
	}

	for _, floor := range l.Floors {
		if floor.Level == ticket.FloorNum {
			floor.free(ticket.Spot.ID)
			break
		}
	}
	return l.Pricing.CalculateFee(ticket, exitTime), nil
}
