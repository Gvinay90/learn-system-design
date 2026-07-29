import threading
import time

import pytest

from task_scheduler import (
    ExponentialBackoff,
    Job,
    Scheduler,
    Status,
)


def test_jobs_run_and_succeed():
    s = Scheduler(workers=2)
    s.start()
    try:
        ran = {"count": 0}

        def task():
            ran["count"] += 1

        s.submit(Job(id="J1", priority=1, run_at=time.monotonic(), task=task))

        res = s.wait_for_result("J1", timeout=1.0)
        assert res is not None, "expected job to complete within timeout"
        assert res.status is Status.SUCCEEDED
        assert res.attempts == 1
        assert ran["count"] == 1
    finally:
        s.stop()


# Among jobs due at the same time, the higher-priority job should run first.
# A single worker keeps execution order deterministic, and both jobs share a
# run_at in the near future so they're both "due" as soon as the scheduler
# starts.
def test_priority_ordering():
    s = Scheduler(workers=1)

    order = []
    mu = threading.Lock()
    due = time.monotonic() + 0.02

    def make_task(name):
        def task():
            with mu:
                order.append(name)
        return task

    s.submit(Job(id="low", priority=1, run_at=due, task=make_task("low")))
    s.submit(Job(id="high", priority=10, run_at=due, task=make_task("high")))

    s.start()
    try:
        assert s.wait_for_result("low", timeout=1.0) is not None
        assert s.wait_for_result("high", timeout=1.0) is not None
    finally:
        s.stop()

    with mu:
        assert order == ["high", "low"], f"expected high before low, got {order}"


def test_delayed_job_does_not_run_early():
    s = Scheduler(workers=1)
    s.start()
    try:
        ran_at = {}
        submitted_at = time.monotonic()
        delay = 0.06

        def task():
            ran_at["t"] = time.monotonic()

        s.submit(Job(id="delayed", priority=1, run_at=submitted_at + delay, task=task))

        res = s.wait_for_result("delayed", timeout=1.0)
        assert res is not None and res.status is Status.SUCCEEDED
        assert ran_at["t"] - submitted_at >= delay
    finally:
        s.stop()


def test_retries_on_failure_then_succeeds():
    s = Scheduler(workers=1)
    s.start()
    try:
        attempts = {"n": 0}

        def task():
            attempts["n"] += 1
            if attempts["n"] < 3:
                raise RuntimeError("transient failure")

        s.submit(Job(
            id="flaky",
            priority=1,
            run_at=time.monotonic(),
            task=task,
            max_retries=3,
            retry=ExponentialBackoff(base=0.002, max=0.01),
        ))

        res = s.wait_for_result("flaky", timeout=1.0)
        assert res is not None, "expected flaky job to eventually reach a terminal state"
        assert res.status is Status.SUCCEEDED, f"got {res.status} (err={res.error})"
        assert res.attempts == 3
    finally:
        s.stop()


def test_exhausts_retries_and_fails():
    s = Scheduler(workers=1)
    s.start()
    try:
        attempts = {"n": 0}

        def task():
            attempts["n"] += 1
            raise RuntimeError("permanent failure")

        s.submit(Job(
            id="doomed",
            priority=1,
            run_at=time.monotonic(),
            task=task,
            max_retries=2,
            retry=ExponentialBackoff(base=0.002, max=0.01),
        ))

        res = s.wait_for_result("doomed", timeout=1.0)
        assert res is not None, "expected doomed job to eventually reach a terminal state"
        assert res.status is Status.FAILED
        # max_retries=2 means 1 initial attempt + 2 retries = 3 attempts total.
        assert res.attempts == 3
        assert res.error is not None and str(res.error) == "permanent failure"
    finally:
        s.stop()


def test_unknown_job_result_not_found():
    s = Scheduler(workers=1)
    assert s.get_result("nonexistent") is None


def test_exponential_backoff_caps_at_max():
    b = ExponentialBackoff(base=0.01, max=0.03)
    assert b.next_delay(0) == pytest.approx(0.01)
    assert b.next_delay(1) == pytest.approx(0.02)
    assert b.next_delay(2) == pytest.approx(0.03)
    assert b.next_delay(5) == pytest.approx(0.03)


# Asserts many jobs submitted concurrently from multiple threads all get
# executed exactly once by the worker pool.
def test_concurrent_submit_and_execute():
    s = Scheduler(workers=4)
    s.start()
    try:
        n = 50

        def job_id(i):
            return f"concurrent-{chr(ord('A') + i % 26)}{i // 26}"

        def submit_one(i):
            s.submit(Job(id=job_id(i), priority=i % 5, run_at=time.monotonic(), task=lambda: None))

        threads = [threading.Thread(target=submit_one, args=(i,)) for i in range(n)]
        for t in threads:
            t.start()
        for t in threads:
            t.join()

        for i in range(n):
            res = s.wait_for_result(job_id(i), timeout=2.0)
            assert res is not None and res.status is Status.SUCCEEDED, f"job {job_id(i)} did not succeed: {res}"
    finally:
        s.stop()


if __name__ == "__main__":
    import sys

    sys.exit(pytest.main([__file__, "-v"]))
