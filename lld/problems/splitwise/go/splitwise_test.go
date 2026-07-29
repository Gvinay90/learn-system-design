package splitwise

import (
	"math"
	"testing"
)

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}

func TestEqualSplitCreatesCorrectBalances(t *testing.T) {
	sw := NewSplitwise()
	alice := sw.AddUser("u1", "Alice")
	bob := sw.AddUser("u2", "Bob")
	carol := sw.AddUser("u3", "Carol")

	_, err := sw.AddExpense("Dinner", alice, 90, []*User{alice, bob, carol}, EqualSplit{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := sw.Ledger.NetBalance(bob.ID, alice.ID); !approxEqual(got, 30) {
		t.Fatalf("expected bob to owe alice 30, got %v", got)
	}
	if got := sw.Ledger.NetBalance(carol.ID, alice.ID); !approxEqual(got, 30) {
		t.Fatalf("expected carol to owe alice 30, got %v", got)
	}
	if got := sw.Ledger.NetBalance(alice.ID, bob.ID); !approxEqual(got, -30) {
		t.Fatalf("expected anti-symmetric balance alice->bob of -30, got %v", got)
	}
}

func TestExactSplitValidation(t *testing.T) {
	sw := NewSplitwise()
	alice := sw.AddUser("u1", "Alice")
	bob := sw.AddUser("u2", "Bob")

	_, err := sw.AddExpense("Rent", alice, 100, []*User{alice, bob},
		ExactSplit{Amounts: map[string]float64{alice.ID: 40, bob.ID: 50}})
	if err == nil {
		t.Fatal("expected error for exact amounts not summing to total, got nil")
	}

	_, err = sw.AddExpense("Rent", alice, 100, []*User{alice, bob},
		ExactSplit{Amounts: map[string]float64{alice.ID: 40, bob.ID: 60}})
	if err != nil {
		t.Fatalf("expected valid exact split to succeed, got %v", err)
	}
	if got := sw.Ledger.NetBalance(bob.ID, alice.ID); !approxEqual(got, 60) {
		t.Fatalf("expected bob to owe alice 60, got %v", got)
	}
}

func TestPercentSplitCorrectness(t *testing.T) {
	sw := NewSplitwise()
	alice := sw.AddUser("u1", "Alice")
	bob := sw.AddUser("u2", "Bob")
	carol := sw.AddUser("u3", "Carol")

	_, err := sw.AddExpense("Trip", alice, 200, []*User{alice, bob, carol},
		PercentSplit{Percentages: map[string]float64{alice.ID: 50, bob.ID: 25, carol.ID: 25}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := sw.Ledger.NetBalance(bob.ID, alice.ID); !approxEqual(got, 50) {
		t.Fatalf("expected bob to owe alice 50, got %v", got)
	}
	if got := sw.Ledger.NetBalance(carol.ID, alice.ID); !approxEqual(got, 50) {
		t.Fatalf("expected carol to owe alice 50, got %v", got)
	}

	_, err = sw.AddExpense("Bad", alice, 200, []*User{alice, bob, carol},
		PercentSplit{Percentages: map[string]float64{alice.ID: 50, bob.ID: 25, carol.ID: 20}})
	if err == nil {
		t.Fatal("expected error for percentages not summing to 100, got nil")
	}
}

func TestSimplifyDebts(t *testing.T) {
	// A owes 300 net, B owes 100 net, C is owed 200, D is owed 200.
	net := map[string]float64{
		"A": -300,
		"B": -100,
		"C": 200,
		"D": 200,
	}

	txns := SimplifyDebts(net)

	if len(txns) > 3 {
		t.Fatalf("expected at most 3 transactions to settle 4 users, got %d: %v", len(txns), txns)
	}

	settled := map[string]float64{}
	for _, tx := range txns {
		settled[tx.From] -= tx.Amount
		settled[tx.To] += tx.Amount
	}
	for user, expected := range net {
		if !approxEqual(settled[user], expected) {
			t.Fatalf("user %s not settled correctly: expected net %v, settled %v", user, expected, settled[user])
		}
	}
}

func TestSimplifyDebtsMinimizesTransactionCount(t *testing.T) {
	// Three-way cycle: A owes B 10, B owes C 10, C owes A 10 -> net balances
	// are all zero, so zero transactions should be needed.
	net := map[string]float64{"A": 0, "B": 0, "C": 0}
	txns := SimplifyDebts(net)
	if len(txns) != 0 {
		t.Fatalf("expected 0 transactions for fully offsetting cycle, got %d: %v", len(txns), txns)
	}
}
