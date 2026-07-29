"""Tic-Tac-Toe LLD — Python reference implementation.

Generic NxN board, N alternating players, incremental win detection via
row/column/diagonal counters. See ../README.md for the design writeup.
"""
from __future__ import annotations

from collections import defaultdict
from dataclasses import dataclass, field
from enum import Enum, auto
from typing import List, Optional


class GameStatus(Enum):
    IN_PROGRESS = auto()
    WON = auto()
    DRAW = auto()


@dataclass
class Player:
    name: str
    symbol: str


class OutOfBoundsError(Exception):
    pass


class CellOccupiedError(Exception):
    pass


class GameOverError(Exception):
    pass


class Board:
    """A generic NxN grid.

    Win detection is O(1) per move: instead of rescanning the whole board,
    we keep running counts per row, column and the two diagonals for
    whichever symbol is placed. A move only ever affects one row, one
    column and (at most) two diagonals, so checking those counts after
    each move is enough to detect a win.
    """

    def __init__(self, size: int = 3):
        self.size = size
        self._cells: List[List[Optional[str]]] = [[None] * size for _ in range(size)]
        self._row_counts = defaultdict(lambda: [0] * size)
        self._col_counts = defaultdict(lambda: [0] * size)
        self._diag_count = defaultdict(int)
        self._anti_diag_count = defaultdict(int)
        self._filled_cells = 0

    def at(self, row: int, col: int) -> Optional[str]:
        return self._cells[row][col]

    def is_full(self) -> bool:
        return self._filled_cells == self.size * self.size

    def _in_bounds(self, row: int, col: int) -> bool:
        return 0 <= row < self.size and 0 <= col < self.size

    def place(self, row: int, col: int, symbol: str) -> bool:
        """Places symbol at (row, col); returns True if this move wins."""
        if not self._in_bounds(row, col):
            raise OutOfBoundsError(f"({row}, {col}) is out of bounds")
        if self._cells[row][col] is not None:
            raise CellOccupiedError(f"({row}, {col}) is already occupied")

        self._cells[row][col] = symbol
        self._filled_cells += 1

        self._row_counts[symbol][row] += 1
        self._col_counts[symbol][col] += 1
        if row == col:
            self._diag_count[symbol] += 1
        if row + col == self.size - 1:
            self._anti_diag_count[symbol] += 1

        return (
            self._row_counts[symbol][row] == self.size
            or self._col_counts[symbol][col] == self.size
            or self._diag_count[symbol] == self.size
            or self._anti_diag_count[symbol] == self.size
        )

    def __str__(self) -> str:
        rows = []
        for r in range(self.size):
            rows.append(" ".join(c if c else "." for c in self._cells[r]))
        return "\n".join(rows)


class Game:
    """Orchestrates turn order, move validation and terminal-state tracking.

    No locking here: unlike the parking lot, a single tic-tac-toe game is
    played by one turn owner at a time, so there is no concurrent-writer
    scenario to guard against.
    """

    def __init__(self, size: int, players: List[Player]):
        self.board = Board(size)
        self.players = players
        self._turn = 0
        self.status = GameStatus.IN_PROGRESS
        self.winner: Optional[Player] = None

    @property
    def current_player(self) -> Player:
        return self.players[self._turn]

    def move(self, row: int, col: int) -> bool:
        if self.status != GameStatus.IN_PROGRESS:
            raise GameOverError("game is already over")

        player = self.players[self._turn]
        won = self.board.place(row, col, player.symbol)

        if won:
            self.status = GameStatus.WON
            self.winner = player
            return True
        if self.board.is_full():
            self.status = GameStatus.DRAW
            return False

        self._turn = (self._turn + 1) % len(self.players)
        return False


def _demo() -> None:
    alice = Player("Alice", "X")
    bob = Player("Bob", "O")
    game = Game(3, [alice, bob])

    moves = [(0, 0), (1, 0), (0, 1), (1, 1), (0, 2)]
    for row, col in moves:
        won = game.move(row, col)
        if won:
            print(f"{game.winner.name} wins!")
    print(game.board)


if __name__ == "__main__":
    _demo()
