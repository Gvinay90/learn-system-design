"""Rate Limiter LLD — Python reference implementation.

Two interchangeable per-client algorithms (token bucket, sliding window)
behind a common RateLimiter interface, plus a RateLimiterRegistry that maps
a client class (free/paid/enterprise) to its configured limiter. See
../README.md for the design writeup.
"""
from __future__ import annotations

import threading
import time
from abc import ABC, abstractmethod
from collections import deque
from typing import Callable, Deque, Dict, Optional


class RateLimiter(ABC):
    """Common strategy interface every algorithm implements, so calling
    code never branches on which algorithm is active."""

    @abstractmethod
    def allow(self, client_id: str) -> bool:
        """Decides whether a request from client_id is permitted right now."""
        raise NotImplementedError


class _Bucket:
    __slots__ = ("tokens", "last_refill")

    def __init__(self, tokens: float, last_refill: float):
        self.tokens = tokens
        self.last_refill = last_refill


class TokenBucketLimiter(RateLimiter):
    """Allows bursts up to capacity tokens, refilling at refill_rate
    tokens/second. Refill is computed lazily on each allow() call rather
    than via a background thread per client."""

    def __init__(self, capacity: float, refill_rate: float, clock: Callable[[], float] = time.monotonic):
        self._capacity = capacity
        self._refill_rate = refill_rate  # tokens per second
        self._clock = clock
        self._buckets: Dict[str, _Bucket] = {}
        self._lock = threading.Lock()

    def _refill(self, bucket: _Bucket, now: float) -> None:
        elapsed = now - bucket.last_refill
        if elapsed <= 0:
            return
        bucket.tokens = min(self._capacity, bucket.tokens + elapsed * self._refill_rate)
        bucket.last_refill = now

    def allow(self, client_id: str) -> bool:
        with self._lock:
            now = self._clock()
            bucket = self._buckets.get(client_id)
            if bucket is None:
                bucket = _Bucket(self._capacity, now)
                self._buckets[client_id] = bucket
            else:
                self._refill(bucket, now)

            if bucket.tokens < 1:
                return False
            bucket.tokens -= 1
            return True


class _Window:
    __slots__ = ("timestamps",)

    def __init__(self) -> None:
        self.timestamps: Deque[float] = deque()


class SlidingWindowLimiter(RateLimiter):
    """Allows at most `limit` requests in any trailing `window_seconds`
    duration (a sliding log), giving a smooth, hard cap with no boundary
    burst allowance."""

    def __init__(self, limit: int, window_seconds: float, clock: Callable[[], float] = time.monotonic):
        self._limit = limit
        self._window = window_seconds
        self._clock = clock
        self._windows: Dict[str, _Window] = {}
        self._lock = threading.Lock()

    def allow(self, client_id: str) -> bool:
        with self._lock:
            now = self._clock()
            window = self._windows.setdefault(client_id, _Window())

            cutoff = now - self._window
            while window.timestamps and window.timestamps[0] <= cutoff:
                window.timestamps.popleft()

            if len(window.timestamps) >= self._limit:
                return False
            window.timestamps.append(now)
            return True


class RateLimiterRegistry:
    """Maps a client class (e.g. free/paid/enterprise) to the RateLimiter
    instance configured for that tier, so different tiers can even run
    different algorithms."""

    def __init__(self) -> None:
        self._per_client_class: Dict[str, RateLimiter] = {}
        self._lock = threading.Lock()

    def register(self, client_class: str, limiter: RateLimiter) -> None:
        with self._lock:
            self._per_client_class[client_class] = limiter

    def get_limiter(self, client_class: str) -> Optional[RateLimiter]:
        with self._lock:
            return self._per_client_class.get(client_class)


def _demo() -> None:
    token_bucket = TokenBucketLimiter(capacity=2, refill_rate=1)
    print(f"token bucket allow(a) = {token_bucket.allow('a')}")
    print(f"token bucket allow(a) = {token_bucket.allow('a')}")
    print(f"token bucket allow(a) [should block] = {token_bucket.allow('a')}")

    sliding_window = SlidingWindowLimiter(limit=2, window_seconds=1)
    print(f"sliding window allow(b) = {sliding_window.allow('b')}")
    print(f"sliding window allow(b) = {sliding_window.allow('b')}")
    print(f"sliding window allow(b) [should block] = {sliding_window.allow('b')}")


if __name__ == "__main__":
    _demo()
