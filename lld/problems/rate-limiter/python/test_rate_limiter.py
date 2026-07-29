import threading

import pytest

from rate_limiter import RateLimiterRegistry, SlidingWindowLimiter, TokenBucketLimiter


class FakeClock:
    """Deterministic clock for tests; avoids sleeping real time."""

    def __init__(self, start: float = 0.0):
        self._now = start

    def __call__(self) -> float:
        return self._now

    def advance(self, delta: float) -> None:
        self._now += delta


@pytest.mark.parametrize(
    "capacity, requests, want_allowed",
    [
        (3, 5, 3),
        (1, 4, 1),
    ],
)
def test_token_bucket_allows_burst_up_to_capacity_then_blocks(capacity, requests, want_allowed):
    clock = FakeClock()
    limiter = TokenBucketLimiter(capacity=capacity, refill_rate=1, clock=clock)

    allowed = sum(1 for _ in range(requests) if limiter.allow("client-a"))
    assert allowed == want_allowed


def test_token_bucket_refills_over_time():
    clock = FakeClock()
    limiter = TokenBucketLimiter(capacity=2, refill_rate=1, clock=clock)  # 1 token/sec

    assert limiter.allow("client-a")
    assert limiter.allow("client-a")
    assert not limiter.allow("client-a")

    clock.advance(1.0)  # refills exactly 1 token
    assert limiter.allow("client-a")
    assert not limiter.allow("client-a")


def test_sliding_window_allows_up_to_max_then_blocks_then_slides():
    clock = FakeClock()
    limiter = SlidingWindowLimiter(limit=3, window_seconds=1, clock=clock)

    for _ in range(3):
        assert limiter.allow("client-a")
    assert not limiter.allow("client-a")

    clock.advance(1.1)  # slide window fully past recorded timestamps
    assert limiter.allow("client-a")


def test_per_client_isolation():
    tb_clock = FakeClock()
    tb = TokenBucketLimiter(capacity=1, refill_rate=1, clock=tb_clock)
    assert tb.allow("client-a")
    assert not tb.allow("client-a")
    assert tb.allow("client-b")

    sw_clock = FakeClock()
    sw = SlidingWindowLimiter(limit=1, window_seconds=1, clock=sw_clock)
    assert sw.allow("client-a")
    assert not sw.allow("client-a")
    assert sw.allow("client-b")


def test_rate_limiter_registry():
    registry = RateLimiterRegistry()
    free = TokenBucketLimiter(capacity=1, refill_rate=1)
    paid = SlidingWindowLimiter(limit=100, window_seconds=1)

    registry.register("free", free)
    registry.register("paid", paid)

    assert registry.get_limiter("free") is free
    assert registry.get_limiter("enterprise") is None


def test_concurrent_allow_does_not_exceed_capacity():
    capacity = 10
    thread_count = 100
    limiter = TokenBucketLimiter(capacity=capacity, refill_rate=0)  # no refill during the burst
    allowed_count = 0
    lock = threading.Lock()

    def worker():
        nonlocal allowed_count
        if limiter.allow("client-a"):
            with lock:
                allowed_count += 1

    threads = [threading.Thread(target=worker) for _ in range(thread_count)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert allowed_count == capacity
