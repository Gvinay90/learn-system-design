import threading

import pytest

from logging_framework import Appender, ConsoleAppender, FileAppender, Level, Logger


class MockAppender(Appender):
    """Test double that records every Record it receives. Guarded by a lock
    so it is safe to share across threads in the concurrency test.
    """

    def __init__(self):
        self._lock = threading.Lock()
        self.records = []

    def append(self, record) -> None:
        with self._lock:
            self.records.append(record)

    def count(self) -> int:
        with self._lock:
            return len(self.records)


@pytest.mark.parametrize(
    "threshold,log_at,want_count",
    [
        (Level.DEBUG, Level.DEBUG, 1),
        (Level.INFO, Level.DEBUG, 0),
        (Level.INFO, Level.INFO, 1),
        (Level.WARN, Level.INFO, 0),
        (Level.WARN, Level.WARN, 1),
        (Level.ERROR, Level.WARN, 0),
        (Level.ERROR, Level.ERROR, 1),
    ],
)
def test_level_filtering(threshold, log_at, want_count):
    appender = MockAppender()
    logger = Logger(threshold, [appender])
    logger.log(log_at, "message")
    assert appender.count() == want_count


def test_multiple_appenders_all_receive_record():
    a1, a2, a3 = MockAppender(), MockAppender(), MockAppender()
    logger = Logger(Level.INFO, [a1, a2, a3])

    logger.info("hello")
    logger.debug("should be filtered")

    assert a1.count() == 1
    assert a2.count() == 1
    assert a3.count() == 1


def test_add_appender_after_construction():
    a1 = MockAppender()
    logger = Logger(Level.DEBUG, [a1])

    a2 = MockAppender()
    logger.add_appender(a2)

    logger.info("hi")
    assert a1.count() == 1
    assert a2.count() == 1


def test_set_level_changes_threshold():
    appender = MockAppender()
    logger = Logger(Level.ERROR, [appender])

    logger.warn("filtered")
    assert appender.count() == 0

    logger.set_level(Level.WARN)
    logger.warn("passes now")
    assert appender.count() == 1


def test_record_format_contains_level_and_message():
    appender = MockAppender()
    logger = Logger(Level.DEBUG, [appender])

    logger.error("disk on fire")

    formatted = appender.records[0].format()
    assert "[ERROR]" in formatted
    assert "disk on fire" in formatted


def test_file_appender_writes_to_temp_file(tmp_path):
    log_path = tmp_path / "app.log"
    appender = FileAppender(str(log_path))
    logger = Logger(Level.INFO, [appender])

    logger.info("first line")
    logger.warn("second line")
    logger.debug("filtered, should not appear")
    appender.close()

    lines = log_path.read_text(encoding="utf-8").splitlines()
    assert len(lines) == 2
    assert "first line" in lines[0]
    assert "second line" in lines[1]
    assert not any("filtered, should not appear" in line for line in lines)


def test_concurrent_log_calls():
    thread_count = 100
    per_thread = 20

    appender = MockAppender()
    logger = Logger(Level.DEBUG, [appender])

    def worker():
        for _ in range(per_thread):
            logger.info("concurrent message")

    threads = [threading.Thread(target=worker) for _ in range(thread_count)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    assert appender.count() == thread_count * per_thread
