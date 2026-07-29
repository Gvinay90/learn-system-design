// Package ridesharing implements the classic Uber-style Ride Sharing LLD problem:
// rider/driver matching (Strategy), trip lifecycle, and pluggable pricing.
package ridesharing

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"
)

type Location struct {
	X, Y float64
}

func distance(a, b Location) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type Rider struct {
	ID   string
	Name string
}

type Driver struct {
	ID        string
	Name      string
	Location  Location
	available bool
}

type TripStatus int

const (
	Requested TripStatus = iota
	Accepted
	InProgress
	Completed
	Cancelled
)

type Trip struct {
	ID          string
	Rider       *Rider
	Driver      *Driver
	Pickup      Location
	Dropoff     Location
	Status      TripStatus
	RequestedAt time.Time
	Fare        float64
}

// DriverMatchingStrategy picks a driver for a trip from the available pool.
type DriverMatchingStrategy interface {
	Match(pickup Location, drivers []*Driver) *Driver
}

// NearestAvailableDriverStrategy picks the closest available driver to the pickup point.
type NearestAvailableDriverStrategy struct{}

func (NearestAvailableDriverStrategy) Match(pickup Location, drivers []*Driver) *Driver {
	var best *Driver
	bestDist := math.MaxFloat64
	for _, d := range drivers {
		if !d.available {
			continue
		}
		dist := distance(pickup, d.Location)
		if dist < bestDist {
			bestDist = dist
			best = d
		}
	}
	return best
}

// PricingStrategy computes the fare for a completed trip.
type PricingStrategy interface {
	CalculateFare(t *Trip) float64
}

// DistanceBasedPricing charges a base fare plus a per-unit-distance rate.
type DistanceBasedPricing struct {
	BaseFare    float64
	PerDistance float64
}

func (p DistanceBasedPricing) CalculateFare(t *Trip) float64 {
	return p.BaseFare + p.PerDistance*distance(t.Pickup, t.Dropoff)
}

var (
	ErrTripNotFound      = errors.New("trip not recognized")
	ErrNoDriverAvailable = errors.New("no driver available")
	ErrInvalidTransition = errors.New("invalid trip status transition")
)

type RideSharingSystem struct {
	Matching DriverMatchingStrategy
	Pricing  PricingStrategy

	mu      sync.Mutex
	drivers []*Driver
	trips   map[string]*Trip
	seq     int
}

func NewRideSharingSystem(drivers []*Driver, matching DriverMatchingStrategy, pricing PricingStrategy) *RideSharingSystem {
	return &RideSharingSystem{
		Matching: matching,
		Pricing:  pricing,
		drivers:  drivers,
		trips:    make(map[string]*Trip),
	}
}

func (s *RideSharingSystem) RequestRide(rider *Rider, pickup, dropoff Location) *Trip {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	trip := &Trip{
		ID:          fmt.Sprintf("TR-%d", s.seq),
		Rider:       rider,
		Pickup:      pickup,
		Dropoff:     dropoff,
		Status:      Requested,
		RequestedAt: time.Now(),
	}
	s.trips[trip.ID] = trip
	return trip
}

// MatchDriver finds the nearest available driver and assigns it atomically:
// the whole find-and-mark-unavailable step is guarded by the system mutex so
// two concurrent calls can never both be assigned the same single driver.
func (s *RideSharingSystem) MatchDriver(tripID string) (*Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	trip, ok := s.trips[tripID]
	if !ok {
		return nil, ErrTripNotFound
	}
	if trip.Status != Requested {
		return nil, ErrInvalidTransition
	}

	driver := s.Matching.Match(trip.Pickup, s.drivers)
	if driver == nil {
		return nil, ErrNoDriverAvailable
	}

	driver.available = false
	trip.Driver = driver
	trip.Status = Accepted
	return driver, nil
}

func (s *RideSharingSystem) StartTrip(tripID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	trip, ok := s.trips[tripID]
	if !ok {
		return ErrTripNotFound
	}
	if trip.Status != Accepted {
		return ErrInvalidTransition
	}
	trip.Status = InProgress
	return nil
}

func (s *RideSharingSystem) CompleteTrip(tripID string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	trip, ok := s.trips[tripID]
	if !ok {
		return 0, ErrTripNotFound
	}
	if trip.Status != InProgress {
		return 0, ErrInvalidTransition
	}
	trip.Status = Completed
	trip.Fare = s.Pricing.CalculateFare(trip)
	if trip.Driver != nil {
		trip.Driver.available = true
	}
	return trip.Fare, nil
}

func (s *RideSharingSystem) CancelTrip(tripID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	trip, ok := s.trips[tripID]
	if !ok {
		return ErrTripNotFound
	}
	if trip.Status != Requested && trip.Status != Accepted {
		return ErrInvalidTransition
	}
	trip.Status = Cancelled
	if trip.Driver != nil {
		trip.Driver.available = true
	}
	return nil
}
