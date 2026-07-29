// Package splitwise implements the classic Splitwise LLD problem:
// users, group expenses, pluggable split strategies (equal/exact/percent),
// a pairwise balance ledger, and greedy debt simplification.
package splitwise

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
)

const epsilon = 1e-6

type User struct {
	ID   string
	Name string
}

type Group struct {
	ID      string
	Name    string
	Members []*User
}

// SplitStrategy computes each participant's owed share for an expense amount.
// Implementations validate their own inputs (e.g. exact amounts summing to
// total, percentages summing to 100) and return an error otherwise.
type SplitStrategy interface {
	Compute(totalAmount float64, participants []*User) (map[string]float64, error)
}

type EqualSplit struct{}

func (EqualSplit) Compute(totalAmount float64, participants []*User) (map[string]float64, error) {
	if len(participants) == 0 {
		return nil, errors.New("equal split requires at least one participant")
	}
	share := totalAmount / float64(len(participants))
	shares := make(map[string]float64, len(participants))
	for _, p := range participants {
		shares[p.ID] = share
	}
	return shares, nil
}

// ExactSplit takes explicit amounts per user ID that must sum to the total.
type ExactSplit struct {
	Amounts map[string]float64
}

func (s ExactSplit) Compute(totalAmount float64, participants []*User) (map[string]float64, error) {
	sum := 0.0
	shares := make(map[string]float64, len(participants))
	for _, p := range participants {
		amt, ok := s.Amounts[p.ID]
		if !ok {
			return nil, fmt.Errorf("exact split missing amount for %s", p.ID)
		}
		shares[p.ID] = amt
		sum += amt
	}
	if math.Abs(sum-totalAmount) > epsilon {
		return nil, fmt.Errorf("exact split amounts sum to %.2f, want %.2f", sum, totalAmount)
	}
	return shares, nil
}

// PercentSplit takes percentages per user ID that must sum to 100.
type PercentSplit struct {
	Percentages map[string]float64
}

func (s PercentSplit) Compute(totalAmount float64, participants []*User) (map[string]float64, error) {
	sum := 0.0
	shares := make(map[string]float64, len(participants))
	for _, p := range participants {
		pct, ok := s.Percentages[p.ID]
		if !ok {
			return nil, fmt.Errorf("percent split missing percentage for %s", p.ID)
		}
		shares[p.ID] = totalAmount * pct / 100.0
		sum += pct
	}
	if math.Abs(sum-100.0) > epsilon {
		return nil, fmt.Errorf("percent split percentages sum to %.2f, want 100", sum)
	}
	return shares, nil
}

type Expense struct {
	ID           string
	Description  string
	PaidBy       *User
	Amount       float64
	Participants []*User
	Strategy     SplitStrategy
}

var ErrInvalidSplit = errors.New("invalid split")

// Ledger tracks net pairwise balances between users. Balances[a][b] > 0
// means b owes a that amount; it is always kept anti-symmetric
// (Balances[a][b] == -Balances[b][a]).
type Ledger struct {
	mu       sync.Mutex
	balances map[string]map[string]float64
}

func NewLedger() *Ledger {
	return &Ledger{balances: make(map[string]map[string]float64)}
}

func (l *Ledger) adjust(a, b string, amount float64) {
	if l.balances[a] == nil {
		l.balances[a] = make(map[string]float64)
	}
	if l.balances[b] == nil {
		l.balances[b] = make(map[string]float64)
	}
	l.balances[a][b] += amount
	l.balances[b][a] -= amount
}

// AddExpense splits the expense per its strategy and records that every
// participant (other than the payer) now owes the payer their share.
func (l *Ledger) AddExpense(e *Expense) error {
	shares, err := e.Strategy.Compute(e.Amount, e.Participants)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSplit, err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, p := range e.Participants {
		if p.ID == e.PaidBy.ID {
			continue
		}
		l.adjust(e.PaidBy.ID, p.ID, shares[p.ID])
	}
	return nil
}

// NetBalance returns how much debtor owes creditor (may be negative if the
// balance actually runs the other way).
func (l *Ledger) NetBalance(debtor, creditor string) float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.balances[creditor][debtor]
}

// NetBalances returns each user's overall net position: positive means the
// user is a net creditor (is owed money), negative means net debtor.
func (l *Ledger) NetBalances() map[string]float64 {
	l.mu.Lock()
	defer l.mu.Unlock()
	net := make(map[string]float64)
	for a, row := range l.balances {
		total := 0.0
		for _, amt := range row {
			total += amt
		}
		net[a] = total
	}
	return net
}

type Transaction struct {
	From   string
	To     string
	Amount float64
}

// SimplifyDebts takes each user's net balance (positive = creditor, negative
// = debtor) and greedily matches the largest creditor with the largest
// debtor, minimizing the number of settling transactions. This is the
// classic greedy algorithm: it does not guarantee the theoretical minimum in
// every adversarial case, but it is optimal for the common case and is what
// real Splitwise-style systems use.
func SimplifyDebts(net map[string]float64) []Transaction {
	type balance struct {
		id     string
		amount float64
	}
	var creditors, debtors []balance
	for id, amt := range net {
		if amt > epsilon {
			creditors = append(creditors, balance{id, amt})
		} else if amt < -epsilon {
			debtors = append(debtors, balance{id, -amt})
		}
	}

	sort.Slice(creditors, func(i, j int) bool { return creditors[i].amount > creditors[j].amount })
	sort.Slice(debtors, func(i, j int) bool { return debtors[i].amount > debtors[j].amount })

	var transactions []Transaction
	i, j := 0, 0
	for i < len(creditors) && j < len(debtors) {
		c, d := &creditors[i], &debtors[j]
		amount := math.Min(c.amount, d.amount)

		transactions = append(transactions, Transaction{From: d.id, To: c.id, Amount: amount})

		c.amount -= amount
		d.amount -= amount
		if c.amount <= epsilon {
			i++
		}
		if d.amount <= epsilon {
			j++
		}
	}
	return transactions
}

// Splitwise ties together users, groups and a ledger into the app-level API.
type Splitwise struct {
	mu     sync.Mutex
	Users  map[string]*User
	Groups map[string]*Group
	Ledger *Ledger
	seq    int
}

func NewSplitwise() *Splitwise {
	return &Splitwise{
		Users:  make(map[string]*User),
		Groups: make(map[string]*Group),
		Ledger: NewLedger(),
	}
}

func (s *Splitwise) AddUser(id, name string) *User {
	s.mu.Lock()
	defer s.mu.Unlock()
	u := &User{ID: id, Name: name}
	s.Users[id] = u
	return u
}

func (s *Splitwise) AddGroup(id, name string, members []*User) *Group {
	s.mu.Lock()
	defer s.mu.Unlock()
	g := &Group{ID: id, Name: name, Members: members}
	s.Groups[id] = g
	return g
}

// AddExpense creates an expense with the given strategy and records it in
// the ledger. Returns the expense (with a generated ID) or an error if the
// split is invalid.
func (s *Splitwise) AddExpense(description string, paidBy *User, amount float64, participants []*User, strategy SplitStrategy) (*Expense, error) {
	s.mu.Lock()
	s.seq++
	e := &Expense{
		ID:           fmt.Sprintf("E-%d", s.seq),
		Description:  description,
		PaidBy:       paidBy,
		Amount:       amount,
		Participants: participants,
		Strategy:     strategy,
	}
	s.mu.Unlock()

	if err := s.Ledger.AddExpense(e); err != nil {
		return nil, err
	}
	return e, nil
}
