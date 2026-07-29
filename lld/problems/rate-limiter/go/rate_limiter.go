// Package ratelimiter implements a single-process, per-client rate limiter
// with two interchangeable algorithms — token bucket and sliding window —
// behind a common RateLimiter interface. See ../README.md for the design
// writeup this package implements.
package ratelimiter

import (
	"sync"
	"time"
)

// RateLimiter is the common strategy interface every algorithm implements,
// so calling code never branches on which algorithm is active.
type RateLimiter interface {
	// Allow decides whether a request from clientID is permitted right now.
	Allow(clientID string) bool
}

// clock returns the current time; it is a field (not a package var) on each
// limiter so tests can inject a fake clock instead of sleeping real time.
type clock func() time.Time

// bucket holds one client's token-bucket state.
type bucket struct {
	tokens     float64
	lastRefill time.Time
}

// TokenBucketLimiter allows bursts up to capacity tokens, refilling at
// refillRate tokens/second. Refill is computed lazily on each Allow call
// rather than via a background goroutine per client.
type TokenBucketLimiter struct {
	capacity   float64
	refillRate float64 // tokens per second
	now        clock

	mu      sync.Mutex
	buckets map[string]*bucket
}

// NewTokenBucketLimiter builds a token bucket limiter with the given
// capacity (max burst size) and refillRate (tokens added per second).
func NewTokenBucketLimiter(capacity, refillRate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:   capacity,
		refillRate: refillRate,
		now:        time.Now,
		buckets:    make(map[string]*bucket),
	}
}

// refill tops b up based on elapsed time since its last refill, capping at
// capacity. Caller must hold l.mu.
func (l *TokenBucketLimiter) refill(b *bucket, now time.Time) {
	elapsed := now.Sub(b.lastRefill).Seconds()
	if elapsed <= 0 {
		return
	}
	b.tokens += elapsed * l.refillRate
	if b.tokens > l.capacity {
		b.tokens = l.capacity
	}
	b.lastRefill = now
}

// Allow consumes one token for clientID if available.
func (l *TokenBucketLimiter) Allow(clientID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.buckets[clientID]
	if !ok {
		b = &bucket{tokens: l.capacity, lastRefill: now}
		l.buckets[clientID] = b
	} else {
		l.refill(b, now)
	}

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// window holds one client's sliding-log state: timestamps of requests still
// inside the trailing window.
type window struct {
	timestamps []time.Time
}

// SlidingWindowLimiter allows at most limit requests in any trailing
// duration of length window (a sliding log), giving a smooth, hard cap with
// no boundary burst allowance.
type SlidingWindowLimiter struct {
	limit  int
	window time.Duration
	now    clock

	mu      sync.Mutex
	windows map[string]*window
}

// NewSlidingWindowLimiter builds a sliding window limiter allowing at most
// limit requests within any trailing duration of the given window size.
func NewSlidingWindowLimiter(limit int, windowSize time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:   limit,
		window:  windowSize,
		now:     time.Now,
		windows: make(map[string]*window),
	}
}

// Allow evicts timestamps older than the window and checks the remaining
// count against limit, recording the new request if it fits.
func (l *SlidingWindowLimiter) Allow(clientID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	w, ok := l.windows[clientID]
	if !ok {
		w = &window{}
		l.windows[clientID] = w
	}

	cutoff := now.Add(-l.window)
	kept := w.timestamps[:0]
	for _, ts := range w.timestamps {
		if ts.After(cutoff) {
			kept = append(kept, ts)
		}
	}
	w.timestamps = kept

	if len(w.timestamps) >= l.limit {
		return false
	}
	w.timestamps = append(w.timestamps, now)
	return true
}

// RateLimiterRegistry maps a client class (e.g. free/paid/enterprise) to the
// RateLimiter instance configured for that tier, so different tiers can even
// run different algorithms.
type RateLimiterRegistry struct {
	mu             sync.RWMutex
	perClientClass map[string]RateLimiter
}

// NewRateLimiterRegistry builds an empty registry.
func NewRateLimiterRegistry() *RateLimiterRegistry {
	return &RateLimiterRegistry{perClientClass: make(map[string]RateLimiter)}
}

// Register associates clientClass with limiter, overwriting any prior
// registration for that class.
func (r *RateLimiterRegistry) Register(clientClass string, limiter RateLimiter) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.perClientClass[clientClass] = limiter
}

// GetLimiter returns the RateLimiter configured for clientClass, and false
// if no limiter has been registered for it.
func (r *RateLimiterRegistry) GetLimiter(clientClass string) (RateLimiter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	limiter, ok := r.perClientClass[clientClass]
	return limiter, ok
}
