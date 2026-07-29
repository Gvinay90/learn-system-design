// Package proxy demonstrates the Proxy pattern: BankAccountProxy controls
// access to a RealBankAccount, enforcing owner authorization before any
// withdrawal or deposit reaches the real object.
package proxy

import (
	"errors"
	"fmt"
)

var ErrUnauthorized = errors.New("unauthorized: requester is not the account owner")
var ErrInsufficientFunds = errors.New("insufficient funds")

// Account is the common interface shared by the real subject and its proxy,
// so callers can use either interchangeably.
type Account interface {
	Deposit(requester string, amount float64) error
	Withdraw(requester string, amount float64) error
	Balance() float64
}

// RealBankAccount is the expensive/sensitive real subject. It performs no
// authorization checks of its own — that's the proxy's job.
type RealBankAccount struct {
	Owner   string
	balance float64
}

func NewRealBankAccount(owner string, initialBalance float64) *RealBankAccount {
	return &RealBankAccount{Owner: owner, balance: initialBalance}
}

func (a *RealBankAccount) Deposit(_ string, amount float64) error {
	a.balance += amount
	return nil
}

func (a *RealBankAccount) Withdraw(_ string, amount float64) error {
	if amount > a.balance {
		return ErrInsufficientFunds
	}
	a.balance -= amount
	return nil
}

func (a *RealBankAccount) Balance() float64 { return a.balance }

// BankAccountProxy is a protection proxy: it enforces that only the account
// owner may deposit or withdraw, and logs every access attempt, without the
// real account needing to know about authorization at all.
type BankAccountProxy struct {
	real      *RealBankAccount
	AccessLog []string
}

func NewBankAccountProxy(real *RealBankAccount) *BankAccountProxy {
	return &BankAccountProxy{real: real}
}

func (p *BankAccountProxy) Deposit(requester string, amount float64) error {
	p.AccessLog = append(p.AccessLog, fmt.Sprintf("deposit attempt by %s", requester))
	if requester != p.real.Owner {
		return ErrUnauthorized
	}
	return p.real.Deposit(requester, amount)
}

func (p *BankAccountProxy) Withdraw(requester string, amount float64) error {
	p.AccessLog = append(p.AccessLog, fmt.Sprintf("withdraw attempt by %s", requester))
	if requester != p.real.Owner {
		return ErrUnauthorized
	}
	return p.real.Withdraw(requester, amount)
}

func (p *BankAccountProxy) Balance() float64 { return p.real.Balance() }
