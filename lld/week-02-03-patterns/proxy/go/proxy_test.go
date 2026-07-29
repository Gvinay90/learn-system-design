package proxy

import "testing"

func TestOwnerCanWithdraw(t *testing.T) {
	real := NewRealBankAccount("alice", 100)
	acc := NewBankAccountProxy(real)

	if err := acc.Withdraw("alice", 40); err != nil {
		t.Fatalf("expected owner withdraw to succeed, got %v", err)
	}
	if acc.Balance() != 60 {
		t.Fatalf("expected balance 60, got %v", acc.Balance())
	}
}

func TestNonOwnerCannotWithdraw(t *testing.T) {
	real := NewRealBankAccount("alice", 100)
	acc := NewBankAccountProxy(real)

	err := acc.Withdraw("mallory", 10)
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if acc.Balance() != 100 {
		t.Fatalf("expected balance unchanged at 100, got %v", acc.Balance())
	}
}

func TestNonOwnerCannotDeposit(t *testing.T) {
	real := NewRealBankAccount("alice", 100)
	acc := NewBankAccountProxy(real)

	err := acc.Deposit("mallory", 500)
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if acc.Balance() != 100 {
		t.Fatalf("expected balance unchanged, got %v", acc.Balance())
	}
}

func TestInsufficientFunds(t *testing.T) {
	real := NewRealBankAccount("alice", 20)
	acc := NewBankAccountProxy(real)

	err := acc.Withdraw("alice", 50)
	if err != ErrInsufficientFunds {
		t.Fatalf("expected ErrInsufficientFunds, got %v", err)
	}
}

func TestAccessLogRecordsAllAttempts(t *testing.T) {
	real := NewRealBankAccount("alice", 100)
	acc := NewBankAccountProxy(real)

	_ = acc.Withdraw("alice", 10)
	_ = acc.Deposit("mallory", 5)

	if len(acc.AccessLog) != 2 {
		t.Fatalf("expected 2 log entries, got %d: %v", len(acc.AccessLog), acc.AccessLog)
	}
}
