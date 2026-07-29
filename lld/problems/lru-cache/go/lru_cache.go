// Package lrucache implements a fixed-capacity LRU cache backed by a
// hash map plus an intrusive doubly linked list, giving O(1) get/put/evict.
package lrucache

import (
	"errors"
	"sync"
)

var ErrInvalidCapacity = errors.New("capacity must be positive")

type node[K comparable, V any] struct {
	key        K
	value      V
	prev, next *node[K, V]
}

// LRUCache is safe for concurrent use; a single mutex guards the map and
// list together since both must stay in sync on every get/put.
type LRUCache[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	items    map[K]*node[K, V]
	head     *node[K, V] // sentinel; head.next is the most-recently-used node
	tail     *node[K, V] // sentinel; tail.prev is the least-recently-used node
}

func NewLRUCache[K comparable, V any](capacity int) (*LRUCache[K, V], error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	head := &node[K, V]{}
	tail := &node[K, V]{}
	head.next = tail
	tail.prev = head
	return &LRUCache[K, V]{
		capacity: capacity,
		items:    make(map[K]*node[K, V], capacity),
		head:     head,
		tail:     tail,
	}, nil
}

func (c *LRUCache[K, V]) unlink(n *node[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// pushFront inserts n right after the head sentinel (the MRU position).
func (c *LRUCache[K, V]) pushFront(n *node[K, V]) {
	n.prev = c.head
	n.next = c.head.next
	c.head.next.prev = n
	c.head.next = n
}

func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	n, ok := c.items[key]
	if !ok {
		var zero V
		return zero, false
	}
	c.unlink(n)
	c.pushFront(n)
	return n.value, true
}

func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if n, ok := c.items[key]; ok {
		n.value = value
		c.unlink(n)
		c.pushFront(n)
		return
	}

	n := &node[K, V]{key: key, value: value}
	c.items[key] = n
	c.pushFront(n)

	if len(c.items) > c.capacity {
		lru := c.tail.prev
		c.unlink(lru)
		delete(c.items, lru.key)
	}
}

func (c *LRUCache[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}
