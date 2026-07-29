package chess

import "testing"

func sq(t *testing.T, s string) Position {
	t.Helper()
	pos, err := ParseSquare(s)
	if err != nil {
		t.Fatalf("bad square %q: %v", s, err)
	}
	return pos
}

func TestPawnOpeningMoveHappyPath(t *testing.T) {
	g := NewGame()
	from, to := sq(t, "e2"), sq(t, "e4")

	if err := g.Move(from, to); err != nil {
		t.Fatalf("expected legal opening move, got err: %v", err)
	}
	if g.Board.At(to) == nil || g.Board.At(to).Type() != Pawn {
		t.Fatalf("expected pawn at e4")
	}
	if g.Board.At(from) != nil {
		t.Fatalf("expected e2 to be empty after move")
	}
	if g.Turn != Black {
		t.Fatalf("expected turn to switch to Black, got %v", g.Turn)
	}
}

func TestIllegalMoveRejected(t *testing.T) {
	g := NewGame()

	// Pawn cannot jump 3 squares.
	if err := g.Move(sq(t, "e2"), sq(t, "e5")); err != ErrIllegalMove {
		t.Fatalf("expected ErrIllegalMove for 3-square pawn move, got %v", err)
	}

	// Rook is blocked by its own pawn at the start of the game.
	if err := g.Move(sq(t, "a1"), sq(t, "a3")); err != ErrIllegalMove {
		t.Fatalf("expected ErrIllegalMove for blocked rook move, got %v", err)
	}
}

func TestCheckDetection(t *testing.T) {
	g := &Game{Board: &Board{}, Turn: White}
	g.Board.set(sq(t, "e1"), NewPiece(King, White))
	g.Board.set(sq(t, "e8"), NewPiece(Rook, Black))
	g.Board.set(sq(t, "a1"), NewPiece(King, Black))

	if !g.IsInCheck(White) {
		t.Fatalf("expected White king on e1 to be in check from rook on e8")
	}
	if g.IsCheckmate(White) {
		t.Fatalf("expected check but not checkmate: king can step aside")
	}
}

// TestFoolsMateCheckmate replays the fastest possible checkmate:
// 1. f3 e5 2. g4 Qh4#
func TestFoolsMateCheckmate(t *testing.T) {
	g := NewGame()
	moves := [][2]string{
		{"f2", "f3"}, // White
		{"e7", "e5"}, // Black
		{"g2", "g4"}, // White
		{"d8", "h4"}, // Black delivers checkmate
	}
	for _, m := range moves {
		if err := g.Move(sq(t, m[0]), sq(t, m[1])); err != nil {
			t.Fatalf("move %s-%s failed: %v", m[0], m[1], err)
		}
	}

	if !g.IsInCheck(White) {
		t.Fatalf("expected White king to be in check after Qh4")
	}
	if !g.IsCheckmate(White) {
		t.Fatalf("expected fool's mate to be checkmate")
	}
}
