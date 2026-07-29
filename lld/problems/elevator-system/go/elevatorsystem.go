// Package elevatorsystem implements the classic Elevator System LLD problem:
// N floors, M elevator cars, hall calls + cabin requests, and a pluggable
// dispatch/scheduling Strategy that picks the best car for a hall call.
package elevatorsystem

import (
	"sort"
	"sync"
)

type Direction int

const (
	Idle Direction = iota
	Up
	Down
)

type CarState int

const (
	StateIdle CarState = iota
	StateMovingUp
	StateMovingDown
	StateDoorOpen
)

// Car represents a single elevator cabin. It tracks its own current floor,
// motion state, direction and the set of floors it still needs to visit.
// A per-car mutex guards all of this so that hall-call dispatch (from the
// system) and cabin requests (from inside the car) can be issued concurrently
// without corrupting the target set or state.
type Car struct {
	ID           int
	NumFloors    int
	CurrentFloor int

	mu        sync.Mutex
	state     CarState
	direction Direction
	targets   map[int]bool
}

func NewCar(id, numFloors, startFloor int) *Car {
	return &Car{
		ID:           id,
		NumFloors:    numFloors,
		CurrentFloor: startFloor,
		state:        StateIdle,
		direction:    Idle,
		targets:      make(map[int]bool),
	}
}

func (c *Car) Snapshot() (floor int, state CarState, dir Direction, pending int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.CurrentFloor, c.state, c.direction, len(c.targets)
}

// AddTarget registers a floor the car must visit (from a hall call assigned
// to this car, or a cabin request made from inside it). If the car is
// currently idle, its direction is derived from the new target.
func (c *Car) AddTarget(floor int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.targets[floor] = true
	if c.state == StateIdle {
		c.setDirectionTowards(floor)
	}
}

func (c *Car) setDirectionTowards(floor int) {
	switch {
	case floor > c.CurrentFloor:
		c.direction = Up
		c.state = StateMovingUp
	case floor < c.CurrentFloor:
		c.direction = Down
		c.state = StateMovingDown
	default:
		// The target is the floor we're already sitting on: service it now,
		// otherwise it would never be removed and pickNextDirection would
		// keep re-selecting it as "nearest" forever, starving every other
		// pending target on this car.
		delete(c.targets, floor)
		c.state = StateDoorOpen
	}
}

// Step advances simulated time by one tick: a car in DoorOpen closes its
// doors and picks a new direction (or goes idle); otherwise it moves one
// floor towards its next target, servicing (removing) any target it lands on.
func (c *Car) Step() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == StateDoorOpen {
		c.pickNextDirection()
		return
	}
	if len(c.targets) == 0 {
		c.state = StateIdle
		c.direction = Idle
		return
	}

	switch c.direction {
	case Up:
		c.CurrentFloor++
	case Down:
		c.CurrentFloor--
	default:
		c.pickNextDirection()
		return
	}

	if c.targets[c.CurrentFloor] {
		delete(c.targets, c.CurrentFloor)
		c.state = StateDoorOpen
	}
}

// pickNextDirection chooses the direction of the nearest remaining target,
// simplified SCAN behaviour: continue in the current direction if a target
// still lies ahead, otherwise reverse, otherwise go idle.
func (c *Car) pickNextDirection() {
	if len(c.targets) == 0 {
		c.state = StateIdle
		c.direction = Idle
		return
	}
	floors := make([]int, 0, len(c.targets))
	for f := range c.targets {
		floors = append(floors, f)
	}
	sort.Ints(floors)

	nearest := floors[0]
	best := abs(nearest - c.CurrentFloor)
	for _, f := range floors {
		if d := abs(f - c.CurrentFloor); d < best {
			best = d
			nearest = f
		}
	}
	c.setDirectionTowards(nearest)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// DispatchStrategy picks the best car in the fleet to service a hall call.
type DispatchStrategy interface {
	SelectCar(cars []*Car, floor int, dir Direction) *Car
}

// NearestCarStrategy prefers an idle car nearest to the call floor; failing
// that, a car already moving in the requested direction that hasn't yet
// passed the call floor (so it can pick it up en route); failing that, the
// overall nearest car regardless of direction compatibility.
type NearestCarStrategy struct{}

func (NearestCarStrategy) SelectCar(cars []*Car, floor int, dir Direction) *Car {
	var best *Car
	bestCost := -1

	for _, c := range cars {
		f, state, d, _ := c.Snapshot()
		cost := abs(f - floor)

		compatible := state == StateIdle ||
			(d == dir && dir == Up && f <= floor) ||
			(d == dir && dir == Down && f >= floor)

		if !compatible {
			cost += c.NumFloors // penalize incompatible cars so idle/en-route cars win when available
		}

		if best == nil || cost < bestCost {
			best = c
			bestCost = cost
		}
	}
	return best
}

// ElevatorSystem coordinates a fleet of cars serving a building with
// NumFloors floors, dispatching hall calls via a pluggable Strategy.
type ElevatorSystem struct {
	Cars      []*Car
	NumFloors int
	Strategy  DispatchStrategy

	// mu guards the select-then-assign sequence for hall calls so that two
	// concurrent hall calls can never both "win" the same decision based on
	// a stale view of car state (analogous to Floor's mutex in parking-lot).
	mu sync.Mutex
}

func NewElevatorSystem(numFloors int, cars []*Car, strategy DispatchStrategy) *ElevatorSystem {
	return &ElevatorSystem{Cars: cars, NumFloors: numFloors, Strategy: strategy}
}

// HallCall handles an external request for a car at `floor` heading `dir`,
// dispatching it to the strategy's chosen car and returning that car.
func (s *ElevatorSystem) HallCall(floor int, dir Direction) *Car {
	s.mu.Lock()
	car := s.Strategy.SelectCar(s.Cars, floor, dir)
	defer s.mu.Unlock()
	car.AddTarget(floor)
	return car
}

// CabinRequest handles an internal request made from inside a specific car.
func (s *ElevatorSystem) CabinRequest(car *Car, destinationFloor int) {
	car.AddTarget(destinationFloor)
}

// Step advances every car in the fleet by one simulated tick.
func (s *ElevatorSystem) Step() {
	for _, c := range s.Cars {
		c.Step()
	}
}
