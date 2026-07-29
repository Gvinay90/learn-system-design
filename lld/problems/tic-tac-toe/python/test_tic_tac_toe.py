import pytest

from tic_tac_toe import (
    CellOccupiedError,
    Game,
    GameOverError,
    GameStatus,
    OutOfBoundsError,
    Player,
)


def new_test_players():
    return [Player("Alice", "X"), Player("Bob", "O")]


def test_row_win():
    g = Game(3, new_test_players())
    moves = [(0, 0), (1, 0), (0, 1), (1, 1), (0, 2)]
    won = False
    for row, col in moves:
        won = g.move(row, col)
    assert won
    assert g.status == GameStatus.WON
    assert g.winner.name == "Alice"


def test_column_win():
    g = Game(3, new_test_players())
    moves = [(0, 0), (0, 1), (1, 0), (1, 1), (2, 0)]
    won = False
    for row, col in moves:
        won = g.move(row, col)
    assert won
    assert g.winner.name == "Alice"


def test_diagonal_win():
    g = Game(3, new_test_players())
    moves = [(0, 0), (0, 1), (1, 1), (0, 2), (2, 2)]
    won = False
    for row, col in moves:
        won = g.move(row, col)
    assert won
    assert g.winner.name == "Alice"


def test_draw_detection():
    g = Game(3, new_test_players())
    moves = [
        (0, 0), (0, 1),
        (0, 2), (1, 1),
        (1, 0), (1, 2),
        (2, 1), (2, 0),
        (2, 2),
    ]
    last_won = False
    for row, col in moves:
        last_won = g.move(row, col)
    assert not last_won
    assert g.status == GameStatus.DRAW


def test_occupied_and_out_of_bounds_rejected():
    g = Game(3, new_test_players())
    g.move(0, 0)
    with pytest.raises(CellOccupiedError):
        g.move(0, 0)
    with pytest.raises(OutOfBoundsError):
        g.move(5, 5)


def test_move_rejected_after_game_over():
    g = Game(3, new_test_players())
    moves = [(0, 0), (1, 0), (0, 1), (1, 1), (0, 2)]
    for row, col in moves:
        g.move(row, col)
    assert g.status == GameStatus.WON
    with pytest.raises(GameOverError):
        g.move(2, 2)


def test_nxn_generalization():
    g = Game(5, new_test_players())
    moves = [
        (0, 4), (0, 0),
        (1, 3), (0, 1),
        (2, 2), (0, 2),
        (3, 1), (0, 3),
        (4, 0),
    ]
    won = False
    winner = None
    for row, col in moves:
        if g.move(row, col):
            won = True
            winner = g.winner
    assert won
    assert winner.name == "Alice"
