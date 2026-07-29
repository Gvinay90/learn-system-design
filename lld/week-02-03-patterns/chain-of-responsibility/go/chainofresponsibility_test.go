package chainofresponsibility

import "testing"

func newTestChain() Approver {
	return BuildChain(NewManager(1000), NewDirector(5000), NewVP(20000))
}

func TestManagerApprovesSmallExpense(t *testing.T) {
	chain := newTestChain()
	who, err := chain.Approve(&ExpenseRequest{Amount: 500, Description: "team lunch"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if who != "Manager" {
		t.Fatalf("expected Manager, got %s", who)
	}
}

func TestDirectorApprovesMidExpense(t *testing.T) {
	chain := newTestChain()
	who, err := chain.Approve(&ExpenseRequest{Amount: 3000, Description: "conference"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if who != "Director" {
		t.Fatalf("expected Director, got %s", who)
	}
}

func TestVPApprovesLargeExpense(t *testing.T) {
	chain := newTestChain()
	who, err := chain.Approve(&ExpenseRequest{Amount: 15000, Description: "new hire equipment"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if who != "VP" {
		t.Fatalf("expected VP, got %s", who)
	}
}

func TestNoApproverForExcessiveExpense(t *testing.T) {
	chain := newTestChain()
	_, err := chain.Approve(&ExpenseRequest{Amount: 100000, Description: "private jet"})
	if err != ErrNoApprover {
		t.Fatalf("expected ErrNoApprover, got %v", err)
	}
}

func TestChainOrderMatters(t *testing.T) {
	// A single-link chain of just the Manager cannot approve anything above its limit.
	chain := BuildChain(NewManager(1000))
	_, err := chain.Approve(&ExpenseRequest{Amount: 2000})
	if err != ErrNoApprover {
		t.Fatalf("expected ErrNoApprover with no further approvers, got %v", err)
	}
}
