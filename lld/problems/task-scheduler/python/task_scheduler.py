"""Task Scheduler LLD — Python reference implementation.

In-memory, priority-based job scheduler with delayed execution and retries,
running on a fixed worker-thread pool (single process — see the HLD track
for a distributed cron/scheduler design). Mirrors go/taskscheduler.go.
"""
from __future__ import annotations

import heapq
import itertools
import threading
import time
from dataclasses import dataclass, field
from enum import Enum, auto
from typing import Callable, Dict, Optional


class Status(Enum):
    PENDING = auto()
    RUNNING = auto()
    RETRYING = auto()
    SUCCEEDED = auto()
    FAILED = auto()


class RetryPolicy:
    """Computes the delay (in seconds) before the (attempt+1)-th retry."""

    def next_delay(self, attempt: int) -> float:
        raise NotImplementedError


@dataclass(frozen=True)
class ExponentialBackoff(RetryPolicy):
    base: float
    max: float

    def next_delay(self, attempt: int) -> float:
        d = self.base
        for _ in range(attempt):
            d *= 2
            if d > self.max:
                return self.max
        return d


_DEFAULT_RETRY_POLICY = ExponentialBackoff(base=0.01, max=0.2)

# Job.task is a zero-arg callable that raises on failure, returns normally
# on success — mirrors Go's `func() error`.
Task = Callable[[], None]


@dataclass
class Job:
    id: str
    priority: int
    run_at: float  # time.monotonic()-based deadline
    task: Task
    max_retries: int = 0
    retry: Optional[RetryPolicy] = None
    attempts: int = field(default=0, init=False)


@dataclass
class JobResult:
    job_id: str
    status: Status
    attempts: int = 0
    error: Optional[BaseException] = None


@dataclass(order=True)
class _QueueEntry:
    """Min-heap entry ordered by run_at (earliest due first), tie-broken by
    higher priority first. `seq` is a tiebreaker so heapq never has to
    compare Job objects (which aren't orderable)."""

    sort_key: tuple = field(compare=True)
    seq: int = field(compare=True)
    job: Job = field(compare=False)


class Scheduler:
    """Runs submitted jobs on a fixed worker-thread pool, respecting
    priority, delayed run_at times, and per-job retry policies."""

    def __init__(self, workers: int):
        self._workers = workers
        self._lock = threading.Lock()
        self._heap: list = []
        self._seq = itertools.count()
        self._results: Dict[str, JobResult] = {}
        self._threads: list = []
        self._stopped = True

    def submit(self, job: Job) -> None:
        """Enqueues a job. Safe to call before or after start(), and
        concurrently from multiple threads."""
        if job.retry is None:
            job.retry = _DEFAULT_RETRY_POLICY
        with self._lock:
            # Negate priority so the heap (min-first) surfaces the highest
            # priority first among equal run_at times.
            heapq.heappush(self._heap, _QueueEntry((job.run_at, -job.priority), next(self._seq), job))
            self._results[job.id] = JobResult(job_id=job.id, status=Status.PENDING)

    def start(self) -> None:
        self._stopped = False
        self._threads = []
        for i in range(self._workers):
            t = threading.Thread(target=self._worker_loop, name=f"scheduler-worker-{i}", daemon=True)
            self._threads.append(t)
            t.start()

    def stop(self) -> None:
        self._stopped = True
        for t in self._threads:
            t.join()

    def _worker_loop(self) -> None:
        while not self._stopped:
            job = self._pop_due()
            if job is not None:
                self._execute(job)
            else:
                time.sleep(0.002)

    def _pop_due(self) -> Optional[Job]:
        """Removes and returns the highest-priority job whose run_at has
        already elapsed, or None if none is ready."""
        with self._lock:
            if not self._heap:
                return None
            if self._heap[0].job.run_at > time.monotonic():
                return None
            return heapq.heappop(self._heap).job

    def _execute(self, job: Job) -> None:
        with self._lock:
            self._results[job.id].status = Status.RUNNING

        err: Optional[BaseException] = None
        try:
            job.task()
        except BaseException as e:  # noqa: BLE001 - mirror Go's generic error
            err = e
        job.attempts += 1

        with self._lock:
            res = self._results[job.id]
            res.attempts = job.attempts
            if err is None:
                res.status = Status.SUCCEEDED
                res.error = None
                return
            res.error = err
            if job.attempts <= job.max_retries:
                res.status = Status.RETRYING
                delay = job.retry.next_delay(job.attempts - 1)
                job.run_at = time.monotonic() + delay
                heapq.heappush(self._heap, _QueueEntry((job.run_at, -job.priority), next(self._seq), job))
                return
            res.status = Status.FAILED

    def get_result(self, job_id: str) -> Optional[JobResult]:
        with self._lock:
            res = self._results.get(job_id)
            if res is None:
                return None
            return JobResult(job_id=res.job_id, status=res.status, attempts=res.attempts, error=res.error)

    def wait_for_result(self, job_id: str, timeout: float) -> Optional[JobResult]:
        """Polls until job_id reaches a terminal state (SUCCEEDED or FAILED)
        or the timeout (in seconds) elapses."""
        deadline = time.monotonic() + timeout
        while True:
            res = self.get_result(job_id)
            if res is not None and res.status in (Status.SUCCEEDED, Status.FAILED):
                return res
            if time.monotonic() > deadline:
                return None
            time.sleep(0.002)


if __name__ == "__main__":
    scheduler = Scheduler(workers=2)
    scheduler.start()

    ran = {"count": 0}

    def demo_task() -> None:
        ran["count"] += 1
        print("demo job executed")

    scheduler.submit(Job(id="demo", priority=1, run_at=time.monotonic(), task=demo_task))

    result = scheduler.wait_for_result("demo", timeout=1.0)
    print(f"Result: {result}")

    scheduler.stop()
