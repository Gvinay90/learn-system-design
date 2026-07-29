import threading

import pytest

from kafka_lld import (
    Broker,
    PartitionNotFoundError,
    TopicNotFoundError,
)


def test_produce_consume_in_order():
    b = Broker()
    b.create_topic("orders", 1)

    for i, v in enumerate(["a", "b", "c"]):
        partition_id, offset = b.produce("orders", "k1", v)
        assert partition_id == 0
        assert offset == i

    messages = b.consume("g1", "orders", 0, 10)
    assert len(messages) == 3
    for i, m in enumerate(messages):
        assert m.offset == i

    more = b.consume("g1", "orders", 0, 10)
    assert len(more) == 0


def test_consumer_groups_track_offsets_independently():
    b = Broker()
    b.create_topic("orders", 1)
    for v in ["a", "b", "c"]:
        b.produce("orders", "", v)

    g1_messages = b.consume("group-1", "orders", 0, 2)
    assert len(g1_messages) == 2

    g2_messages = b.consume("group-2", "orders", 0, 10)
    assert len(g2_messages) == 3

    g1_rest = b.consume("group-1", "orders", 0, 10)
    assert len(g1_rest) == 1
    assert g1_rest[0].offset == 2


def test_edge_cases():
    b = Broker()
    b.create_topic("orders", 1)
    b.produce("orders", "", "only-message")

    messages = b.consume("g1", "orders", 0, 10)
    assert len(messages) == 1

    past = b.consume("g1", "orders", 0, 10)
    assert len(past) == 0

    with pytest.raises(TopicNotFoundError):
        b.produce("unknown-topic", "k", "v")

    with pytest.raises(TopicNotFoundError):
        b.consume("g1", "unknown-topic", 0, 10)

    with pytest.raises(PartitionNotFoundError):
        b.consume("g1", "orders", 5, 10)


def test_concurrent_produce_into_same_partition():
    """Asserts many threads racing to append to the same partition never
    lose a message or assign a duplicate offset — the lock in
    Partition.append must serialize them."""
    b = Broker()
    b.create_topic("orders", 1)

    n = 500
    offsets = []
    offsets_lock = threading.Lock()

    def worker():
        _, offset = b.produce("orders", "", "v")
        with offsets_lock:
            offsets.append(offset)

    threads = [threading.Thread(target=worker) for _ in range(n)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    seen = set(offsets)
    assert len(seen) == n
    for i in range(n):
        assert i in seen
