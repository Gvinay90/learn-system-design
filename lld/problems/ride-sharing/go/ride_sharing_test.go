package ridesharing

import (
	"sync"
	"testing"
)

func newTestSystem() *RideSharingSystem {
	drivers := []*Driver{
		{ID: "D1", Name: "Alice", Location: Location{X: 0, Y: 0}, available: true},
		{ID: "D2", Name: "Bob", Location: Location{X: 10, Y: 10}, available: true},
	}
	return NewRideSharingSystem(drivers, NearestAvailableDriverStrategy{}, DistanceBasedPricing{BaseFare: 2, PerDistance: 1.5})
}

func TestHappyPathLifecycle(t *testing.T) {
	sys := newTestSystem()
	rider := &Rider{ID: "R1", Name: "Riya"}

	trip := sys.RequestRide(rider, Location{X: 0, Y: 0}, Location{X: 3, Y: 4})
	if trip.Status != Requested {
		t.Fatalf("expected Requested, got %v", trip.Status)
	}

	driver, err := sys.MatchDriver(trip.ID)
	if err != nil {
		t.Fatalf("expected match, got err: %v", err)
	}
	if driver.ID != "D1" {
		t.Fatalf("expected nearest driver D1, got %s", driver.ID)
	}
	if trip.Status != Accepted {
		t.Fatalf("expected Accepted, got %v", trip.Status)
	}

	if err := sys.StartTrip(trip.ID); err != nil {
		t.Fatalf("unexpected err starting trip: %v", err)
	}
	if trip.Status != InProgress {
		t.Fatalf("expected InProgress, got %v", trip.Status)
	}

	fare, err := sys.CompleteTrip(trip.ID)
	if err != nil {
		t.Fatalf("unexpected err completing trip: %v", err)
	}
	// base 2 + 5 (distance from (0,0) to (3,4)) * 1.5 = 9.5
	if fare != 9.5 {
		t.Fatalf("expected fare 9.5, got %v", fare)
	}
	if trip.Status != Completed {
		t.Fatalf("expected Completed, got %v", trip.Status)
	}
	if !driver.available {
		t.Fatalf("expected driver freed after trip completion")
	}
}

func TestInvalidTransitionRejected(t *testing.T) {
	sys := newTestSystem()
	rider := &Rider{ID: "R1", Name: "Riya"}
	trip := sys.RequestRide(rider, Location{X: 0, Y: 0}, Location{X: 1, Y: 1})

	if _, err := sys.CompleteTrip(trip.ID); err != ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition completing a Requested trip, got %v", err)
	}

	if _, err := sys.MatchDriver(trip.ID); err != nil {
		t.Fatalf("unexpected err matching driver: %v", err)
	}
	if err := sys.StartTrip(trip.ID); err != nil {
		t.Fatalf("unexpected err starting trip: %v", err)
	}
	if _, err := sys.CompleteTrip(trip.ID); err != nil {
		t.Fatalf("unexpected err completing trip: %v", err)
	}

	if err := sys.CancelTrip(trip.ID); err != ErrInvalidTransition {
		t.Fatalf("expected ErrInvalidTransition cancelling a Completed trip, got %v", err)
	}
}

func TestNoDriverAvailable(t *testing.T) {
	sys := newTestSystem()
	rider := &Rider{ID: "R1", Name: "Riya"}

	t1 := sys.RequestRide(rider, Location{X: 0, Y: 0}, Location{X: 1, Y: 1})
	t2 := sys.RequestRide(rider, Location{X: 0, Y: 0}, Location{X: 1, Y: 1})
	t3 := sys.RequestRide(rider, Location{X: 0, Y: 0}, Location{X: 1, Y: 1})

	if _, err := sys.MatchDriver(t1.ID); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if _, err := sys.MatchDriver(t2.ID); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if _, err := sys.MatchDriver(t3.ID); err != ErrNoDriverAvailable {
		t.Fatalf("expected ErrNoDriverAvailable, got %v", err)
	}
}

// TestConcurrentMatching asserts two goroutines racing for the same single
// available driver never both succeed — the mutex in MatchDriver must serialize them.
func TestConcurrentMatching(t *testing.T) {
	drivers := []*Driver{
		{ID: "D1", Name: "Alice", Location: Location{X: 0, Y: 0}, available: true},
	}
	sys := NewRideSharingSystem(drivers, NearestAvailableDriverStrategy{}, DistanceBasedPricing{BaseFare: 2, PerDistance: 1.5})
	rider := &Rider{ID: "R1", Name: "Riya"}

	trip1 := sys.RequestRide(rider, Location{X: 0, Y: 0}, Location{X: 1, Y: 1})
	trip2 := sys.RequestRide(rider, Location{X: 0, Y: 0}, Location{X: 1, Y: 1})

	var wg sync.WaitGroup
	results := make(chan error, 2)

	for _, tripID := range []string{trip1.ID, trip2.ID} {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			_, err := sys.MatchDriver(id)
			results <- err
		}(tripID)
	}
	wg.Wait()
	close(results)

	successCount := 0
	for err := range results {
		if err == nil {
			successCount++
		}
	}
	if successCount != 1 {
		t.Fatalf("expected exactly 1 success for single available driver, got %d", successCount)
	}
}
