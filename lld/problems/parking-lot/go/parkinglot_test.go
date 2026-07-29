package parkinglot

import (
	"sync"
	"testing"
	"time"
)

func newTestLot() *ParkingLot {
	floor1 := &Floor{Level: 1, Spots: []*ParkingSpot{
		{ID: "F1-M1", Type: Motorcycle},
		{ID: "F1-C1", Type: Compact},
		{ID: "F1-L1", Type: Large},
	}}
	return NewParkingLot([]*Floor{floor1}, HourlyTieredPricing{FirstHourRate: 5, SubsequentRate: 2})
}

func TestParkAndUnpark(t *testing.T) {
	lot := newTestLot()
	v := &Vehicle{LicensePlate: "KA-01", Type: Compact}

	ticket, err := lot.ParkVehicle(v)
	if err != nil {
		t.Fatalf("expected to park, got err: %v", err)
	}
	if ticket.Spot.ID != "F1-C1" {
		t.Fatalf("expected assigned to F1-C1, got %s", ticket.Spot.ID)
	}

	fee, err := lot.UnparkVehicle(ticket.ID, ticket.EntryTime.Add(30*time.Minute))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if fee != 5 {
		t.Fatalf("expected flat first-hour fee of 5, got %v", fee)
	}
}

func TestLotFullForType(t *testing.T) {
	lot := newTestLot()
	_, _ = lot.ParkVehicle(&Vehicle{LicensePlate: "A1", Type: Compact})

	_, err := lot.ParkVehicle(&Vehicle{LicensePlate: "A2", Type: Compact})
	if err != ErrLotFull {
		t.Fatalf("expected ErrLotFull, got %v", err)
	}
}

func TestUnparkUnknownTicket(t *testing.T) {
	lot := newTestLot()
	_, err := lot.UnparkVehicle("bogus", time.Now())
	if err != ErrTicketNotFound {
		t.Fatalf("expected ErrTicketNotFound, got %v", err)
	}
}

func TestHourlyTieredPricing(t *testing.T) {
	lot := newTestLot()
	v := &Vehicle{LicensePlate: "KA-02", Type: Large}
	ticket, _ := lot.ParkVehicle(v)

	fee, _ := lot.UnparkVehicle(ticket.ID, ticket.EntryTime.Add(3*time.Hour))
	// first hour 5 + 2 subsequent hours * 2 = 9
	if fee != 9 {
		t.Fatalf("expected fee 9, got %v", fee)
	}
}

// TestConcurrentParking asserts two goroutines racing for the same single spot
// never both succeed — the mutex in Floor.findAndAssign must serialize them.
func TestConcurrentParking(t *testing.T) {
	lot := newTestLot()
	var wg sync.WaitGroup
	successes := make(chan bool, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_, err := lot.ParkVehicle(&Vehicle{LicensePlate: "race", Type: Motorcycle})
			successes <- err == nil
		}(i)
	}
	wg.Wait()
	close(successes)

	successCount := 0
	for s := range successes {
		if s {
			successCount++
		}
	}
	if successCount != 1 {
		t.Fatalf("expected exactly 1 success for single motorcycle spot, got %d", successCount)
	}
}
