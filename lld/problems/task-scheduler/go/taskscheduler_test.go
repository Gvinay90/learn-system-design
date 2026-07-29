package taskscheduler

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestJobsRunAndSucceed(t *testing.T) {
	s := NewScheduler(2)
	s.Start()
	defer s.Stop()

	var ran int32
	s.Submit(&Job{
		ID:       "J1",
		Priority: 1,
		RunAt:    time.Now(),
		Task: func() error {
			atomic.AddInt32(&ran, 1)
			return nil
		},
	})

	res, ok := s.WaitForResult("J1", time.Second)
	if !ok {
		t.Fatalf("expected job to complete within timeout")
	}
	if res.Status != Succeeded {
		t.Fatalf("expected Succeeded, got %v", res.Status)
	}
	if res.Attempts != 1 {
		t.Fatalf("expected 1 attempt, got %d", res.Attempts)
	}
	if atomic.LoadInt32(&ran) != 1 {
		t.Fatalf("expected task to run exactly once")
	}
}

// TestPriorityOrdering asserts that among jobs due at the same time, the
// higher-priority job runs first. We use a single worker so execution order
// is deterministic and use RunAt in the past for both so they're both "due"
// as soon as the scheduler starts.
func TestPriorityOrdering(t *testing.T) {
	s := NewScheduler(1)

	var mu sync.Mutex
	var order []string
	due := time.Now().Add(20 * time.Millisecond)

	s.Submit(&Job{
		ID:       "low",
		Priority: 1,
		RunAt:    due,
		Task: func() error {
			mu.Lock()
			order = append(order, "low")
			mu.Unlock()
			return nil
		},
	})
	s.Submit(&Job{
		ID:       "high",
		Priority: 10,
		RunAt:    due,
		Task: func() error {
			mu.Lock()
			order = append(order, "high")
			mu.Unlock()
			return nil
		},
	})

	s.Start()
	defer s.Stop()

	if _, ok := s.WaitForResult("low", time.Second); !ok {
		t.Fatalf("expected low priority job to complete")
	}
	if _, ok := s.WaitForResult("high", time.Second); !ok {
		t.Fatalf("expected high priority job to complete")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(order) != 2 || order[0] != "high" || order[1] != "low" {
		t.Fatalf("expected high-priority job to run before low-priority job, got %v", order)
	}
}

func TestDelayedJobDoesNotRunEarly(t *testing.T) {
	s := NewScheduler(1)
	s.Start()
	defer s.Stop()

	var ranAt time.Time
	submittedAt := time.Now()
	delay := 60 * time.Millisecond

	s.Submit(&Job{
		ID:       "delayed",
		Priority: 1,
		RunAt:    submittedAt.Add(delay),
		Task: func() error {
			ranAt = time.Now()
			return nil
		},
	})

	res, ok := s.WaitForResult("delayed", time.Second)
	if !ok || res.Status != Succeeded {
		t.Fatalf("expected delayed job to eventually succeed, got %+v ok=%v", res, ok)
	}
	if ranAt.Sub(submittedAt) < delay {
		t.Fatalf("expected job to run no earlier than %v after submit, ran after %v", delay, ranAt.Sub(submittedAt))
	}
}

func TestRetriesOnFailureThenSucceeds(t *testing.T) {
	s := NewScheduler(1)
	s.Start()
	defer s.Stop()

	var attempts int32
	s.Submit(&Job{
		ID:         "flaky",
		Priority:   1,
		RunAt:      time.Now(),
		MaxRetries: 3,
		Retry:      ExponentialBackoff{Base: 2 * time.Millisecond, Max: 10 * time.Millisecond},
		Task: func() error {
			n := atomic.AddInt32(&attempts, 1)
			if n < 3 {
				return errors.New("transient failure")
			}
			return nil
		},
	})

	res, ok := s.WaitForResult("flaky", time.Second)
	if !ok {
		t.Fatalf("expected flaky job to eventually reach a terminal state")
	}
	if res.Status != Succeeded {
		t.Fatalf("expected Succeeded after retries, got %v (err=%v)", res.Status, res.Err)
	}
	if res.Attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", res.Attempts)
	}
}

func TestExhaustsRetriesAndFails(t *testing.T) {
	s := NewScheduler(1)
	s.Start()
	defer s.Stop()

	var attempts int32
	permanentErr := errors.New("permanent failure")
	s.Submit(&Job{
		ID:         "doomed",
		Priority:   1,
		RunAt:      time.Now(),
		MaxRetries: 2,
		Retry:      ExponentialBackoff{Base: 2 * time.Millisecond, Max: 10 * time.Millisecond},
		Task: func() error {
			atomic.AddInt32(&attempts, 1)
			return permanentErr
		},
	})

	res, ok := s.WaitForResult("doomed", time.Second)
	if !ok {
		t.Fatalf("expected doomed job to eventually reach a terminal state")
	}
	if res.Status != Failed {
		t.Fatalf("expected Failed, got %v", res.Status)
	}
	// MaxRetries=2 means 1 initial attempt + 2 retries = 3 attempts total.
	if res.Attempts != 3 {
		t.Fatalf("expected 3 total attempts (1 initial + 2 retries), got %d", res.Attempts)
	}
	if res.Err == nil || res.Err.Error() != permanentErr.Error() {
		t.Fatalf("expected the last error preserved, got %v", res.Err)
	}
}

func TestUnknownJobResultNotFound(t *testing.T) {
	s := NewScheduler(1)
	if _, ok := s.GetResult("nonexistent"); ok {
		t.Fatalf("expected no result for unknown job id")
	}
}

func TestExponentialBackoffCapsAtMax(t *testing.T) {
	b := ExponentialBackoff{Base: 10 * time.Millisecond, Max: 30 * time.Millisecond}
	if d := b.NextDelay(0); d != 10*time.Millisecond {
		t.Fatalf("expected 10ms for attempt 0, got %v", d)
	}
	if d := b.NextDelay(1); d != 20*time.Millisecond {
		t.Fatalf("expected 20ms for attempt 1, got %v", d)
	}
	if d := b.NextDelay(2); d != 30*time.Millisecond {
		t.Fatalf("expected backoff capped at 30ms for attempt 2, got %v", d)
	}
	if d := b.NextDelay(5); d != 30*time.Millisecond {
		t.Fatalf("expected backoff to stay capped at 30ms for larger attempts, got %v", d)
	}
}

// TestConcurrentSubmitAndExecute asserts many jobs submitted concurrently
// from multiple goroutines all get executed exactly once by the worker pool.
func TestConcurrentSubmitAndExecute(t *testing.T) {
	s := NewScheduler(4)
	s.Start()
	defer s.Stop()

	const n = 50
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := jobID(i)
			s.Submit(&Job{
				ID:       id,
				Priority: i % 5,
				RunAt:    time.Now(),
				Task:     func() error { return nil },
			})
		}(i)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		id := jobID(i)
		res, ok := s.WaitForResult(id, 2*time.Second)
		if !ok || res.Status != Succeeded {
			t.Fatalf("expected job %s to succeed, got %+v ok=%v", id, res, ok)
		}
	}
}

func jobID(i int) string {
	return "concurrent-" + string(rune('A'+i%26)) + string(rune('0'+i/26))
}
