import threading

from lru_cache import LRUCache


def test_get_put_update_in_place():
    cache = LRUCache(2)
    cache.put("a", 1)
    cache.put("b", 2)

    value, found = cache.get("a")
    assert found and value == 1

    cache.put("a", 100)
    value, found = cache.get("a")
    assert found and value == 100
    assert len(cache) == 2


def test_evicts_actual_lru():
    cache = LRUCache(2)
    cache.put("a", 1)
    cache.put("b", 2)

    cache.get("a")  # touch a so b becomes least-recently-used
    cache.put("c", 3)

    _, found_b = cache.get("b")
    assert not found_b

    value_a, found_a = cache.get("a")
    assert found_a and value_a == 1

    value_c, found_c = cache.get("c")
    assert found_c and value_c == 3
    assert len(cache) == 2


def test_missing_key_and_capacity_one():
    cache = LRUCache(1)
    _, found = cache.get("missing")
    assert not found
    assert len(cache) == 0

    cache.put("a", 1)
    cache.put("b", 2)
    _, found_a = cache.get("a")
    assert not found_a
    value_b, found_b = cache.get("b")
    assert found_b and value_b == 2


def test_concurrent_access():
    capacity = 50
    thread_count = 100
    cache = LRUCache(capacity)

    def worker(n: int) -> None:
        cache.put(n, n * n)
        cache.get(n)
        cache.get(n - 1)

    threads = [threading.Thread(target=worker, args=(i,)) for i in range(thread_count)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert len(cache) <= capacity
    assert len(cache) == cache._list_length()
