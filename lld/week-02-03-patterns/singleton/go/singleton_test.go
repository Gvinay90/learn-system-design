package singleton

import (
	"sync"
	"testing"
)

func TestGetInstanceReturnsSameInstance(t *testing.T) {
	resetForTest()
	a := GetInstance()
	b := GetInstance()
	if a != b {
		t.Fatalf("expected same pointer, got different instances")
	}
	if a.ID() != b.ID() {
		t.Fatalf("expected same id, got %d and %d", a.ID(), b.ID())
	}
}

func TestSetAndGet(t *testing.T) {
	resetForTest()
	cfg := GetInstance()
	cfg.Set("env", "production")

	v, ok := cfg.Get("env")
	if !ok || v != "production" {
		t.Fatalf("expected env=production, got %q ok=%v", v, ok)
	}

	if _, ok := cfg.Get("missing"); ok {
		t.Fatalf("expected missing key to be absent")
	}
}

// TestConcurrentFirstAccess spins up many goroutines racing to call
// GetInstance for the first time; sync.Once must guarantee exactly one
// AppConfig gets constructed, so every goroutine must observe the same id.
func TestConcurrentFirstAccess(t *testing.T) {
	resetForTest()

	const goroutines = 200
	var wg sync.WaitGroup
	ids := make([]int, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			ids[idx] = GetInstance().ID()
		}(i)
	}
	wg.Wait()

	first := ids[0]
	for i, id := range ids {
		if id != first {
			t.Fatalf("goroutine %d observed id %d, expected %d (multiple instances constructed)", i, id, first)
		}
	}
}

func TestDemo(t *testing.T) {
	resetForTest()
	cfg := GetInstance()
	cfg.Set("region", "us-east-1")
	v, _ := cfg.Get("region")
	if v != "us-east-1" {
		t.Fatalf("demo failed, got %q", v)
	}
}
