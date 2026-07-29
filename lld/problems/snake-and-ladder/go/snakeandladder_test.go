package snakeandladder

import "testing"

func smallBoard(t *testing.T, snakes, ladders []Entity) *Board {
	t.Helper()
	board, err := NewBoard(10, snakes, ladders)
	if err != nil {
		t.Fatalf("unexpected error building board: %v", err)
	}
	return board
}

func TestHappyPathDeterministicWin(t *testing.T) {
	board := smallBoard(t, []Entity{{Start: 9, End: 3}}, []Entity{{Start: 4, End: 8}})
	// Cycle 4,3,3 drives: P1->4(ladder->8), P2->3, P1->11 overshoot(stays 8),
	// P2->7, P1->11 overshoot(stays 8), P2->10 exact -> wins.
	dice := NewScriptedDice(4, 3, 3)
	game, err := NewGame(board, dice, []string{"P1", "P2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	winner, err := game.Play(6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if winner.Name != "P2" {
		t.Fatalf("expected P2 to win, got %s", winner.Name)
	}
	if winner.Position != board.Size {
		t.Fatalf("expected winner at last cell %d, got %d", board.Size, winner.Position)
	}
}

func TestExactLandingOvershootDoesNotMove(t *testing.T) {
	board := smallBoard(t, nil, nil)
	dice := NewScriptedDice(8, 1, 3)
	game, err := NewGame(board, dice, []string{"P1", "P2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, _ := game.PlayTurn() // P1: 0 + 8 = 8
	if p.Position != 8 {
		t.Fatalf("expected P1 at 8, got %d", p.Position)
	}
	game.PlayTurn() // P2: 0 + 1 = 1

	p, err = game.PlayTurn() // P1: 8 + 3 = 11, overshoots the last cell (10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Position != 8 {
		t.Fatalf("expected overshoot roll to leave P1 at 8, got %d", p.Position)
	}
	if game.Winner != nil {
		t.Fatalf("expected no winner from an overshoot roll, got %v", game.Winner)
	}
}

func TestLadderMovesPlayerForward(t *testing.T) {
	board := smallBoard(t, nil, []Entity{{Start: 4, End: 8}})
	game, err := NewGame(board, NewScriptedDice(4), []string{"P1", "P2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, _ := game.PlayTurn()
	if p.Position != 8 {
		t.Fatalf("expected ladder from 4 to carry P1 to 8, got %d", p.Position)
	}
}

func TestSnakeMovesPlayerBackward(t *testing.T) {
	board := smallBoard(t, []Entity{{Start: 6, End: 2}}, nil)
	game, err := NewGame(board, NewScriptedDice(6), []string{"P1", "P2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, _ := game.PlayTurn()
	if p.Position != 2 {
		t.Fatalf("expected snake from 6 to send P1 back to 2, got %d", p.Position)
	}
}

func TestTurnRotationAmongMultiplePlayers(t *testing.T) {
	board, err := NewBoard(100, nil, nil) // large enough that nobody wins mid-test
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	game, err := NewGame(board, NewScriptedDice(1), []string{"P1", "P2", "P3"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var order []string
	for i := 0; i < 6; i++ {
		p, err := game.PlayTurn()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		order = append(order, p.Name)
	}

	expected := []string{"P1", "P2", "P3", "P1", "P2", "P3"}
	for i, name := range expected {
		if order[i] != name {
			t.Fatalf("turn order mismatch at index %d: expected %s, got %s", i, name, order[i])
		}
	}
}

func TestNotEnoughPlayers(t *testing.T) {
	board := smallBoard(t, nil, nil)
	_, err := NewGame(board, NewScriptedDice(1), []string{"solo"})
	if err != ErrNotEnoughPlayers {
		t.Fatalf("expected ErrNotEnoughPlayers, got %v", err)
	}
}

func TestPlayTurnAfterWinReturnsError(t *testing.T) {
	board := smallBoard(t, nil, nil)
	dice := NewScriptedDice(10, 1)
	game, err := NewGame(board, dice, []string{"P1", "P2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	winner, err := game.Play(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if winner.Name != "P1" {
		t.Fatalf("expected P1 to win in one roll, got %s", winner.Name)
	}

	if _, err := game.PlayTurn(); err != ErrGameAlreadyOver {
		t.Fatalf("expected ErrGameAlreadyOver, got %v", err)
	}
}
