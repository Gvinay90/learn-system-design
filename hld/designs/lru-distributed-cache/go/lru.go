// Package lrudistributedcache implements the single-node LRU eviction core
// (HashMap + doubly linked list, O(1) Get/Put) that each shard of the
// distributed cache design in the README runs internally. The distributed
// design layers consistent-hash partitioning and replication on top of many
// copies of exactly this data structure — see README.md for the full design.
package lrudistributedcache

import "sync"

// node is a doubly linked list node holding one cache entry.
type node struct {
	key, value string
	prev, next *node
}

// LRU is a fixed-capacity, thread-safe least-recently-used cache. Get and
// Put are O(1): the hash map gives O(1) lookup, and the doubly linked list
// gives O(1) move-to-front / eviction without shifting any other elements.
type LRU struct {
	capacity int

	mu    sync.Mutex
	items map[string]*node
	head  *node // most recently used (sentinel)
	tail  *node // least recently used (sentinel)
}

// New creates an LRU cache that holds at most capacity entries. capacity
// must be >= 1.
func New(capacity int) *LRU {
	if capacity < 1 {
		capacity = 1
	}
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head
	return &LRU{
		capacity: capacity,
		items:    make(map[string]*node, capacity),
		head:     head,
		tail:     tail,
	}
}

// Get returns the value for key and marks it most-recently-used. The bool
// reports whether the key was present.
func (l *LRU) Get(key string) (string, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	n, ok := l.items[key]
	if !ok {
		return "", false
	}
	l.moveToFront(n)
	return n.value, true
}

// Put inserts or updates key's value, marking it most-recently-used. If the
// cache is at capacity and key is new, the least-recently-used entry (the
// node just before the tail sentinel) is evicted.
func (l *LRU) Put(key, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if n, ok := l.items[key]; ok {
		n.value = value
		l.moveToFront(n)
		return
	}

	if len(l.items) >= l.capacity {
		l.evictLRU()
	}

	n := &node{key: key, value: value}
	l.items[key] = n
	l.insertAtFront(n)
}

// Len returns the current number of entries.
func (l *LRU) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.items)
}

// --- internal list operations (caller must hold l.mu) ---

func (l *LRU) moveToFront(n *node) {
	l.unlink(n)
	l.insertAtFront(n)
}

func (l *LRU) unlink(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (l *LRU) insertAtFront(n *node) {
	n.prev = l.head
	n.next = l.head.next
	l.head.next.prev = n
	l.head.next = n
}

func (l *LRU) evictLRU() {
	lru := l.tail.prev
	if lru == l.head {
		return // empty
	}
	l.unlink(lru)
	delete(l.items, lru.key)
}
