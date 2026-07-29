// Package concurrency demonstrates the core concurrency primitives used in
// LLD interviews: goroutines, a bounded channel as the shared buffer, a
// worker-pool of consumers, and sync.WaitGroup for lifecycle coordination.
package concurrency

import (
	"fmt"
	"sync"
)

// Item is a unit of work flowing from producers to workers.
type Item struct {
	ProducerID int
	Seq        int
}

// ID uniquely identifies an item across all producers so consumers can
// dedupe/verify delivery.
func (i Item) ID() string {
	return fmt.Sprintf("p%d-%d", i.ProducerID, i.Seq)
}

// Pool runs a bounded producer-consumer pipeline: numProducers goroutines
// each emit itemsPerProducer items into a buffered channel of size
// bufferSize (the bounded buffer), and numWorkers goroutines drain it
// concurrently, invoking handle for every item they receive.
type Pool struct {
	NumWorkers int
	BufferSize int
}

// Run starts producers and workers, blocks until every item has been
// produced and consumed, then returns. handle is called once per item from
// whichever worker goroutine received it, so callers must make it
// concurrency-safe themselves (e.g. guard shared state with a mutex).
func (p *Pool) Run(numProducers, itemsPerProducer int, handle func(Item)) {
	buffer := make(chan Item, p.BufferSize)

	var producers sync.WaitGroup
	producers.Add(numProducers)
	for pid := 0; pid < numProducers; pid++ {
		go func(pid int) {
			defer producers.Done()
			for seq := 0; seq < itemsPerProducer; seq++ {
				buffer <- Item{ProducerID: pid, Seq: seq}
			}
		}(pid)
	}

	var workers sync.WaitGroup
	workers.Add(p.NumWorkers)
	for w := 0; w < p.NumWorkers; w++ {
		go func() {
			defer workers.Done()
			for item := range buffer {
				handle(item)
			}
		}()
	}

	// Closing the channel only after all producers finish signals workers
	// to drain the remaining buffered items and exit the range loop.
	producers.Wait()
	close(buffer)
	workers.Wait()
}

// SafeSet is a mutex-guarded string set used to collect consumed item IDs
// concurrently without losing updates (the map itself is not goroutine-safe).
type SafeSet struct {
	mu   sync.Mutex
	seen map[string]int
}

func NewSafeSet() *SafeSet {
	return &SafeSet{seen: make(map[string]int)}
}

func (s *SafeSet) Add(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen[id]++
}

func (s *SafeSet) Snapshot() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]int, len(s.seen))
	for k, v := range s.seen {
		out[k] = v
	}
	return out
}
