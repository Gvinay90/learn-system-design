// Package paymentgateway implements a single-process payment gateway LLD:
// idempotency-key deduplication, a pluggable processor with retries, and a
// simple in-memory ledger.
package paymentgateway

import (
	"fmt"
	"sync"
	"time"
)

type Status int

const (
	Pending Status = iota
	Success
	Failed
)

func (s Status) String() string {
	switch s {
	case Success:
		return "SUCCESS"
	case Failed:
		return "FAILED"
	default:
		return "PENDING"
	}
}

type PaymentRequest struct {
	IdempotencyKey string
	PayerID        string
	PayeeID        string
	Amount         float64
	Currency       string
}

type Attempt struct {
	Number  int
	Success bool
	Err     string
	At      time.Time
}

type PaymentResult struct {
	ID        string
	Request   PaymentRequest
	Status    Status
	CreatedAt time.Time
	Attempts  []Attempt
}

type LedgerEntry struct {
	PaymentID string
	PayerID   string
	PayeeID   string
	Amount    float64
	At        time.Time
}

// Ledger is a simplified single-entry-per-payment record. A real
// double-entry ledger would post two rows (debit payer / credit payee);
// one row is kept here to stay focused on the idempotency/retry logic.
type Ledger struct {
	mu      sync.Mutex
	entries []LedgerEntry
}

func NewLedger() *Ledger {
	return &Ledger{}
}

func (l *Ledger) record(e LedgerEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, e)
}

func (l *Ledger) Entries() []LedgerEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]LedgerEntry, len(l.entries))
	copy(out, l.entries)
	return out
}

// PaymentProcessor is the Strategy interface for actually moving money.
// Real implementations would call a card network / bank rail; tests use a
// fake that can simulate transient failures.
type PaymentProcessor interface {
	Process(req PaymentRequest) error
}

type RetryPolicy struct {
	MaxAttempts int
	Delay       time.Duration
}

func (r RetryPolicy) delayFor(attempt int) time.Duration {
	return time.Duration(attempt) * r.Delay
}

// idempotencyEntry tracks an in-flight or completed charge for a key.
// done is closed once result is populated, letting concurrent callers with
// the same key block until the single owner finishes instead of racing.
type idempotencyEntry struct {
	done   chan struct{}
	result *PaymentResult
}

type IdempotencyStore struct {
	mu      sync.Mutex
	entries map[string]*idempotencyEntry
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{entries: make(map[string]*idempotencyEntry)}
}

// reserveOrWait returns (storedResult, true) if this call is the first to
// see key and must perform the charge, or (result, false) if another
// caller already owns/owned it -- in which case it blocks until that
// caller finishes and returns the same result.
func (s *IdempotencyStore) reserveOrWait(key string) (*PaymentResult, bool) {
	s.mu.Lock()
	if e, ok := s.entries[key]; ok {
		s.mu.Unlock()
		<-e.done
		return e.result, false
	}
	e := &idempotencyEntry{done: make(chan struct{})}
	s.entries[key] = e
	s.mu.Unlock()
	return nil, true
}

func (s *IdempotencyStore) complete(key string, result *PaymentResult) {
	s.mu.Lock()
	e := s.entries[key]
	s.mu.Unlock()
	e.result = result
	close(e.done)
}

type PaymentGateway struct {
	processor PaymentProcessor
	retry     RetryPolicy
	store     *IdempotencyStore
	ledger    *Ledger

	mu  sync.Mutex
	seq int
}

func NewPaymentGateway(processor PaymentProcessor, retry RetryPolicy) *PaymentGateway {
	return &PaymentGateway{
		processor: processor,
		retry:     retry,
		store:     NewIdempotencyStore(),
		ledger:    NewLedger(),
	}
}

func (g *PaymentGateway) Ledger() *Ledger {
	return g.ledger
}

// Charge processes req, or returns the previously stored terminal result
// if req.IdempotencyKey has already been used -- same key always yields
// the same result, even if that result was FAILED (some real gateways
// allow retrying a failed idempotency key; this exercise keeps the
// simpler "terminal result is final" semantics).
func (g *PaymentGateway) Charge(req PaymentRequest) *PaymentResult {
	result, isOwner := g.store.reserveOrWait(req.IdempotencyKey)
	if !isOwner {
		return result
	}

	result = g.processWithRetry(req)
	g.store.complete(req.IdempotencyKey, result)
	return result
}

func (g *PaymentGateway) processWithRetry(req PaymentRequest) *PaymentResult {
	g.mu.Lock()
	g.seq++
	id := fmt.Sprintf("PAY-%d", g.seq)
	g.mu.Unlock()

	result := &PaymentResult{ID: id, Request: req, Status: Pending, CreatedAt: time.Now()}

	for attempt := 1; attempt <= g.retry.MaxAttempts; attempt++ {
		err := g.processor.Process(req)
		a := Attempt{Number: attempt, Success: err == nil, At: time.Now()}
		if err != nil {
			a.Err = err.Error()
		}
		result.Attempts = append(result.Attempts, a)

		if err == nil {
			result.Status = Success
			g.ledger.record(LedgerEntry{
				PaymentID: id,
				PayerID:   req.PayerID,
				PayeeID:   req.PayeeID,
				Amount:    req.Amount,
				At:        time.Now(),
			})
			return result
		}

		if attempt < g.retry.MaxAttempts {
			time.Sleep(g.retry.delayFor(attempt))
		}
	}

	result.Status = Failed
	return result
}
