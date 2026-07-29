// Package taskscheduler implements an in-memory, priority-based job scheduler
// with delayed execution and retries (single process — see the HLD track for
// a distributed cron/scheduler design).
package taskscheduler

import (
	"container/heap"
	"sync"
	"time"
)

type Task func() error

type Status int

const (
	Pending Status = iota
	Running
	Retrying
	Succeeded
	Failed
)

// RetryPolicy computes the delay before the (attempt+1)-th retry.
type RetryPolicy interface {
	NextDelay(attempt int) time.Duration
}

type ExponentialBackoff struct {
	Base time.Duration
	Max  time.Duration
}

func (b ExponentialBackoff) NextDelay(attempt int) time.Duration {
	d := b.Base
	for i := 0; i < attempt; i++ {
		d *= 2
		if d > b.Max {
			return b.Max
		}
	}
	return d
}

type Job struct {
	ID         string
	Priority   int
	RunAt      time.Time
	Task       Task
	MaxRetries int
	Retry      RetryPolicy

	attempts int
}

type JobResult struct {
	JobID    string
	Status   Status
	Attempts int
	Err      error
}

// jobQueue is a min-heap ordered by RunAt (earliest due first), tie-broken by
// higher Priority first.
type jobQueue []*Job

func (q jobQueue) Len() int { return len(q) }
func (q jobQueue) Less(i, j int) bool {
	if q[i].RunAt.Equal(q[j].RunAt) {
		return q[i].Priority > q[j].Priority
	}
	return q[i].RunAt.Before(q[j].RunAt)
}
func (q jobQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *jobQueue) Push(x any)   { *q = append(*q, x.(*Job)) }
func (q *jobQueue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[:n-1]
	return item
}

var defaultRetryPolicy = ExponentialBackoff{Base: 10 * time.Millisecond, Max: 200 * time.Millisecond}

// Scheduler runs submitted jobs on a fixed worker pool, respecting priority,
// delayed RunAt times, and per-job retry policies.
type Scheduler struct {
	mu      sync.Mutex
	queue   jobQueue
	results map[string]*JobResult
	workers int
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

func NewScheduler(workers int) *Scheduler {
	return &Scheduler{results: make(map[string]*JobResult), workers: workers}
}

// Submit enqueues a job. Safe to call before or after Start, and concurrently
// from multiple goroutines.
func (s *Scheduler) Submit(j *Job) {
	if j.Retry == nil {
		j.Retry = defaultRetryPolicy
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	heap.Push(&s.queue, j)
	s.results[j.ID] = &JobResult{JobID: j.ID, Status: Pending}
}

func (s *Scheduler) Start() {
	s.stopCh = make(chan struct{})
	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.workerLoop()
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

func (s *Scheduler) workerLoop() {
	defer s.wg.Done()
	ticker := time.NewTicker(2 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			if job := s.popDue(); job != nil {
				s.execute(job)
			}
		}
	}
}

// popDue removes and returns the highest-priority job whose RunAt has
// already elapsed, or nil if none is ready. Locking here (rather than around
// the whole execute) keeps job execution off the critical section so workers
// don't serialize on long-running tasks.
func (s *Scheduler) popDue() *Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.queue.Len() == 0 {
		return nil
	}
	if s.queue[0].RunAt.After(time.Now()) {
		return nil
	}
	return heap.Pop(&s.queue).(*Job)
}

func (s *Scheduler) execute(j *Job) {
	s.mu.Lock()
	s.results[j.ID].Status = Running
	s.mu.Unlock()

	err := j.Task()
	j.attempts++

	s.mu.Lock()
	defer s.mu.Unlock()
	res := s.results[j.ID]
	res.Attempts = j.attempts
	if err == nil {
		res.Status = Succeeded
		res.Err = nil
		return
	}
	res.Err = err
	if j.attempts <= j.MaxRetries {
		res.Status = Retrying
		delay := j.Retry.NextDelay(j.attempts - 1)
		j.RunAt = time.Now().Add(delay)
		heap.Push(&s.queue, j)
		return
	}
	res.Status = Failed
}

func (s *Scheduler) GetResult(id string) (JobResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.results[id]
	if !ok {
		return JobResult{}, false
	}
	return *r, true
}

// WaitForResult polls until job id reaches a terminal state (Succeeded or
// Failed) or the timeout elapses.
func (s *Scheduler) WaitForResult(id string, timeout time.Duration) (JobResult, bool) {
	deadline := time.Now().Add(timeout)
	for {
		if r, ok := s.GetResult(id); ok && (r.Status == Succeeded || r.Status == Failed) {
			return r, true
		}
		if time.Now().After(deadline) {
			return JobResult{}, false
		}
		time.Sleep(2 * time.Millisecond)
	}
}
