"""Producer-consumer worker pool LLD primitive — Python reference implementation.

Multiple producer threads push items into a bounded queue.Queue; a fixed
pool of worker threads drains it concurrently. See ../README.md for the
design writeup.
"""
from __future__ import annotations

import queue
import threading
from dataclasses import dataclass
from typing import Callable


@dataclass(frozen=True)
class Item:
    producer_id: int
    seq: int

    @property
    def id(self) -> str:
        return f"p{self.producer_id}-{self.seq}"


_POISON = Item(-1, -1)


class BoundedPipeline:
    """Runs producers and a worker pool over a bounded queue.Queue.

    Producers push items until exhausted; one poison-pill sentinel per
    worker is enqueued afterwards so every worker thread exits cleanly.
    """

    def __init__(self, num_workers: int, buffer_size: int):
        self.num_workers = num_workers
        self.buffer_size = buffer_size

    def run(self, num_producers: int, items_per_producer: int, handle: Callable[[Item], None]) -> None:
        buffer: "queue.Queue[Item]" = queue.Queue(maxsize=self.buffer_size)

        def produce(producer_id: int) -> None:
            for seq in range(items_per_producer):
                buffer.put(Item(producer_id, seq))

        def consume() -> None:
            while True:
                item = buffer.get()
                if item is _POISON:
                    buffer.task_done()
                    return
                handle(item)
                buffer.task_done()

        producers = [threading.Thread(target=produce, args=(pid,)) for pid in range(num_producers)]
        workers = [threading.Thread(target=consume) for _ in range(self.num_workers)]

        for w in workers:
            w.start()
        for p in producers:
            p.start()
        for p in producers:
            p.join()

        for _ in workers:
            buffer.put(_POISON)
        for w in workers:
            w.join()


def _demo() -> None:
    consumed: set[str] = set()
    lock = threading.Lock()

    def handle(item: Item) -> None:
        with lock:
            consumed.add(item.id)

    pipeline = BoundedPipeline(num_workers=4, buffer_size=8)
    pipeline.run(num_producers=3, items_per_producer=20, handle=handle)
    print(f"Consumed {len(consumed)} unique items via producer-consumer pipeline")


if __name__ == "__main__":
    _demo()
