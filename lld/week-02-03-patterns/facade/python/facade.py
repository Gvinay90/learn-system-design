"""Facade pattern — HomeTheaterFacade coordinates five subsystems behind two calls.

See ../README.md for the design writeup.
"""
from __future__ import annotations


class Amplifier:
    def __init__(self) -> None:
        self.on = False
        self.volume = 0

    def turn_on(self) -> None:
        self.on = True

    def turn_off(self) -> None:
        self.on = False
        self.volume = 0

    def set_volume(self, volume: int) -> None:
        self.volume = volume


class DvdPlayer:
    def __init__(self) -> None:
        self.on = False
        self.movie = ""

    def turn_on(self) -> None:
        self.on = True

    def turn_off(self) -> None:
        self.on = False
        self.movie = ""

    def play(self, movie: str) -> None:
        self.movie = movie

    def stop(self) -> None:
        self.movie = ""


class Projector:
    def __init__(self) -> None:
        self.on = False
        self.input = ""

    def turn_on(self) -> None:
        self.on = True

    def turn_off(self) -> None:
        self.on = False
        self.input = ""

    def set_input(self, source: str) -> None:
        self.input = source


class Screen:
    def __init__(self) -> None:
        self.down = False

    def lower(self) -> None:
        self.down = True

    def raise_(self) -> None:
        self.down = False


class Lights:
    def __init__(self) -> None:
        self.dimmed = False

    def dim(self) -> None:
        self.dimmed = True

    def brighten(self) -> None:
        self.dimmed = False


class HomeTheaterFacade:
    """Hides the coordination order of five subsystems behind two calls."""

    def __init__(self) -> None:
        self.amp = Amplifier()
        self.dvd = DvdPlayer()
        self.projector = Projector()
        self.screen = Screen()
        self.lights = Lights()

    def watch_movie(self, movie: str) -> None:
        self.lights.dim()
        self.screen.lower()
        self.projector.turn_on()
        self.projector.set_input("dvd")
        self.amp.turn_on()
        self.amp.set_volume(7)
        self.dvd.turn_on()
        self.dvd.play(movie)

    def end_movie(self) -> None:
        self.dvd.stop()
        self.dvd.turn_off()
        self.amp.turn_off()
        self.projector.turn_off()
        self.screen.raise_()
        self.lights.brighten()


def _demo() -> None:
    theater = HomeTheaterFacade()
    theater.watch_movie("Inception")
    print(f"Now playing: {theater.dvd.movie}")
    theater.end_movie()
    print(f"Movie ended, lights dimmed: {theater.lights.dimmed}")


if __name__ == "__main__":
    _demo()
