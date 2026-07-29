package tictactoe

import "testing"

func newTestPlayers() []*Player {
	return []*Player{
		{Name: "Alice", Symbol: "X"},
		{Name: "Bob", Symbol: "O"},
	}
}

func TestRowWin(t *testing.T) {
	g := NewGame(3, newTestPlayers())
	moves := [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}}
	var won bool
	var err error
	for _, m := range moves {
		won, err = g.Move(m[0], m[1])
		if err != nil {
			t.Fatalf("unexpected error on move %v: %v", m, err)
		}
	}
	if !won {
		t.Fatalf("expected the last move to win")
	}
	if g.Status() != Won {
		t.Fatalf("expected status Won, got %v", g.Status())
	}
	if g.Winner().Name != "Alice" {
		t.Fatalf("expected Alice to win, got %s", g.Winner().Name)
	}
}

func TestColumnWin(t *testing.T) {
	g := NewGame(3, newTestPlayers())
	moves := [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}, {2, 0}}
	var won bool
	for _, m := range moves {
		w, err := g.Move(m[0], m[1])
		if err != nil {
			t.Fatalf("unexpected error on move %v: %v", m, err)
		}
		won = w
	}
	if !won || g.Winner().Name != "Alice" {
		t.Fatalf("expected Alice to win via column, got winner=%v status=%v", g.Winner(), g.Status())
	}
}

func TestDiagonalWin(t *testing.T) {
	g := NewGame(3, newTestPlayers())
	moves := [][2]int{{0, 0}, {0, 1}, {1, 1}, {0, 2}, {2, 2}}
	var won bool
	for _, m := range moves {
		w, err := g.Move(m[0], m[1])
		if err != nil {
			t.Fatalf("unexpected error on move %v: %v", m, err)
		}
		won = w
	}
	if !won || g.Winner().Name != "Alice" {
		t.Fatalf("expected Alice to win via diagonal, got winner=%v status=%v", g.Winner(), g.Status())
	}
}

func TestDrawDetection(t *testing.T) {
	g := NewGame(3, newTestPlayers())
	// X O X
	// X O O
	// O X X
	moves := [][2]int{
		{0, 0}, {0, 1}, // X O
		{0, 2}, {1, 1}, // X O
		{1, 0}, {1, 2}, // X O
		{2, 1}, {2, 0}, // X O
		{2, 2},         // X
	}
	var lastWon bool
	for i, m := range moves {
		won, err := g.Move(m[0], m[1])
		if err != nil {
			t.Fatalf("unexpected error on move %d %v: %v", i, m, err)
		}
		lastWon = won
	}
	if lastWon {
		t.Fatalf("expected no winner in a drawn game")
	}
	if g.Status() != Draw {
		t.Fatalf("expected status Draw, got %v", g.Status())
	}
}

func TestOccupiedAndOutOfBoundsRejected(t *testing.T) {
	g := NewGame(3, newTestPlayers())
	if _, err := g.Move(0, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := g.Move(0, 0); err != ErrCellOccupied {
		t.Fatalf("expected ErrCellOccupied, got %v", err)
	}
	if _, err := g.Move(5, 5); err != ErrOutOfBounds {
		t.Fatalf("expected ErrOutOfBounds, got %v", err)
	}
}

func TestMoveRejectedAfterGameOver(t *testing.T) {
	g := NewGame(3, newTestPlayers())
	moves := [][2]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}}
	for _, m := range moves {
		if _, err := g.Move(m[0], m[1]); err != nil {
			t.Fatalf("unexpected error on move %v: %v", m, err)
		}
	}
	if g.Status() != Won {
		t.Fatalf("expected game to be won before extra move")
	}
	if _, err := g.Move(2, 2); err != ErrGameOver {
		t.Fatalf("expected ErrGameOver, got %v", err)
	}
}

func TestNxNGeneralization(t *testing.T) {
	g := NewGame(5, newTestPlayers())
	// Alice plays the full anti-diagonal of a 5x5 board: (0,4) (1,3) (2,2) (3,1) (4,0)
	// Bob's moves are irrelevant fillers off that diagonal.
	moves := [][2]int{
		{0, 4}, {0, 0},
		{1, 3}, {0, 1},
		{2, 2}, {0, 2},
		{3, 1}, {0, 3},
		{4, 0},
	}
	var won bool
	var winner *Player
	for _, m := range moves {
		w, err := g.Move(m[0], m[1])
		if err != nil {
			t.Fatalf("unexpected error on move %v: %v", m, err)
		}
		if w {
			won = true
			winner = g.Winner()
		}
	}
	if !won {
		t.Fatalf("expected a win on the 5x5 anti-diagonal")
	}
	if winner.Name != "Alice" {
		t.Fatalf("expected Alice to win, got %s", winner.Name)
	}
}
