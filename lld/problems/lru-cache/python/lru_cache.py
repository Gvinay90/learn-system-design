"""LRU Cache LLD — Python reference implementation.

Hash map + hand-rolled doubly linked list (sentinel head/tail) for O(1)
get/put/evict. See ../README.md for the design writeup.
"""
from __future__ import annotations

import threading
from typing import Dict, Generic, Optional, Tuple, TypeVar

K = TypeVar("K")
V = TypeVar("V")


class Node(Generic[K, V]):
    def __init__(self, key: Optional[K], value: Optional[V]):
        self.key = key
        self.value = value
        self.prev: Optional["Node[K, V]"] = None
        self.next: Optional["Node[K, V]"] = None


class LRUCache(Generic[K, V]):
    def __init__(self, capacity: int):
        if capacity <= 0:
            raise ValueError("capacity must be positive")
        self._capacity = capacity
        self._items: Dict[K, Node[K, V]] = {}
        self._lock = threading.Lock()
        # Sentinels avoid None-checks at the ends of the list.
        self._head: Node[K, V] = Node(None, None)
        self._tail: Node[K, V] = Node(None, None)
        self._head.next = self._tail
        self._tail.prev = self._head

    def _unlink(self, node: Node[K, V]) -> None:
        node.prev.next = node.next
        node.next.prev = node.prev

    def _push_front(self, node: Node[K, V]) -> None:
        node.prev = self._head
        node.next = self._head.next
        self._head.next.prev = node
        self._head.next = node

    def get(self, key: K) -> Tuple[Optional[V], bool]:
        with self._lock:
            node = self._items.get(key)
            if node is None:
                return None, False
            self._unlink(node)
            self._push_front(node)
            return node.value, True

    def put(self, key: K, value: V) -> None:
        with self._lock:
            node = self._items.get(key)
            if node is not None:
                node.value = value
                self._unlink(node)
                self._push_front(node)
                return

            node = Node(key, value)
            self._items[key] = node
            self._push_front(node)

            if len(self._items) > self._capacity:
                lru = self._tail.prev
                self._unlink(lru)
                del self._items[lru.key]

    def __len__(self) -> int:
        with self._lock:
            return len(self._items)

    def _list_length(self) -> int:
        with self._lock:
            count = 0
            node = self._head.next
            while node is not self._tail:
                count += 1
                node = node.next
            return count


def _demo() -> None:
    cache: LRUCache[str, int] = LRUCache(2)
    cache.put("a", 1)
    cache.put("b", 2)
    print(f"get(a) = {cache.get('a')}")
    cache.put("c", 3)
    print(f"get(b) after eviction = {cache.get('b')}")
    print(f"get(c) = {cache.get('c')}")


if __name__ == "__main__":
    _demo()
