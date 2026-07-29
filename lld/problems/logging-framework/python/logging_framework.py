"""Logging Framework LLD — Python reference implementation.

A Logger holds a minimum severity threshold and a list of Appenders
(Strategy pattern) that it dispatches formatted records to. Appenders are
interchangeable output destinations (console, file, ...); new ones can be
added without touching Logger itself.
"""
from __future__ import annotations

import sys
import threading
from abc import ABC, abstractmethod
from dataclasses import dataclass
from datetime import datetime, timezone
from enum import IntEnum
from typing import List, TextIO


class Level(IntEnum):
    """Log severity, ordered DEBUG < INFO < WARN < ERROR."""

    DEBUG = 0
    INFO = 1
    WARN = 2
    ERROR = 3


@dataclass(frozen=True)
class Record:
    """A single formatted log entry handed to every Appender."""

    timestamp: datetime
    level: Level
    message: str

    def format(self) -> str:
        return f"{self.timestamp.isoformat()} [{self.level.name}] {self.message}"


class Appender(ABC):
    """A pluggable output destination for log records."""

    @abstractmethod
    def append(self, record: Record) -> None:
        raise NotImplementedError


class ConsoleAppender(Appender):
    """Writes formatted records to a text stream (stdout by default)."""

    def __init__(self, stream: TextIO = sys.stdout):
        self._stream = stream

    def append(self, record: Record) -> None:
        print(record.format(), file=self._stream)


class FileAppender(Appender):
    """Appends formatted records, one per line, to a file at a
    caller-supplied path. The file is opened once (append mode) and kept
    open for the lifetime of the appender.
    """

    def __init__(self, path: str):
        self._lock = threading.Lock()
        self._file = open(path, "a", encoding="utf-8")

    def append(self, record: Record) -> None:
        with self._lock:
            self._file.write(record.format() + "\n")
            self._file.flush()

    def close(self) -> None:
        with self._lock:
            self._file.close()


class Logger:
    """Dispatches records to its appenders when the record's level meets or
    exceeds the logger's minimum threshold. Thread-safe via an internal
    lock around dispatch.
    """

    def __init__(self, level: Level, appenders: List[Appender] | None = None):
        self._lock = threading.Lock()
        self._level = level
        self._appenders: List[Appender] = list(appenders) if appenders else []

    def add_appender(self, appender: Appender) -> None:
        with self._lock:
            self._appenders.append(appender)

    def set_level(self, level: Level) -> None:
        with self._lock:
            self._level = level

    def log(self, level: Level, message: str) -> None:
        with self._lock:
            if level < self._level:
                return
            record = Record(timestamp=datetime.now(timezone.utc), level=level, message=message)
            for appender in self._appenders:
                appender.append(record)

    def debug(self, message: str) -> None:
        self.log(Level.DEBUG, message)

    def info(self, message: str) -> None:
        self.log(Level.INFO, message)

    def warn(self, message: str) -> None:
        self.log(Level.WARN, message)

    def error(self, message: str) -> None:
        self.log(Level.ERROR, message)


def _demo() -> None:
    logger = Logger(Level.INFO, [ConsoleAppender()])
    logger.debug("this is filtered out below INFO threshold")
    logger.info("service started")
    logger.warn("cache miss rate elevated")
    logger.error("failed to connect to downstream")


if __name__ == "__main__":
    _demo()
