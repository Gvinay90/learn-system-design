import pytest

from snake_and_ladder import (
    Board,
    Entity,
    Game,
    GameAlreadyOverError,
    NotEnoughPlayersError,
    ScriptedDice,
)


def small_board(snakes=None, ladders=None) -> Board:
    return Board(10, snakes or [], ladders or [])


def test_happy_path_deterministic_win():
    board = small_board(snakes=[Entity(9, 3)], ladders=[Entity(4, 8)])
    # Cycle 4,3,3 drives: P1->4(ladder->8), P2->3, P1->11 overshoot(stays 8),
    # P2->7, P1->11 overshoot(stays 8), P2->10 exact -> wins.
    dice = ScriptedDice(4, 3, 3)
    game = Game(board, dice, ["P1", "P2"])

    winner = game.play(6)
    assert winner.name == "P2"
    assert winner.position == board.size


def test_exact_landing_overshoot_does_not_move():
    board = small_board()
    dice = ScriptedDice(8, 1, 3)
    game = Game(board, dice, ["P1", "P2"])

    p = game.play_turn()  # P1: 0 + 8 = 8
    assert p.position == 8
    game.play_turn()  # P2: 0 + 1 = 1

    p = game.play_turn()  # P1: 8 + 3 = 11, overshoots the last cell (10)
    assert p.position == 8
    assert game.winner is None


def test_ladder_moves_player_forward():
    board = small_board(ladders=[Entity(4, 8)])
    game = Game(board, ScriptedDice(4), ["P1", "P2"])

    p = game.play_turn()
    assert p.position == 8


def test_snake_moves_player_backward():
    board = small_board(snakes=[Entity(6, 2)])
    game = Game(board, ScriptedDice(6), ["P1", "P2"])

    p = game.play_turn()
    assert p.position == 2


def test_turn_rotation_among_multiple_players():
    board = Board(100)  # large enough that nobody wins mid-test
    game = Game(board, ScriptedDice(1), ["P1", "P2", "P3"])

    order = [game.play_turn().name for _ in range(6)]

    assert order == ["P1", "P2", "P3", "P1", "P2", "P3"]


def test_not_enough_players():
    board = small_board()
    with pytest.raises(NotEnoughPlayersError):
        Game(board, ScriptedDice(1), ["solo"])


def test_play_turn_after_win_raises():
    board = small_board()
    dice = ScriptedDice(10, 1)
    game = Game(board, dice, ["P1", "P2"])

    winner = game.play(1)
    assert winner.name == "P1"

    with pytest.raises(GameAlreadyOverError):
        game.play_turn()


def test_invalid_snake_and_ladder_configs_rejected():
    from snake_and_ladder import BoardConfigError

    with pytest.raises(BoardConfigError):
        small_board(snakes=[Entity(3, 9)])  # snake must go start > end

    with pytest.raises(BoardConfigError):
        small_board(ladders=[Entity(9, 3)])  # ladder must go start < end

    with pytest.raises(BoardConfigError):
        small_board(snakes=[Entity(9, 3)], ladders=[Entity(9, 4)])  # duplicate cell
