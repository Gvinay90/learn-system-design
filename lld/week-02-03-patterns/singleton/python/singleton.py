"""Singleton LLD — thread-safe AppConfig using a module-level lock with
double-checked locking so concurrent first-access constructs exactly one
instance. See ../README.md for the design writeup.
"""
from __future__ import annotations

import itertools
import threading
from typing import Dict, Optional


class AppConfig:
    _instance: Optional["AppConfig"] = None
    _lock = threading.Lock()
    _id_seq = itertools.count(1)

    def __init__(self) -> None:
        self._id = next(AppConfig._id_seq)
        self._settings: Dict[str, str] = {}
        self._settings_lock = threading.Lock()

    @classmethod
    def get_instance(cls) -> "AppConfig":
        if cls._instance is None:
            with cls._lock:
                if cls._instance is None:
                    cls._instance = cls()
        return cls._instance

    @property
    def id(self) -> int:
        return self._id

    def set(self, key: str, value: str) -> None:
        with self._settings_lock:
            self._settings[key] = value

    def get(self, key: str) -> Optional[str]:
        with self._settings_lock:
            return self._settings.get(key)

    @classmethod
    def _reset_for_test(cls) -> None:
        cls._instance = None


def _demo() -> None:
    config = AppConfig.get_instance()
    config.set("env", "production")
    print(f"Config instance id: {config.id}")
    print(f"env={AppConfig.get_instance().get('env')}")


if __name__ == "__main__":
    _demo()
