package lrudistributedcache

import (
	"fmt"
	"sync"
	"testing"
)

func TestEvictionOrderIsLeastRecentlyUsed(t *testing.T) {
	c := New(3)
	c.Put("a", "1")
	c.Put("b", "2")
	c.Put("c", "3")

	// Touch "a" so "b" becomes the least recently used.
	if _, ok := c.Get("a"); !ok {
		t.Fatalf("expected a to be present")
	}

	c.Put("d", "4") // should evict "b" (least recently used)

	if _, ok := c.Get("b"); ok {
		t.Fatalf("expected b to have been evicted")
	}
	if v, ok := c.Get("a"); !ok || v != "1" {
		t.Fatalf("expected a to survive with value 1, got %q (ok=%v)", v, ok)
	}
	if v, ok := c.Get("c"); !ok || v != "3" {
		t.Fatalf("expected c to survive with value 3, got %q (ok=%v)", v, ok)
	}
	if v, ok := c.Get("d"); !ok || v != "4" {
		t.Fatalf("expected d to be present with value 4, got %q (ok=%v)", v, ok)
	}
	if c.Len() != 3 {
		t.Fatalf("expected capacity to be maintained at 3, got %d", c.Len())
	}
}

func TestCapacityOneEdgeCase(t *testing.T) {
	c := New(1)
	c.Put("a", "1")
	if v, ok := c.Get("a"); !ok || v != "1" {
		t.Fatalf("expected a=1, got %q (ok=%v)", v, ok)
	}

	c.Put("b", "2") // should evict "a" immediately
	if _, ok := c.Get("a"); ok {
		t.Fatalf("expected a to be evicted in capacity-1 cache")
	}
	if v, ok := c.Get("b"); !ok || v != "2" {
		t.Fatalf("expected b=2, got %q (ok=%v)", v, ok)
	}
	if c.Len() != 1 {
		t.Fatalf("expected len 1, got %d", c.Len())
	}
}

func TestUpdatingExistingKeyDoesNotEvict(t *testing.T) {
	c := New(2)
	c.Put("a", "1")
	c.Put("b", "2")
	c.Put("a", "updated") // update, not a new insert

	if c.Len() != 2 {
		t.Fatalf("expected len to stay 2 after updating existing key, got %d", c.Len())
	}
	if v, ok := c.Get("a"); !ok || v != "updated" {
		t.Fatalf("expected a=updated, got %q (ok=%v)", v, ok)
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatalf("expected b to still be present")
	}
}

func TestGetMissingKey(t *testing.T) {
	c := New(2)
	if _, ok := c.Get("missing"); ok {
		t.Fatalf("expected miss for absent key")
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := New(50)
	var wg sync.WaitGroup

	// Many goroutines hammering Put/Get concurrently; the race detector
	// (`go test -race`) plus a sane final state is the real assertion here.
	for g := 0; g < 20; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				key := fmt.Sprintf("g%d-k%d", g, i%10)
				c.Put(key, "v")
				c.Get(key)
			}
		}(g)
	}
	wg.Wait()

	if c.Len() > 50 {
		t.Fatalf("expected len to stay within capacity 50, got %d", c.Len())
	}
}
