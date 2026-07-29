// Package chainofresponsibility demonstrates the Chain of Responsibility
// pattern: an expense request travels along a chain of approvers (Manager,
// Director, VP), each of which either approves it (if within its limit) or
// forwards it to the next approver in the chain.
package chainofresponsibility

import "errors"

var ErrNoApprover = errors.New("no approver in the chain can approve this amount")

// ExpenseRequest is the request that travels along the chain.
type ExpenseRequest struct {
	Amount      float64
	Description string
}

// Approver is a handler in the chain. Each concrete approver decides whether
// it can approve the request, or passes it to the next link.
type Approver interface {
	SetNext(next Approver)
	Approve(req *ExpenseRequest) (approvedBy string, err error)
}

// baseApprover holds the shared "approval limit + next link" behavior so
// concrete approvers only need to supply their limit and name.
type baseApprover struct {
	name  string
	limit float64
	next  Approver
}

func (b *baseApprover) SetNext(next Approver) { b.next = next }

func (b *baseApprover) tryApprove(req *ExpenseRequest) (string, error) {
	if req.Amount <= b.limit {
		return b.name, nil
	}
	if b.next != nil {
		return b.next.Approve(req)
	}
	return "", ErrNoApprover
}

type Manager struct{ baseApprover }

func NewManager(limit float64) *Manager {
	return &Manager{baseApprover{name: "Manager", limit: limit}}
}

func (m *Manager) Approve(req *ExpenseRequest) (string, error) { return m.tryApprove(req) }

type Director struct{ baseApprover }

func NewDirector(limit float64) *Director {
	return &Director{baseApprover{name: "Director", limit: limit}}
}

func (d *Director) Approve(req *ExpenseRequest) (string, error) { return d.tryApprove(req) }

type VP struct{ baseApprover }

func NewVP(limit float64) *VP {
	return &VP{baseApprover{name: "VP", limit: limit}}
}

func (v *VP) Approve(req *ExpenseRequest) (string, error) { return v.tryApprove(req) }

// BuildChain wires approvers in order and returns the head of the chain.
func BuildChain(approvers ...Approver) Approver {
	for i := 0; i < len(approvers)-1; i++ {
		approvers[i].SetNext(approvers[i+1])
	}
	return approvers[0]
}
