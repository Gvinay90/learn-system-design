// Package consistenthashing implements a consistent hash ring with virtual
// nodes, the building block referenced by sharding, load-balancing, and
// distributed-cache designs throughout this repo. See the README for the
// full write-up (why virtual nodes, how much data moves on add/remove).
package consistenthashing

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

// Ring is a consistent hash ring. Each real node is mapped to VirtualNodes
// points on the ring (via "<node>#<replica>" hashing) so that load spreads
// evenly across nodes instead of depending on the luck of a single hash
// point per node.
type Ring struct {
	VirtualNodes int

	mu       sync.RWMutex
	sorted   []uint32          // sorted hash points on the ring
	hashToID map[uint32]string // ring point -> real node ID
	nodes    map[string]bool   // set of real node IDs currently in the ring
}

// NewRing creates a ring with the given number of virtual nodes (replicas)
// per real node. More virtual nodes -> smoother load distribution, at the
// cost of more memory and a slightly larger sorted-slice lookup.
func NewRing(virtualNodes int) *Ring {
	if virtualNodes <= 0 {
		virtualNodes = 100
	}
	return &Ring{
		VirtualNodes: virtualNodes,
		hashToID:     make(map[uint32]string),
		nodes:        make(map[string]bool),
	}
}

func hashKey(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

// AddNode adds a real node to the ring, creating VirtualNodes points for it.
func (r *Ring) AddNode(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.nodes[node] {
		return // already present, no-op
	}
	r.nodes[node] = true

	for i := 0; i < r.VirtualNodes; i++ {
		point := hashKey(fmt.Sprintf("%s#%d", node, i))
		r.hashToID[point] = node
		r.sorted = append(r.sorted, point)
	}
	sort.Slice(r.sorted, func(i, j int) bool { return r.sorted[i] < r.sorted[j] })
}

// RemoveNode removes a real node and all its virtual points from the ring.
func (r *Ring) RemoveNode(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.nodes[node] {
		return
	}
	delete(r.nodes, node)

	kept := r.sorted[:0:0]
	for _, point := range r.sorted {
		if r.hashToID[point] == node {
			delete(r.hashToID, point)
			continue
		}
		kept = append(kept, point)
	}
	r.sorted = kept
}

// Get returns the node responsible for key: hash the key, then walk
// clockwise (the first virtual point >= the key's hash; wrap to the start
// of the slice if the key hashes past the last point).
func (r *Ring) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.sorted) == 0 {
		return "", false
	}

	h := hashKey(key)
	idx := sort.Search(len(r.sorted), func(i int) bool { return r.sorted[i] >= h })
	if idx == len(r.sorted) {
		idx = 0 // wrap around the ring
	}
	return r.hashToID[r.sorted[idx]], true
}

// Nodes returns the current set of real node IDs in the ring.
func (r *Ring) Nodes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.nodes))
	for n := range r.nodes {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
