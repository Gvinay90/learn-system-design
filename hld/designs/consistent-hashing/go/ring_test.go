package consistenthashing

import (
	"fmt"
	"math"
	"testing"
)

func TestGetIsDeterministicForSameKey(t *testing.T) {
	r := NewRing(100)
	r.AddNode("nodeA")
	r.AddNode("nodeB")
	r.AddNode("nodeC")

	node, ok := r.Get("user:1234")
	if !ok {
		t.Fatalf("expected a node to own the key")
	}
	for i := 0; i < 10; i++ {
		got, _ := r.Get("user:1234")
		if got != node {
			t.Fatalf("expected repeated Get for same key to return same node, got %q then %q", node, got)
		}
	}
}

func TestGetOnEmptyRing(t *testing.T) {
	r := NewRing(100)
	if _, ok := r.Get("anything"); ok {
		t.Fatalf("expected no owner on an empty ring")
	}
}

func TestDistributionIsReasonablyEven(t *testing.T) {
	r := NewRing(200) // enough virtual nodes to smooth distribution
	nodes := []string{"n1", "n2", "n3", "n4"}
	for _, n := range nodes {
		r.AddNode(n)
	}

	counts := make(map[string]int)
	const numKeys = 100_000
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key-%d", i)
		node, _ := r.Get(key)
		counts[node]++
	}

	expected := float64(numKeys) / float64(len(nodes))
	for _, n := range nodes {
		got := float64(counts[n])
		deviation := math.Abs(got-expected) / expected
		if deviation > 0.15 { // allow 15% deviation from perfectly even
			t.Errorf("node %q got %d keys (expected ~%.0f, %.1f%% deviation) — distribution too skewed",
				n, counts[n], expected, deviation*100)
		}
	}
}

func TestMinimalRemappingOnNodeAdd(t *testing.T) {
	r := NewRing(100)
	nodes := []string{"n1", "n2", "n3", "n4"}
	for _, n := range nodes {
		r.AddNode(n)
	}

	const numKeys = 10_000
	before := make(map[string]string, numKeys)
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key-%d", i)
		node, _ := r.Get(key)
		before[key] = node
	}

	r.AddNode("n5") // 5th node added

	moved := 0
	for key, oldNode := range before {
		newNode, _ := r.Get(key)
		if newNode != oldNode {
			moved++
		}
	}

	// With naive mod-N hashing, adding a 5th node to 4 would remap ~80% of
	// keys. Consistent hashing should only remap roughly 1/5 (~20%) of them
	// — the ones that now land on the new node's virtual points.
	fractionMoved := float64(moved) / float64(numKeys)
	if fractionMoved > 0.35 {
		t.Errorf("expected roughly 1/5 (~20%%) of keys to move on adding a 5th node, got %.1f%% moved",
			fractionMoved*100)
	}
	if moved == 0 {
		t.Errorf("expected some keys to move to the new node, got 0")
	}
}

func TestMinimalRemappingOnNodeRemove(t *testing.T) {
	r := NewRing(100)
	nodes := []string{"n1", "n2", "n3", "n4", "n5"}
	for _, n := range nodes {
		r.AddNode(n)
	}

	const numKeys = 10_000
	before := make(map[string]string, numKeys)
	for i := 0; i < numKeys; i++ {
		key := fmt.Sprintf("key-%d", i)
		node, _ := r.Get(key)
		before[key] = node
	}

	r.RemoveNode("n3")

	moved := 0
	for key, oldNode := range before {
		newNode, _ := r.Get(key)
		if newNode != oldNode {
			moved++
			if newNode == "n3" {
				t.Fatalf("key %q remapped to removed node n3", key)
			}
		} else if oldNode == "n3" {
			t.Fatalf("key %q still on removed node n3 after removal", key)
		}
	}

	// Only keys that were owned by n3 should move; that's ~1/5 of the keys.
	fractionMoved := float64(moved) / float64(numKeys)
	if fractionMoved > 0.35 {
		t.Errorf("expected roughly 1/5 (~20%%) of keys to move on removing a node, got %.1f%% moved",
			fractionMoved*100)
	}
}

func TestRemoveThenGetRedistributesToRemainingNodes(t *testing.T) {
	r := NewRing(50)
	r.AddNode("only-node")
	r.RemoveNode("only-node")

	if _, ok := r.Get("some-key"); ok {
		t.Fatalf("expected no owner after removing the only node")
	}

	r.AddNode("new-node")
	node, ok := r.Get("some-key")
	if !ok || node != "new-node" {
		t.Fatalf("expected some-key to route to new-node, got %q (ok=%v)", node, ok)
	}
}
