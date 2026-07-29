import pytest

from chess import (
    Board,
    Color,
    Game,
    IllegalMoveError,
    PieceType,
    new_piece,
    parse_square,
)


def sq(s):
    return parse_square(s)


def test_pawn_opening_move_happy_path():
    g = Game()
    frm, to = sq("e2"), sq("e4")

    g.move(frm, to)
    assert g.board.at(to) is not None and g.board.at(to).get_type() is PieceType.PAWN
    assert g.board.at(frm) is None
    assert g.turn is Color.BLACK


def test_illegal_move_rejected():
    g = Game()

    with pytest.raises(IllegalMoveError):
        g.move(sq("e2"), sq("e5"))  # pawn cannot jump 3 squares

    with pytest.raises(IllegalMoveError):
        g.move(sq("a1"), sq("a3"))  # rook blocked by its own pawn


def test_check_detection():
    board = Board()
    board.set(sq("e1"), new_piece(PieceType.KING, Color.WHITE))
    board.set(sq("e8"), new_piece(PieceType.ROOK, Color.BLACK))
    board.set(sq("a1"), new_piece(PieceType.KING, Color.BLACK))
    g = Game(board=board, turn=Color.WHITE)

    assert g.is_in_check(Color.WHITE)
    assert not g.is_checkmate(Color.WHITE)  # king can step aside


def test_fools_mate_checkmate():
    g = Game()
    moves = [("f2", "f3"), ("e7", "e5"), ("g2", "g4"), ("d8", "h4")]
    for src, dst in moves:
        g.move(sq(src), sq(dst))

    assert g.is_in_check(Color.WHITE)
    assert g.is_checkmate(Color.WHITE)
