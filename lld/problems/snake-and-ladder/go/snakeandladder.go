// Package snakeandladder implements the classic Snake and Ladder LLD problem:
// a configurable board with snakes and ladders, an injectable dice (Strategy
// pattern) for deterministic testing, and a turn-based game engine that
// enforces the exact-landing-to-win rule.
package snakeandladder

import (
	"errors"
	"fmt"
	"math/rand"
)

// Dice rolls a value each turn. Implementations can be random or scripted.
type Dice interface {
	Roll() int
}

// StandardDice is a fair single six-sided die, seedable for reproducible runs.
type StandardDice struct {
	rng *rand.Rand
}

func NewStandardDice(seed int64) *StandardDice {
	return &StandardDice{rng: rand.New(rand.NewSource(seed))}
}

func (d *StandardDice) Roll() int {
	return d.rng.Intn(6) + 1
}

// ScriptedDice replays a fixed sequence of rolls, cycling once exhausted.
// This is the key to writing deterministic, non-flaky tests for the engine.
type ScriptedDice struct {
	Rolls []int
	pos   int
}

func NewScriptedDice(rolls ...int) *ScriptedDice {
	return &ScriptedDice{Rolls: rolls}
}

func (d *ScriptedDice) Roll() int {
	if len(d.Rolls) == 0 {
		return 1
	}
	v := d.Rolls[d.pos%len(d.Rolls)]
	d.pos++
	return v
}

// Entity is a board hazard/shortcut: a snake (head->tail) or ladder (bottom->top).
type Entity struct {
	Start int
	End   int
}

// Board holds the size and the snake/ladder map keyed by the cell a player lands on.
type Board struct {
	Size     int
	entities map[int]int
}

// NewBoard builds a board of the given size (last cell = size) with the supplied
// snakes and ladders. Entities with overlapping start/end are rejected to avoid
// ambiguous or infinite-loop configurations.
func NewBoard(size int, snakes, ladders []Entity) (*Board, error) {
	b := &Board{Size: size, entities: make(map[int]int)}
	for _, s := range snakes {
		if s.Start <= s.End {
			return nil, fmt.Errorf("snake start %d must be greater than end %d", s.Start, s.End)
		}
		if err := b.addEntity(s.Start, s.End); err != nil {
			return nil, err
		}
	}
	for _, l := range ladders {
		if l.Start >= l.End {
			return nil, fmt.Errorf("ladder start %d must be less than end %d", l.Start, l.End)
		}
		if err := b.addEntity(l.Start, l.End); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (b *Board) addEntity(from, to int) error {
	if from <= 1 || from >= b.Size {
		return fmt.Errorf("entity start %d must be within (1, %d)", from, b.Size)
	}
	if _, exists := b.entities[from]; exists {
		return fmt.Errorf("cell %d already has a snake or ladder", from)
	}
	b.entities[from] = to
	return nil
}

// resolve applies any snake/ladder at the landed-on cell, returning the final resting cell.
func (b *Board) resolve(cell int) int {
	if dest, ok := b.entities[cell]; ok {
		return dest
	}
	return cell
}

type Player struct {
	Name     string
	Position int
}

var (
	ErrNotEnoughPlayers = errors.New("need at least two players")
	ErrGameAlreadyOver  = errors.New("game already has a winner")
)

// Game drives turn order, dice rolls, and win detection for a single match.
type Game struct {
	Board   *Board
	Dice    Dice
	Players []*Player
	turn    int
	Winner  *Player
}

func NewGame(board *Board, dice Dice, playerNames []string) (*Game, error) {
	if len(playerNames) < 2 {
		return nil, ErrNotEnoughPlayers
	}
	players := make([]*Player, len(playerNames))
	for i, name := range playerNames {
		players[i] = &Player{Name: name, Position: 0}
	}
	return &Game{Board: board, Dice: dice, Players: players}, nil
}

// PlayTurn rolls the dice for the current player, applies the exact-landing rule
// and any snake/ladder at the destination, then rotates turn order to the next
// player. It returns the player who just moved.
func (g *Game) PlayTurn() (*Player, error) {
	if g.Winner != nil {
		return nil, ErrGameAlreadyOver
	}

	player := g.Players[g.turn]
	roll := g.Dice.Roll()
	target := player.Position + roll

	// Overshooting the last cell is not a legal move: the player stays put.
	if target <= g.Board.Size {
		player.Position = g.Board.resolve(target)
	}

	if player.Position == g.Board.Size {
		g.Winner = player
	} else {
		g.turn = (g.turn + 1) % len(g.Players)
	}
	return player, nil
}

// CurrentPlayer returns whose turn it is next.
func (g *Game) CurrentPlayer() *Player {
	return g.Players[g.turn]
}

// Play runs turns until a winner emerges, guarding against runaway loops.
func (g *Game) Play(maxTurns int) (*Player, error) {
	for i := 0; i < maxTurns; i++ {
		if _, err := g.PlayTurn(); err != nil {
			return nil, err
		}
		if g.Winner != nil {
			return g.Winner, nil
		}
	}
	return nil, fmt.Errorf("no winner after %d turns", maxTurns)
}
