// Package chess implements a simplified Chess LLD: 8x8 board, per-piece move
// generation (Template/Strategy pattern via the Piece interface), turn
// alternation, and check/checkmate detection. En passant, castling, and pawn
// promotion are intentionally out of scope (see README talking points).
package chess

import (
	"errors"
	"fmt"
)

type Color int

const (
	White Color = iota
	Black
)

func (c Color) Opposite() Color {
	if c == White {
		return Black
	}
	return White
}

func (c Color) String() string {
	if c == White {
		return "White"
	}
	return "Black"
}

type PieceType int

const (
	Pawn PieceType = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

type Position struct {
	Row, Col int
}

func (p Position) inBounds() bool {
	return p.Row >= 0 && p.Row < 8 && p.Col >= 0 && p.Col < 8
}

// ParseSquare converts algebraic notation ("e2") into a Position.
func ParseSquare(s string) (Position, error) {
	if len(s) != 2 {
		return Position{}, fmt.Errorf("invalid square %q", s)
	}
	col := int(s[0] - 'a')
	row := 8 - int(s[1]-'0')
	pos := Position{Row: row, Col: col}
	if !pos.inBounds() {
		return Position{}, fmt.Errorf("invalid square %q", s)
	}
	return pos, nil
}

// Piece is implemented by each piece type; CanMove encodes that piece's
// movement pattern plus path-clearance, but not whether the move leaves the
// mover's own king in check (Game is responsible for that).
type Piece interface {
	Color() Color
	Type() PieceType
	CanMove(b *Board, from, to Position) bool
}

type basePiece struct {
	color Color
}

func (p basePiece) Color() Color { return p.color }

func targetOK(b *Board, color Color, to Position) bool {
	occupant := b.At(to)
	return occupant == nil || occupant.Color() != color
}

type PawnPiece struct{ basePiece }
type KnightPiece struct{ basePiece }
type BishopPiece struct{ basePiece }
type RookPiece struct{ basePiece }
type QueenPiece struct{ basePiece }
type KingPiece struct{ basePiece }

func NewPiece(t PieceType, c Color) Piece {
	base := basePiece{color: c}
	switch t {
	case Pawn:
		return PawnPiece{base}
	case Knight:
		return KnightPiece{base}
	case Bishop:
		return BishopPiece{base}
	case Rook:
		return RookPiece{base}
	case Queen:
		return QueenPiece{base}
	case King:
		return KingPiece{base}
	}
	panic("unknown piece type")
}

func (PawnPiece) Type() PieceType   { return Pawn }
func (KnightPiece) Type() PieceType { return Knight }
func (BishopPiece) Type() PieceType { return Bishop }
func (RookPiece) Type() PieceType   { return Rook }
func (QueenPiece) Type() PieceType  { return Queen }
func (KingPiece) Type() PieceType   { return King }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (p PawnPiece) CanMove(b *Board, from, to Position) bool {
	dir := 1
	startRow := 1
	if p.color == White {
		dir = -1
		startRow = 6
	}
	dr := to.Row - from.Row
	dc := to.Col - from.Col

	if dc == 0 {
		if dr == dir {
			return b.At(to) == nil
		}
		if dr == 2*dir && from.Row == startRow {
			mid := Position{Row: from.Row + dir, Col: from.Col}
			return b.At(mid) == nil && b.At(to) == nil
		}
		return false
	}
	if abs(dc) == 1 && dr == dir {
		occupant := b.At(to)
		return occupant != nil && occupant.Color() != p.color
	}
	return false
}

func (p KnightPiece) CanMove(b *Board, from, to Position) bool {
	dr, dc := abs(to.Row-from.Row), abs(to.Col-from.Col)
	if !((dr == 1 && dc == 2) || (dr == 2 && dc == 1)) {
		return false
	}
	return targetOK(b, p.color, to)
}

func (p BishopPiece) CanMove(b *Board, from, to Position) bool {
	dr, dc := to.Row-from.Row, to.Col-from.Col
	if abs(dr) != abs(dc) || dr == 0 {
		return false
	}
	return b.clearPath(from, to) && targetOK(b, p.color, to)
}

func (p RookPiece) CanMove(b *Board, from, to Position) bool {
	dr, dc := to.Row-from.Row, to.Col-from.Col
	if (dr == 0) == (dc == 0) {
		return false
	}
	return b.clearPath(from, to) && targetOK(b, p.color, to)
}

func (p QueenPiece) CanMove(b *Board, from, to Position) bool {
	dr, dc := to.Row-from.Row, to.Col-from.Col
	straight := (dr == 0) != (dc == 0)
	diagonal := dr != 0 && abs(dr) == abs(dc)
	if !straight && !diagonal {
		return false
	}
	return b.clearPath(from, to) && targetOK(b, p.color, to)
}

func (p KingPiece) CanMove(b *Board, from, to Position) bool {
	dr, dc := abs(to.Row-from.Row), abs(to.Col-from.Col)
	if dr > 1 || dc > 1 || (dr == 0 && dc == 0) {
		return false
	}
	return targetOK(b, p.color, to)
}

// Board holds the 8x8 grid of squares. Squares is a fixed array (not a
// slice), so Clone can deep-copy it via plain assignment for check-safety
// simulation without a manual nested loop.
type Board struct {
	squares [8][8]Piece
}

func NewBoard() *Board {
	b := &Board{}
	backRank := []PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for col, t := range backRank {
		b.squares[0][col] = NewPiece(t, Black)
		b.squares[7][col] = NewPiece(t, White)
		b.squares[1][col] = NewPiece(Pawn, Black)
		b.squares[6][col] = NewPiece(Pawn, White)
	}
	return b
}

func (b *Board) At(pos Position) Piece {
	return b.squares[pos.Row][pos.Col]
}

func (b *Board) set(pos Position, p Piece) {
	b.squares[pos.Row][pos.Col] = p
}

// move relocates a piece with no legality checks; callers must validate first.
func (b *Board) move(from, to Position) Piece {
	captured := b.At(to)
	b.set(to, b.At(from))
	b.set(from, nil)
	return captured
}

func (b *Board) Clone() *Board {
	nb := &Board{}
	nb.squares = b.squares
	return nb
}

func sign(n int) int {
	if n > 0 {
		return 1
	}
	if n < 0 {
		return -1
	}
	return 0
}

// clearPath reports whether all squares strictly between from and to (on a
// straight or diagonal line) are empty. Callers already verified the line is
// straight/diagonal; the destination square itself is not checked here.
func (b *Board) clearPath(from, to Position) bool {
	stepR, stepC := sign(to.Row-from.Row), sign(to.Col-from.Col)
	r, c := from.Row+stepR, from.Col+stepC
	for r != to.Row || c != to.Col {
		if b.squares[r][c] != nil {
			return false
		}
		r += stepR
		c += stepC
	}
	return true
}

func (b *Board) FindKing(color Color) (Position, bool) {
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			p := b.squares[r][c]
			if p != nil && p.Type() == King && p.Color() == color {
				return Position{Row: r, Col: c}, true
			}
		}
	}
	return Position{}, false
}

// IsSquareAttacked reports whether any piece of attacker's color can move
// onto pos. Only meaningful when pos is occupied (used for check detection,
// where pos is always the defending king's square).
func (b *Board) IsSquareAttacked(pos Position, attacker Color) bool {
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			p := b.squares[r][c]
			if p != nil && p.Color() == attacker && p.CanMove(b, Position{Row: r, Col: c}, pos) {
				return true
			}
		}
	}
	return false
}

func (b *Board) IsInCheck(color Color) bool {
	kingPos, ok := b.FindKing(color)
	if !ok {
		return false
	}
	return b.IsSquareAttacked(kingPos, color.Opposite())
}

type Move struct {
	From, To Position
	Piece    Piece
	Captured Piece
}

var (
	ErrNotYourTurn           = errors.New("not your turn")
	ErrNoPieceAtSource       = errors.New("no piece at source square")
	ErrOwnPieceAtTarget      = errors.New("cannot capture your own piece")
	ErrIllegalMove           = errors.New("illegal move for this piece")
	ErrMoveLeavesKingInCheck = errors.New("move would leave your own king in check")
)

// Game orchestrates turn alternation and rule enforcement on top of Board.
type Game struct {
	Board   *Board
	Turn    Color
	History []Move
}

func NewGame() *Game {
	return &Game{Board: NewBoard(), Turn: White}
}

// Move validates and applies a move, then switches turns.
func (g *Game) Move(from, to Position) error {
	piece := g.Board.At(from)
	if piece == nil {
		return ErrNoPieceAtSource
	}
	if piece.Color() != g.Turn {
		return ErrNotYourTurn
	}
	if target := g.Board.At(to); target != nil && target.Color() == piece.Color() {
		return ErrOwnPieceAtTarget
	}
	if !piece.CanMove(g.Board, from, to) {
		return ErrIllegalMove
	}

	trial := g.Board.Clone()
	trial.move(from, to)
	if trial.IsInCheck(piece.Color()) {
		return ErrMoveLeavesKingInCheck
	}

	captured := g.Board.move(from, to)
	g.History = append(g.History, Move{From: from, To: to, Piece: piece, Captured: captured})
	g.Turn = g.Turn.Opposite()
	return nil
}

func (g *Game) IsInCheck(color Color) bool {
	return g.Board.IsInCheck(color)
}

// IsCheckmate reports whether color is in check with no legal move (by any
// of its pieces, to any square) that escapes check. This brute-forces every
// (piece, destination) pair rather than generating a minimal legal-move set,
// which is intentionally simple at interview scope.
func (g *Game) IsCheckmate(color Color) bool {
	if !g.Board.IsInCheck(color) {
		return false
	}
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			piece := g.Board.squares[r][c]
			if piece == nil || piece.Color() != color {
				continue
			}
			from := Position{Row: r, Col: c}
			for tr := 0; tr < 8; tr++ {
				for tc := 0; tc < 8; tc++ {
					to := Position{Row: tr, Col: tc}
					if !piece.CanMove(g.Board, from, to) {
						continue
					}
					trial := g.Board.Clone()
					trial.move(from, to)
					if !trial.IsInCheck(color) {
						return false
					}
				}
			}
		}
	}
	return true
}
