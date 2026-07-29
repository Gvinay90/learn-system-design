package paymentgateway

import (
	"errors"
	"sync"
	"testing"
	"time"
)

type fakeProcessor struct {
	mu         sync.Mutex
	callCount  int
	failTimes  int
	alwaysFail bool
}

func (f *fakeProcessor) Process(req PaymentRequest) error {
	f.mu.Lock()
	f.callCount++
	count := f.callCount
	f.mu.Unlock()

	if f.alwaysFail {
		return errors.New("simulated permanent failure")
	}
	if count <= f.failTimes {
		return errors.New("simulated transient failure")
	}
	return nil
}

func (f *fakeProcessor) calls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.callCount
}

func testRequest(key string) PaymentRequest {
	return PaymentRequest{
		IdempotencyKey: key,
		PayerID:        "payer-1",
		PayeeID:        "payee-1",
		Amount:         100,
		Currency:       "INR",
	}
}

func TestHappyPathChargeRecordsLedgerEntry(t *testing.T) {
	proc := &fakeProcessor{}
	gw := NewPaymentGateway(proc, RetryPolicy{MaxAttempts: 3, Delay: time.Millisecond})

	result := gw.Charge(testRequest("key-1"))
	if result.Status != Success {
		t.Fatalf("expected SUCCESS, got %v", result.Status)
	}

	entries := gw.Ledger().Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 ledger entry, got %d", len(entries))
	}
	if entries[0].PaymentID != result.ID || entries[0].Amount != 100 {
		t.Fatalf("ledger entry mismatch: %+v", entries[0])
	}
}

func TestSameIdempotencyKeyDoesNotReprocess(t *testing.T) {
	proc := &fakeProcessor{}
	gw := NewPaymentGateway(proc, RetryPolicy{MaxAttempts: 3, Delay: time.Millisecond})

	first := gw.Charge(testRequest("key-2"))
	second := gw.Charge(testRequest("key-2"))

	if proc.calls() != 1 {
		t.Fatalf("expected processor invoked exactly once, got %d", proc.calls())
	}
	if first.ID != second.ID || first.Status != second.Status {
		t.Fatalf("expected identical cached result, got %+v vs %+v", first, second)
	}
}

func TestRetryPolicySucceedsAfterTransientFailures(t *testing.T) {
	proc := &fakeProcessor{failTimes: 2}
	gw := NewPaymentGateway(proc, RetryPolicy{MaxAttempts: 3, Delay: time.Millisecond})

	result := gw.Charge(testRequest("key-3"))
	if result.Status != Success {
		t.Fatalf("expected SUCCESS after retries, got %v", result.Status)
	}
	if len(result.Attempts) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(result.Attempts))
	}
}

func TestRetryPolicyExhaustsToFailed(t *testing.T) {
	proc := &fakeProcessor{alwaysFail: true}
	gw := NewPaymentGateway(proc, RetryPolicy{MaxAttempts: 3, Delay: time.Millisecond})

	result := gw.Charge(testRequest("key-4"))
	if result.Status != Failed {
		t.Fatalf("expected FAILED after exhausting retries, got %v", result.Status)
	}
	if len(result.Attempts) != 3 {
		t.Fatalf("expected 3 attempts, got %d", len(result.Attempts))
	}

	// Same key retried by the caller should return the stored FAILED result,
	// not reprocess.
	again := gw.Charge(testRequest("key-4"))
	if again.Status != Failed || proc.calls() != 3 {
		t.Fatalf("expected cached FAILED result without reprocessing, calls=%d", proc.calls())
	}
}

func TestConcurrentChargesSameKeyProcessOnce(t *testing.T) {
	proc := &fakeProcessor{}
	gw := NewPaymentGateway(proc, RetryPolicy{MaxAttempts: 3, Delay: time.Millisecond})

	const n = 20
	var wg sync.WaitGroup
	results := make([]*PaymentResult, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = gw.Charge(testRequest("shared-key"))
		}(i)
	}
	wg.Wait()

	if proc.calls() != 1 {
		t.Fatalf("expected processor invoked exactly once, got %d", proc.calls())
	}
	firstID := results[0].ID
	for i, r := range results {
		if r.ID != firstID || r.Status != Success {
			t.Fatalf("goroutine %d got divergent result: %+v", i, r)
		}
	}
}
