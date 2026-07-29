// Package tictactoe implements the classic Tic-Tac-Toe LLD problem on a
// generic NxN board: move validation, alternating turns between N players,
// and win detection via incrementally maintained row/column/diagonal counts.
package tictactoe

import (
	"errors"
	"fmt"
)

type Symbol string

const Empty Symbol = ""

type Player struct {
	Name   string
	Symbol Symbol
}

type GameStatus int

const (
	InProgress GameStatus = iota
	Won
	Draw
)

var (
	ErrOutOfBounds  = errors.New("move is out of bounds")
	ErrCellOccupied = errors.New("cell is already occupied")
	ErrGameOver     = errors.New("game is already over")
)

// Board is a generic NxN grid. Win detection is O(1) per move: instead of
// rescanning the whole board, we keep running counts per row, column and
// the two diagonals for whichever symbol is placed. A move only ever
// affects one row, one column and (at most) two diagonals, so checking
// those counts after each move is enough to detect a win.
type Board struct {
	size        int
	cells       [][]Symbol
	rowCounts   map[Symbol][]int
	colCounts   map[Symbol][]int
	diagCount   map[Symbol]int
	antiCount   map[Symbol]int
	filledCells int
}

func NewBoard(size int) *Board {
	cells := make([][]Symbol, size)
	for i := range cells {
		cells[i] = make([]Symbol, size)
	}
	return &Board{
		size:      size,
		cells:     cells,
		rowCounts: make(map[Symbol][]int),
		colCounts: make(map[Symbol][]int),
		diagCount: make(map[Symbol]int),
		antiCount: make(map[Symbol]int),
	}
}

func (b *Board) Size() int { return b.size }

func (b *Board) At(row, col int) Symbol { return b.cells[row][col] }

func (b *Board) inBounds(row, col int) bool {
	return row >= 0 && row < b.size && col >= 0 && col < b.size
}

func (b *Board) IsFull() bool { return b.filledCells == b.size*b.size }

// place records sym at (row, col) and returns true if that move completes
// a row, column, or diagonal of length size for sym.
func (b *Board) place(row, col int, sym Symbol) (bool, error) {
	if !b.inBounds(row, col) {
		return false, ErrOutOfBounds
	}
	if b.cells[row][col] != Empty {
		return false, ErrCellOccupied
	}
	b.cells[row][col] = sym
	b.filledCells++

	if _, ok := b.rowCounts[sym]; !ok {
		b.rowCounts[sym] = make([]int, b.size)
		b.colCounts[sym] = make([]int, b.size)
	}
	b.rowCounts[sym][row]++
	b.colCounts[sym][col]++
	if row == col {
		b.diagCount[sym]++
	}
	if row+col == b.size-1 {
		b.antiCount[sym]++
	}

	won := b.rowCounts[sym][row] == b.size ||
		b.colCounts[sym][col] == b.size ||
		b.diagCount[sym] == b.size ||
		b.antiCount[sym] == b.size
	return won, nil
}

func (b *Board) String() string {
	s := ""
	for r := 0; r < b.size; r++ {
		for c := 0; c < b.size; c++ {
			ch := b.cells[r][c]
			if ch == Empty {
				ch = "."
			}
			s += fmt.Sprintf("%s ", ch)
		}
		s += "\n"
	}
	return s
}

// Game orchestrates turn order, move validation and terminal-state tracking
// on top of a Board. It has no locking: unlike the parking lot, a single
// tic-tac-toe game is played by one turn owner at a time, so there is no
// concurrent-writer scenario to guard against.
type Game struct {
	board   *Board
	players []*Player
	turn    int
	status  GameStatus
	winner  *Player
}

func NewGame(size int, players []*Player) *Game {
	return &Game{board: NewBoard(size), players: players}
}

func (g *Game) Board() *Board { return g.board }

func (g *Game) Status() GameStatus { return g.status }

func (g *Game) Winner() *Player { return g.winner }

func (g *Game) CurrentPlayer() *Player { return g.players[g.turn] }

// Move plays a cell for the current player, advances the turn, and updates
// game status. Returns the player who just moved won and any validation error.
func (g *Game) Move(row, col int) (bool, error) {
	if g.status != InProgress {
		return false, ErrGameOver
	}

	player := g.players[g.turn]
	won, err := g.board.place(row, col, player.Symbol)
	if err != nil {
		return false, err
	}

	if won {
		g.status = Won
		g.winner = player
		return true, nil
	}
	if g.board.IsFull() {
		g.status = Draw
		return false, nil
	}

	g.turn = (g.turn + 1) % len(g.players)
	return false, nil
}
