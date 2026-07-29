package ratelimiter

import (
	"testing"
	"time"
)

func TestTokenBucketAllowsUpToCapacityThenBlocks(t *testing.T) {
	b := NewTokenBucket(3, 1) // capacity 3, refill 1/sec

	for i := 0; i < 3; i++ {
		if !b.Allow() {
			t.Fatalf("expected request %d to be allowed (within initial burst capacity)", i)
		}
	}
	if b.Allow() {
		t.Fatalf("expected 4th immediate request to be blocked, bucket should be empty")
	}
}

func TestTokenBucketRefillsOverTime(t *testing.T) {
	b := NewTokenBucket(1, 10) // capacity 1, refill 10/sec -> refills in 100ms
	if !b.Allow() {
		t.Fatalf("expected first request to be allowed")
	}
	if b.Allow() {
		t.Fatalf("expected immediate second request to be blocked")
	}
	time.Sleep(150 * time.Millisecond)
	if !b.Allow() {
		t.Fatalf("expected request to be allowed after refill window")
	}
}

func TestSlidingWindowCounterBlocksOverLimit(t *testing.T) {
	s := NewSlidingWindowCounter(2, time.Second)
	base := time.Unix(0, 0)

	if !s.AllowAt(base) {
		t.Fatalf("expected 1st request allowed")
	}
	if !s.AllowAt(base.Add(100 * time.Millisecond)) {
		t.Fatalf("expected 2nd request allowed")
	}
	if s.AllowAt(base.Add(200 * time.Millisecond)) {
		t.Fatalf("expected 3rd request in same window to be blocked")
	}
}

func TestSlidingWindowCounterSmoothsAcrossWindowBoundary(t *testing.T) {
	s := NewSlidingWindowCounter(2, time.Second)
	base := time.Unix(0, 0)

	// Fill the previous window fully.
	s.AllowAt(base.Add(900 * time.Millisecond))
	s.AllowAt(base.Add(950 * time.Millisecond))

	// Right at the window rollover, elapsed fraction into the new window is
	// ~0, so the full previous-window count is carried over and should still
	// count against the new window's limit.
	if s.AllowAt(base.Add(1000 * time.Millisecond)) {
		t.Fatalf("expected request right at window rollover to be blocked due to carried-over weight")
	}
}
