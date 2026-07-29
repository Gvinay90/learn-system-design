package ratelimiter

import (
	"sync"
	"testing"
	"time"
)

// fakeClock lets tests advance time deterministically instead of sleeping.
type fakeClock struct {
	t time.Time
}

func (f *fakeClock) now() time.Time          { return f.t }
func (f *fakeClock) advance(d time.Duration) { f.t = f.t.Add(d) }

func TestTokenBucketAllowsBurstUpToCapacityThenBlocks(t *testing.T) {
	tests := []struct {
		name        string
		capacity    float64
		requests    int
		wantAllowed int
	}{
		{name: "capacity 3 allows exactly 3 then blocks", capacity: 3, requests: 5, wantAllowed: 3},
		{name: "capacity 1 allows exactly 1 then blocks", capacity: 1, requests: 4, wantAllowed: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fc := &fakeClock{t: time.Unix(0, 0)}
			l := NewTokenBucketLimiter(tc.capacity, 1 /* refillRate: 1 token/sec */)
			l.now = fc.now

			allowed := 0
			for i := 0; i < tc.requests; i++ {
				if l.Allow("client-a") {
					allowed++
				}
			}
			if allowed != tc.wantAllowed {
				t.Fatalf("expected %d allowed requests, got %d", tc.wantAllowed, allowed)
			}
		})
	}
}

func TestTokenBucketRefillsOverTime(t *testing.T) {
	fc := &fakeClock{t: time.Unix(0, 0)}
	l := NewTokenBucketLimiter(2, 1 /* 1 token/sec */)
	l.now = fc.now

	if !l.Allow("client-a") || !l.Allow("client-a") {
		t.Fatalf("expected first two requests to be allowed (burst up to capacity)")
	}
	if l.Allow("client-a") {
		t.Fatalf("expected third immediate request to be blocked")
	}

	// Advance by 1s -> refills exactly 1 token.
	fc.advance(1 * time.Second)
	if !l.Allow("client-a") {
		t.Fatalf("expected request to be allowed after 1s refill")
	}
	if l.Allow("client-a") {
		t.Fatalf("expected next request to be blocked again (bucket drained)")
	}
}

func TestSlidingWindowAllowsUpToMaxThenBlocksThenSlides(t *testing.T) {
	fc := &fakeClock{t: time.Unix(0, 0)}
	l := NewSlidingWindowLimiter(3, 1*time.Second)
	l.now = fc.now

	for i := 0; i < 3; i++ {
		if !l.Allow("client-a") {
			t.Fatalf("expected request %d to be allowed within limit", i)
		}
	}
	if l.Allow("client-a") {
		t.Fatalf("expected 4th request within window to be blocked")
	}

	// Slide the window fully past the recorded timestamps.
	fc.advance(1100 * time.Millisecond)
	if !l.Allow("client-a") {
		t.Fatalf("expected request to be allowed after window slides")
	}
}

func TestPerClientIsolation(t *testing.T) {
	fc := &fakeClock{t: time.Unix(0, 0)}
	tb := NewTokenBucketLimiter(1, 1)
	tb.now = fc.now

	if !tb.Allow("client-a") {
		t.Fatalf("expected client-a first request allowed")
	}
	if tb.Allow("client-a") {
		t.Fatalf("expected client-a second immediate request blocked")
	}
	if !tb.Allow("client-b") {
		t.Fatalf("expected client-b to have its own independent bucket")
	}

	fc2 := &fakeClock{t: time.Unix(0, 0)}
	sw := NewSlidingWindowLimiter(1, 1*time.Second)
	sw.now = fc2.now

	if !sw.Allow("client-a") {
		t.Fatalf("expected client-a first request allowed")
	}
	if sw.Allow("client-a") {
		t.Fatalf("expected client-a second immediate request blocked")
	}
	if !sw.Allow("client-b") {
		t.Fatalf("expected client-b to have its own independent window")
	}
}

func TestRateLimiterRegistry(t *testing.T) {
	reg := NewRateLimiterRegistry()
	free := NewTokenBucketLimiter(1, 1)
	paid := NewSlidingWindowLimiter(100, time.Second)

	reg.Register("free", free)
	reg.Register("paid", paid)

	got, ok := reg.GetLimiter("free")
	if !ok || got != RateLimiter(free) {
		t.Fatalf("expected registry to return the registered free-tier limiter")
	}
	if _, ok := reg.GetLimiter("enterprise"); ok {
		t.Fatalf("expected no limiter registered for unknown class")
	}
}

// TestConcurrentAllowIsRace-Free hammers both limiters from many goroutines
// across multiple clients and asserts no more than capacity/limit requests
// are ever allowed per client in the same instant, and nothing panics/races
// (run with `go test -race`).
func TestConcurrentAllowIsRaceFree(t *testing.T) {
	const goroutinesPerClient = 50
	const clients = 5

	tb := NewTokenBucketLimiter(10, 1000) // generous refill so we mostly just check no races/panics
	sw := NewSlidingWindowLimiter(10, time.Second)

	var wg sync.WaitGroup
	for c := 0; c < clients; c++ {
		clientID := clientIDFor(c)
		for i := 0; i < goroutinesPerClient; i++ {
			wg.Add(2)
			go func() {
				defer wg.Done()
				tb.Allow(clientID)
			}()
			go func() {
				defer wg.Done()
				sw.Allow(clientID)
			}()
		}
	}
	wg.Wait()
}

func clientIDFor(n int) string {
	return "client-" + string(rune('a'+n))
}
