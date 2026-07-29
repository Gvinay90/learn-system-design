// Package ratelimiter demonstrates two distributed-rate-limiter algorithms
// (Token Bucket, Sliding Window Counter) as single-process building blocks.
// A real distributed limiter runs this logic atomically in Redis (e.g. via
// a Lua script) so all API servers share one view of the counters — see
// the README for the full distributed design.
package ratelimiter

import (
	"sync"
	"time"
)

// TokenBucket allows a burst up to Capacity, refilling at RefillRate tokens/sec.
type TokenBucket struct {
	Capacity   float64
	RefillRate float64 // tokens per second

	mu       sync.Mutex
	tokens   float64
	lastFill time.Time
}

func NewTokenBucket(capacity, refillRate float64) *TokenBucket {
	return &TokenBucket{
		Capacity:   capacity,
		RefillRate: refillRate,
		tokens:     capacity,
		lastFill:   time.Now(),
	}
}

// Allow consumes one token if available, refilling based on elapsed time first.
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastFill).Seconds()
	b.tokens += elapsed * b.RefillRate
	if b.tokens > b.Capacity {
		b.tokens = b.Capacity
	}
	b.lastFill = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// SlidingWindowCounter approximates a sliding window using the current and
// previous fixed windows, weighted by how far into the current window we are.
// This gives O(1) memory (two counters) instead of storing every timestamp.
type SlidingWindowCounter struct {
	Limit      int
	WindowSize time.Duration

	mu            sync.Mutex
	currentWindow int64
	currentCount  int
	previousCount int
}

func NewSlidingWindowCounter(limit int, windowSize time.Duration) *SlidingWindowCounter {
	return &SlidingWindowCounter{Limit: limit, WindowSize: windowSize}
}

func (s *SlidingWindowCounter) windowID(t time.Time) int64 {
	return t.UnixNano() / int64(s.WindowSize)
}

// Allow reports whether the request at time `now` is within the limit, using
// a weighted estimate: count = currentCount + previousCount*(1 - elapsedFractionOfCurrentWindow).
func (s *SlidingWindowCounter) AllowAt(now time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	windowID := s.windowID(now)

	if windowID != s.currentWindow {
		if windowID == s.currentWindow+1 {
			s.previousCount = s.currentCount
		} else {
			s.previousCount = 0
		}
		s.currentCount = 0
		s.currentWindow = windowID
	}

	windowStart := time.Unix(0, windowID*int64(s.WindowSize))
	elapsedFraction := float64(now.Sub(windowStart)) / float64(s.WindowSize)
	if elapsedFraction > 1 {
		elapsedFraction = 1
	}
	estimate := float64(s.currentCount) + float64(s.previousCount)*(1-elapsedFraction)

	if estimate >= float64(s.Limit) {
		return false
	}
	s.currentCount++
	return true
}

func (s *SlidingWindowCounter) Allow() bool {
	return s.AllowAt(time.Now())
}
