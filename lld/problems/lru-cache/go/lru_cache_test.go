package lrucache

import (
	"sync"
	"testing"
)

func TestGetPutUpdateInPlace(t *testing.T) {
	c, err := NewLRUCache[string, int](2)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	c.Put("a", 1)
	c.Put("b", 2)

	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("expected a=1, got %v ok=%v", v, ok)
	}

	c.Put("a", 100)
	if v, ok := c.Get("a"); !ok || v != 100 {
		t.Fatalf("expected a=100 after update, got %v ok=%v", v, ok)
	}
	if c.Len() != 2 {
		t.Fatalf("expected len 2, got %d", c.Len())
	}
}

func TestEvictsActualLRU(t *testing.T) {
	c, _ := NewLRUCache[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)

	// Touch "a" so "b" becomes the least-recently-used entry.
	c.Get("a")
	c.Put("c", 3)

	if _, ok := c.Get("b"); ok {
		t.Fatalf("expected b to be evicted")
	}
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("expected a to survive eviction, got %v ok=%v", v, ok)
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Fatalf("expected c to survive eviction, got %v ok=%v", v, ok)
	}
	if c.Len() != 2 {
		t.Fatalf("expected len 2, got %d", c.Len())
	}
}

func TestMissingKeyAndCapacityOne(t *testing.T) {
	c, _ := NewLRUCache[string, int](1)

	if _, ok := c.Get("missing"); ok {
		t.Fatalf("expected miss for absent key")
	}
	if c.Len() != 0 {
		t.Fatalf("expected len 0 after a miss, got %d", c.Len())
	}

	c.Put("a", 1)
	c.Put("b", 2)
	if _, ok := c.Get("a"); ok {
		t.Fatalf("expected a evicted at capacity 1")
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatalf("expected b=2, got %v ok=%v", v, ok)
	}
}

// TestConcurrentAccess runs many goroutines performing puts/gets on a shared
// cache and asserts no panic occurs and the final size never exceeds capacity.
func TestConcurrentAccess(t *testing.T) {
	const capacity = 50
	const goroutines = 100

	c, _ := NewLRUCache[int, int](capacity)
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Put(n, n*n)
			c.Get(n)
			c.Get(n - 1)
		}(i)
	}
	wg.Wait()

	if c.Len() > capacity {
		t.Fatalf("expected size <= capacity %d, got %d", capacity, c.Len())
	}
	if len(c.items) != c.Len() {
		t.Fatalf("map size %d and reported len %d disagree", len(c.items), c.Len())
	}

	count := 0
	for n := c.head.next; n != c.tail; n = n.next {
		count++
	}
	if count != c.Len() {
		t.Fatalf("linked list length %d does not match cache len %d", count, c.Len())
	}
}
