"""Chess LLD — Python reference implementation.

8x8 board, per-piece move generation (Strategy pattern via the Piece
protocol), turn alternation, and check/checkmate detection. En passant,
castling, and pawn promotion are intentionally out of scope (see
../README.md for talking points).
"""
from __future__ import annotations

from dataclasses import dataclass, field
from enum import Enum, auto
from typing import List, Optional, Protocol


class Color(Enum):
    WHITE = auto()
    BLACK = auto()

    def opposite(self) -> "Color":
        return Color.BLACK if self is Color.WHITE else Color.WHITE

    def __str__(self) -> str:
        return "White" if self is Color.WHITE else "Black"


class PieceType(Enum):
    PAWN = auto()
    KNIGHT = auto()
    BISHOP = auto()
    ROOK = auto()
    QUEEN = auto()
    KING = auto()


@dataclass(frozen=True)
class Position:
    row: int
    col: int

    def in_bounds(self) -> bool:
        return 0 <= self.row < 8 and 0 <= self.col < 8


def parse_square(s: str) -> Position:
    """Converts algebraic notation ("e2") into a Position."""
    if len(s) != 2:
        raise ValueError(f"invalid square {s!r}")
    col = ord(s[0]) - ord("a")
    row = 8 - int(s[1])
    pos = Position(row, col)
    if not pos.in_bounds():
        raise ValueError(f"invalid square {s!r}")
    return pos


class Piece(Protocol):
    color: Color

    def get_type(self) -> PieceType: ...
    def can_move(self, board: "Board", frm: Position, to: Position) -> bool: ...


def _target_ok(board: "Board", color: Color, to: Position) -> bool:
    occupant = board.at(to)
    return occupant is None or occupant.color != color


class _BasePiece:
    def __init__(self, color: Color):
        self.color = color


class PawnPiece(_BasePiece):
    def get_type(self) -> PieceType:
        return PieceType.PAWN

    def can_move(self, board: "Board", frm: Position, to: Position) -> bool:
        direction = -1 if self.color is Color.WHITE else 1
        start_row = 6 if self.color is Color.WHITE else 1
        dr, dc = to.row - frm.row, to.col - frm.col

        if dc == 0:
            if dr == direction:
                return board.at(to) is None
            if dr == 2 * direction and frm.row == start_row:
                mid = Position(frm.row + direction, frm.col)
                return board.at(mid) is None and board.at(to) is None
            return False
        if abs(dc) == 1 and dr == direction:
            occupant = board.at(to)
            return occupant is not None and occupant.color != self.color
        return False


class KnightPiece(_BasePiece):
    def get_type(self) -> PieceType:
        return PieceType.KNIGHT

    def can_move(self, board: "Board", frm: Position, to: Position) -> bool:
        dr, dc = abs(to.row - frm.row), abs(to.col - frm.col)
        if not ((dr == 1 and dc == 2) or (dr == 2 and dc == 1)):
            return False
        return _target_ok(board, self.color, to)


class BishopPiece(_BasePiece):
    def get_type(self) -> PieceType:
        return PieceType.BISHOP

    def can_move(self, board: "Board", frm: Position, to: Position) -> bool:
        dr, dc = to.row - frm.row, to.col - frm.col
        if abs(dr) != abs(dc) or dr == 0:
            return False
        return board.clear_path(frm, to) and _target_ok(board, self.color, to)


class RookPiece(_BasePiece):
    def get_type(self) -> PieceType:
        return PieceType.ROOK

    def can_move(self, board: "Board", frm: Position, to: Position) -> bool:
        dr, dc = to.row - frm.row, to.col - frm.col
        if (dr == 0) == (dc == 0):
            return False
        return board.clear_path(frm, to) and _target_ok(board, self.color, to)


class QueenPiece(_BasePiece):
    def get_type(self) -> PieceType:
        return PieceType.QUEEN

    def can_move(self, board: "Board", frm: Position, to: Position) -> bool:
        dr, dc = to.row - frm.row, to.col - frm.col
        straight = (dr == 0) != (dc == 0)
        diagonal = dr != 0 and abs(dr) == abs(dc)
        if not straight and not diagonal:
            return False
        return board.clear_path(frm, to) and _target_ok(board, self.color, to)


class KingPiece(_BasePiece):
    def get_type(self) -> PieceType:
        return PieceType.KING

    def can_move(self, board: "Board", frm: Position, to: Position) -> bool:
        dr, dc = abs(to.row - frm.row), abs(to.col - frm.col)
        if dr > 1 or dc > 1 or (dr == 0 and dc == 0):
            return False
        return _target_ok(board, self.color, to)


_PIECE_CLASSES = {
    PieceType.PAWN: PawnPiece,
    PieceType.KNIGHT: KnightPiece,
    PieceType.BISHOP: BishopPiece,
    PieceType.ROOK: RookPiece,
    PieceType.QUEEN: QueenPiece,
    PieceType.KING: KingPiece,
}


def new_piece(piece_type: PieceType, color: Color) -> Piece:
    return _PIECE_CLASSES[piece_type](color)


def _sign(n: int) -> int:
    return (n > 0) - (n < 0)


class Board:
    """Holds the 8x8 grid of squares."""

    def __init__(self):
        self.squares: List[List[Optional[Piece]]] = [[None] * 8 for _ in range(8)]

    @staticmethod
    def new_game() -> "Board":
        b = Board()
        back_rank = [
            PieceType.ROOK, PieceType.KNIGHT, PieceType.BISHOP, PieceType.QUEEN,
            PieceType.KING, PieceType.BISHOP, PieceType.KNIGHT, PieceType.ROOK,
        ]
        for col, t in enumerate(back_rank):
            b.squares[0][col] = new_piece(t, Color.BLACK)
            b.squares[7][col] = new_piece(t, Color.WHITE)
            b.squares[1][col] = new_piece(PieceType.PAWN, Color.BLACK)
            b.squares[6][col] = new_piece(PieceType.PAWN, Color.WHITE)
        return b

    def at(self, pos: Position) -> Optional[Piece]:
        return self.squares[pos.row][pos.col]

    def set(self, pos: Position, piece: Optional[Piece]) -> None:
        self.squares[pos.row][pos.col] = piece

    def move(self, frm: Position, to: Position) -> Optional[Piece]:
        """Relocates a piece with no legality checks; callers must validate first."""
        captured = self.at(to)
        self.set(to, self.at(frm))
        self.set(frm, None)
        return captured

    def clone(self) -> "Board":
        nb = Board()
        nb.squares = [row[:] for row in self.squares]
        return nb

    def clear_path(self, frm: Position, to: Position) -> bool:
        step_r, step_c = _sign(to.row - frm.row), _sign(to.col - frm.col)
        r, c = frm.row + step_r, frm.col + step_c
        while r != to.row or c != to.col:
            if self.squares[r][c] is not None:
                return False
            r += step_r
            c += step_c
        return True

    def find_king(self, color: Color) -> Optional[Position]:
        for r in range(8):
            for c in range(8):
                p = self.squares[r][c]
                if p is not None and p.get_type() is PieceType.KING and p.color is color:
                    return Position(r, c)
        return None

    def is_square_attacked(self, pos: Position, attacker: Color) -> bool:
        for r in range(8):
            for c in range(8):
                p = self.squares[r][c]
                if p is not None and p.color is attacker and p.can_move(self, Position(r, c), pos):
                    return True
        return False

    def is_in_check(self, color: Color) -> bool:
        king_pos = self.find_king(color)
        if king_pos is None:
            return False
        return self.is_square_attacked(king_pos, color.opposite())


@dataclass
class Move:
    frm: Position
    to: Position
    piece: Piece
    captured: Optional[Piece]


class NotYourTurnError(Exception):
    pass


class NoPieceAtSourceError(Exception):
    pass


class OwnPieceAtTargetError(Exception):
    pass


class IllegalMoveError(Exception):
    pass


class MoveLeavesKingInCheckError(Exception):
    pass


class Game:
    """Orchestrates turn alternation and rule enforcement on top of Board."""

    def __init__(self, board: Optional[Board] = None, turn: Color = Color.WHITE):
        self.board = board if board is not None else Board.new_game()
        self.turn = turn
        self.history: List[Move] = []

    def move(self, frm: Position, to: Position) -> None:
        piece = self.board.at(frm)
        if piece is None:
            raise NoPieceAtSourceError()
        if piece.color != self.turn:
            raise NotYourTurnError()
        target = self.board.at(to)
        if target is not None and target.color == piece.color:
            raise OwnPieceAtTargetError()
        if not piece.can_move(self.board, frm, to):
            raise IllegalMoveError()

        trial = self.board.clone()
        trial.move(frm, to)
        if trial.is_in_check(piece.color):
            raise MoveLeavesKingInCheckError()

        captured = self.board.move(frm, to)
        self.history.append(Move(frm, to, piece, captured))
        self.turn = self.turn.opposite()

    def is_in_check(self, color: Color) -> bool:
        return self.board.is_in_check(color)

    def is_checkmate(self, color: Color) -> bool:
        """Brute-forces every (piece, destination) pair; simple at interview scope."""
        if not self.board.is_in_check(color):
            return False
        for r in range(8):
            for c in range(8):
                piece = self.board.squares[r][c]
                if piece is None or piece.color != color:
                    continue
                frm = Position(r, c)
                for tr in range(8):
                    for tc in range(8):
                        to = Position(tr, tc)
                        if not piece.can_move(self.board, frm, to):
                            continue
                        trial = self.board.clone()
                        trial.move(frm, to)
                        if not trial.is_in_check(color):
                            return False
        return True


if __name__ == "__main__":
    game = Game()
    game.move(parse_square("e2"), parse_square("e4"))
    print(f"White plays e2-e4; turn is now {game.turn}")

    game.move(parse_square("e7"), parse_square("e5"))
    print(f"Black plays e7-e5; turn is now {game.turn}")

    mate = Game()
    for src, dst in [("f2", "f3"), ("e7", "e5"), ("g2", "g4"), ("d8", "h4")]:
        mate.move(parse_square(src), parse_square(dst))
    print(
        f"After Fool's Mate: White in check = {mate.is_in_check(Color.WHITE)}, "
        f"checkmate = {mate.is_checkmate(Color.WHITE)}"
    )
