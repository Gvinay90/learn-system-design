"""Snake and Ladder LLD — Python reference implementation.

A configurable board with snakes and ladders, an injectable dice (strategy
pattern) for deterministic testing, and a turn-based game engine that
enforces the exact-landing-to-win rule.
See ../go/snakeandladder.go for the original design writeup.
"""
from __future__ import annotations

import random
from dataclasses import dataclass
from typing import Dict, List, Optional, Protocol, Sequence


class Dice(Protocol):
    """Rolls a value each turn. Implementations can be random or scripted."""

    def roll(self) -> int:
        ...


class StandardDice:
    """A fair single six-sided die, seedable for reproducible runs."""

    def __init__(self, seed: Optional[int] = None):
        self._rng = random.Random(seed)

    def roll(self) -> int:
        return self._rng.randint(1, 6)


class ScriptedDice:
    """Replays a fixed sequence of rolls, cycling once exhausted.

    This is the key to writing deterministic, non-flaky tests for the
    engine.
    """

    def __init__(self, *rolls: int):
        self.rolls = list(rolls)
        self._pos = 0

    def roll(self) -> int:
        if not self.rolls:
            return 1
        v = self.rolls[self._pos % len(self.rolls)]
        self._pos += 1
        return v


@dataclass(frozen=True)
class Entity:
    """A board hazard/shortcut: a snake (head->tail) or ladder (bottom->top)."""

    start: int
    end: int


class BoardConfigError(Exception):
    pass


class Board:
    """Holds the size and the snake/ladder map keyed by the cell a player
    lands on.
    """

    def __init__(self, size: int, snakes: Sequence[Entity] = (), ladders: Sequence[Entity] = ()):
        self.size = size
        self._entities: Dict[int, int] = {}

        for s in snakes:
            if s.start <= s.end:
                raise BoardConfigError(f"snake start {s.start} must be greater than end {s.end}")
            self._add_entity(s.start, s.end)
        for l in ladders:
            if l.start >= l.end:
                raise BoardConfigError(f"ladder start {l.start} must be less than end {l.end}")
            self._add_entity(l.start, l.end)

    def _add_entity(self, start: int, end: int) -> None:
        if start <= 1 or start >= self.size:
            raise BoardConfigError(f"entity start {start} must be within (1, {self.size})")
        if start in self._entities:
            raise BoardConfigError(f"cell {start} already has a snake or ladder")
        self._entities[start] = end

    def resolve(self, cell: int) -> int:
        """Applies any snake/ladder at the landed-on cell, returning the
        final resting cell.
        """
        return self._entities.get(cell, cell)


class Player:
    def __init__(self, name: str):
        self.name = name
        self.position = 0


class NotEnoughPlayersError(Exception):
    pass


class GameAlreadyOverError(Exception):
    pass


class TooManyTurnsError(Exception):
    pass


class Game:
    """Drives turn order, dice rolls, and win detection for a single match."""

    def __init__(self, board: Board, dice: Dice, player_names: Sequence[str]):
        if len(player_names) < 2:
            raise NotEnoughPlayersError("need at least two players")
        self.board = board
        self.dice = dice
        self.players: List[Player] = [Player(name) for name in player_names]
        self._turn = 0
        self.winner: Optional[Player] = None

    def play_turn(self) -> Player:
        """Rolls the dice for the current player, applies the
        exact-landing rule and any snake/ladder at the destination, then
        rotates turn order to the next player. Returns the player who just
        moved.
        """
        if self.winner is not None:
            raise GameAlreadyOverError("game already has a winner")

        player = self.players[self._turn]
        roll = self.dice.roll()
        target = player.position + roll

        # Overshooting the last cell is not a legal move: the player stays put.
        if target <= self.board.size:
            player.position = self.board.resolve(target)

        if player.position == self.board.size:
            self.winner = player
        else:
            self._turn = (self._turn + 1) % len(self.players)
        return player

    @property
    def current_player(self) -> Player:
        """Returns whose turn it is next."""
        return self.players[self._turn]

    def play(self, max_turns: int) -> Player:
        """Runs turns until a winner emerges, guarding against runaway loops."""
        for _ in range(max_turns):
            self.play_turn()
            if self.winner is not None:
                return self.winner
        raise TooManyTurnsError(f"no winner after {max_turns} turns")


def _demo() -> None:
    board = Board(
        100,
        snakes=[Entity(99, 41), Entity(70, 34), Entity(52, 29)],
        ladders=[Entity(4, 25), Entity(21, 82), Entity(50, 95)],
    )
    game = Game(board, StandardDice(seed=42), ["Alice", "Bob"])
    winner = game.play(1000)
    print(f"{winner.name} wins at position {winner.position}!")


if __name__ == "__main__":
    _demo()
