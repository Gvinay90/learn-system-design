"""Observer pattern — WeatherStation subject notifies subscribed observers.

See ../README.md for the design writeup.
"""
from __future__ import annotations

from typing import List, Protocol


class Observer(Protocol):
    def update(self, temp_c: float, humidity: float) -> None: ...


class WeatherStation:
    """The subject: holds the current readings and notifies all subscribed
    observers whenever they change."""

    def __init__(self) -> None:
        self._observers: List[Observer] = []
        self.temp_c = 0.0
        self.humidity = 0.0

    def subscribe(self, observer: Observer) -> None:
        self._observers.append(observer)

    def unsubscribe(self, observer: Observer) -> None:
        if observer in self._observers:
            self._observers.remove(observer)

    def set_measurements(self, temp_c: float, humidity: float) -> None:
        self.temp_c = temp_c
        self.humidity = humidity
        for observer in self._observers:
            observer.update(temp_c, humidity)


class PhoneDisplay:
    def __init__(self) -> None:
        self.last_temp_c = 0.0
        self.last_humidity = 0.0
        self.update_count = 0

    def update(self, temp_c: float, humidity: float) -> None:
        self.last_temp_c = temp_c
        self.last_humidity = humidity
        self.update_count += 1


class LoggingDisplay:
    def __init__(self) -> None:
        self.history: List[str] = []

    def update(self, temp_c: float, humidity: float) -> None:
        self.history.append(f"temp={temp_c:.1f} humidity={humidity:.1f}")


def _demo() -> None:
    station = WeatherStation()
    phone = PhoneDisplay()
    logger = LoggingDisplay()
    station.subscribe(phone)
    station.subscribe(logger)

    station.set_measurements(21.5, 60)
    print(f"Phone display: {phone.last_temp_c}C, {phone.last_humidity}%")


if __name__ == "__main__":
    _demo()
