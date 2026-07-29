import threading

from concurrency import BoundedPipeline


def test_no_lost_or_duplicated_items():
    num_producers = 5
    items_per_producer = 200
    num_workers = 8
    buffer_size = 16

    counts: dict[str, int] = {}
    lock = threading.Lock()

    def handle(item):
        with lock:
            counts[item.id] = counts.get(item.id, 0) + 1

    pipeline = BoundedPipeline(num_workers=num_workers, buffer_size=buffer_size)
    pipeline.run(num_producers, items_per_producer, handle)

    want_total = num_producers * items_per_producer
    assert len(counts) == want_total

    for pid in range(num_producers):
        for seq in range(items_per_producer):
            item_id = f"p{pid}-{seq}"
            assert item_id in counts, f"item {item_id} was never consumed"
            assert counts[item_id] == 1, f"item {item_id} consumed {counts[item_id]} times"


def test_single_producer_single_worker():
    consumed: list[str] = []
    lock = threading.Lock()

    def handle(item):
        with lock:
            consumed.append(item.id)

    pipeline = BoundedPipeline(num_workers=1, buffer_size=1)
    pipeline.run(1, 50, handle)

    assert len(consumed) == 50
    assert len(set(consumed)) == 50
