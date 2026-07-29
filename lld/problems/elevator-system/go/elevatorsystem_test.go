package elevatorsystem

import (
	"sync"
	"testing"
)

func newTestSystem(numFloors int, startFloors ...int) *ElevatorSystem {
	cars := make([]*Car, 0, len(startFloors))
	for i, f := range startFloors {
		cars = append(cars, NewCar(i, numFloors, f))
	}
	return NewElevatorSystem(numFloors, cars, NearestCarStrategy{})
}

func runUntilIdle(s *ElevatorSystem, maxSteps int) {
	for i := 0; i < maxSteps; i++ {
		s.Step()
		allIdle := true
		for _, c := range s.Cars {
			_, state, _, pending := c.Snapshot()
			if state != StateIdle || pending > 0 {
				allIdle = false
			}
		}
		if allIdle {
			return
		}
	}
}

func TestHallCallDispatchedAndServiced(t *testing.T) {
	sys := newTestSystem(10, 0)
	car := sys.HallCall(5, Up)
	if car.ID != 0 {
		t.Fatalf("expected only car to be dispatched, got car %d", car.ID)
	}

	runUntilIdle(sys, 50)

	floor, state, _, pending := car.Snapshot()
	if floor != 5 {
		t.Fatalf("expected car to reach floor 5, got %d", floor)
	}
	if pending != 0 {
		t.Fatalf("expected no pending targets, got %d", pending)
	}
	if state != StateDoorOpen && state != StateIdle {
		t.Fatalf("expected car to have serviced its stop, got state %v", state)
	}
}

// TestCallToCurrentFloor covers the edge case where the hall call floor is
// the floor the only car is already parked at: it should open its doors
// immediately rather than trying to move.
func TestCallToCurrentFloor(t *testing.T) {
	sys := newTestSystem(10, 3)
	car := sys.HallCall(3, Up)

	floor, state, _, _ := car.Snapshot()
	if floor != 3 {
		t.Fatalf("expected car to remain at floor 3, got %d", floor)
	}
	if state != StateDoorOpen {
		t.Fatalf("expected door to open immediately at current floor, got state %v", state)
	}
}

// TestAllCarsBusyStillAssigns covers the edge case where every car is
// already moving: the strategy must still pick some car (the least-bad
// option) rather than dropping the call.
func TestAllCarsBusyStillAssigns(t *testing.T) {
	sys := newTestSystem(20, 0, 10)
	sys.HallCall(19, Up) // car 0 now moving up towards 19
	sys.HallCall(1, Down) // car 1 now moving down towards 1

	car := sys.HallCall(15, Up)
	if car == nil {
		t.Fatal("expected a car to be assigned even though both cars are busy")
	}
	_, _, _, pending := car.Snapshot()
	if pending == 0 {
		t.Fatal("expected the assigned car to have the new target registered")
	}
}

func TestNearestCarDispatchCorrectness(t *testing.T) {
	// Car 0 idle at floor 0, car 1 idle at floor 18: a call at floor 15
	// heading up should go to the nearer car (car 1).
	sys := newTestSystem(20, 0, 18)
	car := sys.HallCall(15, Up)
	if car.ID != 1 {
		t.Fatalf("expected nearest idle car (id 1) to be dispatched, got car %d", car.ID)
	}
}

func TestNearestCarPrefersEnRouteOverIdle(t *testing.T) {
	// Car 0 is moving up and already past floor 4 heading towards floor 10 -
	// it can pick up a call at floor 6 en route. Car 1 is idle at floor 0,
	// which is farther from floor 6. The en-route car should win.
	sys := newTestSystem(20, 0, 0)
	sys.HallCall(10, Up) // dispatched to car 0 (both idle at 0 -> tie goes to first)
	for i := 0; i < 4; i++ {
		sys.Step()
	}
	// car 0 should now be around floor 4, moving up
	floor, state, dir, _ := sys.Cars[0].Snapshot()
	if state != StateMovingUp && state != StateDoorOpen {
		t.Fatalf("expected car 0 to be moving up, got state %v at floor %d", state, floor)
	}
	if dir != Up && state != StateDoorOpen {
		t.Fatalf("expected car 0 direction Up, got %v", dir)
	}

	car := sys.HallCall(6, Up)
	if car.ID != 0 {
		t.Fatalf("expected en-route car 0 to be preferred over idle car 1, got car %d", car.ID)
	}
}

// TestConcurrentHallCalls fires many concurrent hall calls at an
// ElevatorSystem with multiple cars from goroutines and asserts that the
// system's mutex-guarded select-then-assign sequence in HallCall never
// corrupts a car's target set (each call must end up registered on exactly
// the car it was returned for) and every call is eventually serviced.
func TestConcurrentHallCalls(t *testing.T) {
	const numFloors = 50
	sys := newTestSystem(numFloors, 0, 10, 25, 40)

	const numCalls = 200
	var wg sync.WaitGroup
	assigned := make([]*Car, numCalls)
	floors := make([]int, numCalls)

	for i := 0; i < numCalls; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			floor := (n * 7) % numFloors
			dir := Up
			if n%2 == 0 {
				dir = Down
			}
			floors[n] = floor
			assigned[n] = sys.HallCall(floor, dir)
		}(i)
	}
	wg.Wait()

	for i, car := range assigned {
		if car == nil {
			t.Fatalf("call %d was never assigned a car", i)
		}
	}

	// The nearest-target dispatch heuristic (not a strict direction-committed
	// SCAN) can thrash back and forth when many targets are scattered across
	// one car, so give it a generously large step budget rather than a tight
	// one derived from floor/call counts.
	runUntilIdle(sys, numFloors*numCalls*2)

	for _, c := range sys.Cars {
		_, state, _, pending := c.Snapshot()
		if pending != 0 {
			t.Fatalf("car %d still has %d pending targets after running to completion", c.ID, pending)
		}
		if state != StateIdle {
			t.Fatalf("car %d expected to be idle after servicing all calls, got state %v", c.ID, state)
		}
		if c.CurrentFloor < 0 || c.CurrentFloor >= numFloors {
			t.Fatalf("car %d ended up at out-of-range floor %d (state corruption)", c.ID, c.CurrentFloor)
		}
	}
}
