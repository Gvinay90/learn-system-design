package concurrency

import (
	"fmt"
	"testing"
)

// TestNoLostOrDuplicatedItems drives multiple producers and multiple worker
// consumers through a bounded channel and asserts every produced item is
// consumed exactly once: no loss, no duplication, regardless of scheduling.
func TestNoLostOrDuplicatedItems(t *testing.T) {
	const numProducers = 5
	const itemsPerProducer = 200
	const numWorkers = 8
	const bufferSize = 16

	results := NewSafeSet()
	pool := &Pool{NumWorkers: numWorkers, BufferSize: bufferSize}
	pool.Run(numProducers, itemsPerProducer, func(item Item) {
		results.Add(item.ID())
	})

	snapshot := results.Snapshot()
	wantTotal := numProducers * itemsPerProducer
	if len(snapshot) != wantTotal {
		t.Fatalf("expected %d unique items consumed, got %d", wantTotal, len(snapshot))
	}

	for pid := 0; pid < numProducers; pid++ {
		for seq := 0; seq < itemsPerProducer; seq++ {
			id := fmt.Sprintf("p%d-%d", pid, seq)
			count, ok := snapshot[id]
			if !ok {
				t.Fatalf("item %s was never consumed", id)
			}
			if count != 1 {
				t.Fatalf("item %s consumed %d times, want exactly 1", id, count)
			}
		}
	}
}

// TestSingleProducerSingleWorker sanity-checks the pipeline degrades
// correctly to the simplest 1:1 case.
func TestSingleProducerSingleWorker(t *testing.T) {
	results := NewSafeSet()
	pool := &Pool{NumWorkers: 1, BufferSize: 1}
	pool.Run(1, 50, func(item Item) {
		results.Add(item.ID())
	})

	if got := len(results.Snapshot()); got != 50 {
		t.Fatalf("expected 50 items consumed, got %d", got)
	}
}
