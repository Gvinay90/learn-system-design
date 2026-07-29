"""Payment Gateway LLD — Python reference implementation.

Idempotency-key deduplication, a pluggable PaymentProcessor (Strategy) with
retries, and a simple in-memory ledger. See ../README.md for the design
writeup.
"""
from __future__ import annotations

import threading
import time
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum, auto
from typing import List, Optional, Protocol


class PaymentStatus(Enum):
    PENDING = auto()
    SUCCESS = auto()
    FAILED = auto()


@dataclass
class PaymentRequest:
    idempotency_key: str
    payer_id: str
    payee_id: str
    amount: float
    currency: str


@dataclass
class Attempt:
    number: int
    success: bool
    error: Optional[str]
    at: datetime


@dataclass
class PaymentResult:
    id: str
    request: PaymentRequest
    status: PaymentStatus = PaymentStatus.PENDING
    created_at: datetime = field(default_factory=datetime.now)
    attempts: List[Attempt] = field(default_factory=list)


@dataclass
class LedgerEntry:
    payment_id: str
    payer_id: str
    payee_id: str
    amount: float
    at: datetime


class Ledger:
    """One row per successful payment (debit-payer/credit-payee collapsed
    into a single entry to keep the example focused on idempotency/retries).
    """

    def __init__(self) -> None:
        self._entries: List[LedgerEntry] = []
        self._lock = threading.Lock()

    def record(self, entry: LedgerEntry) -> None:
        with self._lock:
            self._entries.append(entry)

    def entries(self) -> List[LedgerEntry]:
        with self._lock:
            return list(self._entries)


class PaymentProcessor(Protocol):
    def process(self, request: PaymentRequest) -> None:
        """Raises PaymentProcessingError on failure."""
        ...


class PaymentProcessingError(Exception):
    pass


class FakePaymentProcessor:
    def __init__(self, fail_times: int = 0, always_fail: bool = False):
        self.fail_times = fail_times
        self.always_fail = always_fail
        self.call_count = 0
        self._lock = threading.Lock()

    def process(self, request: PaymentRequest) -> None:
        with self._lock:
            self.call_count += 1
            count = self.call_count
        if self.always_fail:
            raise PaymentProcessingError("simulated permanent failure")
        if count <= self.fail_times:
            raise PaymentProcessingError("simulated transient failure")


@dataclass
class RetryPolicy:
    max_attempts: int
    delay_seconds: float

    def delay_for(self, attempt: int) -> float:
        return attempt * self.delay_seconds


class _IdempotencyEntry:
    def __init__(self) -> None:
        self.done = threading.Event()
        self.result: Optional[PaymentResult] = None


class IdempotencyStore:
    def __init__(self) -> None:
        self._entries: dict[str, _IdempotencyEntry] = {}
        self._lock = threading.Lock()

    def reserve_or_wait(self, key: str) -> tuple[Optional[PaymentResult], bool]:
        with self._lock:
            entry = self._entries.get(key)
            if entry is None:
                entry = _IdempotencyEntry()
                self._entries[key] = entry
                is_owner = True
            else:
                is_owner = False
        if is_owner:
            return None, True
        entry.done.wait()
        return entry.result, False

    def complete(self, key: str, result: PaymentResult) -> None:
        with self._lock:
            entry = self._entries[key]
        entry.result = result
        entry.done.set()


class PaymentGateway:
    def __init__(self, processor: PaymentProcessor, retry_policy: RetryPolicy):
        self._processor = processor
        self._retry_policy = retry_policy
        self._store = IdempotencyStore()
        self._ledger = Ledger()
        self._seq = 0
        self._seq_lock = threading.Lock()

    @property
    def ledger(self) -> Ledger:
        return self._ledger

    def charge(self, request: PaymentRequest) -> PaymentResult:
        """Same idempotency_key always returns the same terminal result,
        SUCCESS or FAILED, without reprocessing (some real gateways allow
        retrying a failed key; this exercise keeps the simpler "terminal
        result is final" semantics).
        """
        cached, is_owner = self._store.reserve_or_wait(request.idempotency_key)
        if not is_owner:
            return cached

        result = self._process_with_retry(request)
        self._store.complete(request.idempotency_key, result)
        return result

    def _process_with_retry(self, request: PaymentRequest) -> PaymentResult:
        with self._seq_lock:
            self._seq += 1
            payment_id = f"PAY-{self._seq}"

        result = PaymentResult(id=payment_id, request=request)

        for attempt in range(1, self._retry_policy.max_attempts + 1):
            error: Optional[str] = None
            try:
                self._processor.process(request)
            except PaymentProcessingError as e:
                error = str(e)
            result.attempts.append(Attempt(number=attempt, success=error is None, error=error, at=datetime.now()))

            if error is None:
                result.status = PaymentStatus.SUCCESS
                self._ledger.record(LedgerEntry(
                    payment_id=payment_id,
                    payer_id=request.payer_id,
                    payee_id=request.payee_id,
                    amount=request.amount,
                    at=datetime.now(),
                ))
                return result

            if attempt < self._retry_policy.max_attempts:
                time.sleep(self._retry_policy.delay_for(attempt))

        result.status = PaymentStatus.FAILED
        return result


def _demo() -> None:
    gateway = PaymentGateway(FakePaymentProcessor(), RetryPolicy(max_attempts=3, delay_seconds=0.001))

    request = PaymentRequest("order-123", "payer-1", "payee-1", 250.0, "INR")
    result = gateway.charge(request)
    print(f"Charge result: {result.status.name} (id={result.id})")

    cached = gateway.charge(request)
    print(f"Repeat with same key returns cached id={cached.id}")

    print(f"Ledger entries: {len(gateway.ledger.entries())}")


if __name__ == "__main__":
    _demo()
