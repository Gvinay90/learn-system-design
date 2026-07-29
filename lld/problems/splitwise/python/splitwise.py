"""Splitwise LLD — Python reference implementation.

Users, group expenses, pluggable split strategies (equal/exact/percent),
a pairwise balance ledger, and greedy debt simplification.
See ../go/splitwise.go for the original design writeup.
"""
from __future__ import annotations

import threading
from dataclasses import dataclass, field
from typing import Dict, List, Optional, Protocol

EPSILON = 1e-6


@dataclass(frozen=True)
class User:
    id: str
    name: str


@dataclass
class Group:
    id: str
    name: str
    members: List[User]


class InvalidSplitError(Exception):
    pass


class SplitStrategy(Protocol):
    """Computes each participant's owed share for an expense amount.

    Implementations validate their own inputs (e.g. exact amounts summing
    to total, percentages summing to 100) and raise InvalidSplitError
    otherwise.
    """

    def compute(self, total_amount: float, participants: List[User]) -> Dict[str, float]:
        ...


class EqualSplit:
    def compute(self, total_amount: float, participants: List[User]) -> Dict[str, float]:
        if not participants:
            raise InvalidSplitError("equal split requires at least one participant")
        share = total_amount / len(participants)
        return {p.id: share for p in participants}


@dataclass
class ExactSplit:
    """Takes explicit amounts per user ID that must sum to the total."""

    amounts: Dict[str, float]

    def compute(self, total_amount: float, participants: List[User]) -> Dict[str, float]:
        total = 0.0
        shares: Dict[str, float] = {}
        for p in participants:
            if p.id not in self.amounts:
                raise InvalidSplitError(f"exact split missing amount for {p.id}")
            amt = self.amounts[p.id]
            shares[p.id] = amt
            total += amt
        if abs(total - total_amount) > EPSILON:
            raise InvalidSplitError(f"exact split amounts sum to {total:.2f}, want {total_amount:.2f}")
        return shares


@dataclass
class PercentSplit:
    """Takes percentages per user ID that must sum to 100."""

    percentages: Dict[str, float]

    def compute(self, total_amount: float, participants: List[User]) -> Dict[str, float]:
        total_pct = 0.0
        shares: Dict[str, float] = {}
        for p in participants:
            if p.id not in self.percentages:
                raise InvalidSplitError(f"percent split missing percentage for {p.id}")
            pct = self.percentages[p.id]
            shares[p.id] = total_amount * pct / 100.0
            total_pct += pct
        if abs(total_pct - 100.0) > EPSILON:
            raise InvalidSplitError(f"percent split percentages sum to {total_pct:.2f}, want 100")
        return shares


@dataclass
class Expense:
    id: str
    description: str
    paid_by: User
    amount: float
    participants: List[User]
    strategy: SplitStrategy


class Ledger:
    """Tracks net pairwise balances between users. balances[a][b] > 0 means
    b owes a that amount; it is always kept anti-symmetric
    (balances[a][b] == -balances[b][a]).
    """

    def __init__(self):
        self._lock = threading.Lock()
        self._balances: Dict[str, Dict[str, float]] = {}

    def _adjust(self, a: str, b: str, amount: float) -> None:
        self._balances.setdefault(a, {})
        self._balances.setdefault(b, {})
        self._balances[a][b] = self._balances[a].get(b, 0.0) + amount
        self._balances[b][a] = self._balances[b].get(a, 0.0) - amount

    def add_expense(self, expense: Expense) -> None:
        """Splits the expense per its strategy and records that every
        participant (other than the payer) now owes the payer their
        share.
        """
        try:
            shares = expense.strategy.compute(expense.amount, expense.participants)
        except InvalidSplitError as e:
            raise InvalidSplitError(f"invalid split: {e}") from e

        with self._lock:
            for p in expense.participants:
                if p.id == expense.paid_by.id:
                    continue
                self._adjust(expense.paid_by.id, p.id, shares[p.id])

    def net_balance(self, debtor: str, creditor: str) -> float:
        """Returns how much debtor owes creditor (may be negative if the
        balance actually runs the other way).
        """
        with self._lock:
            return self._balances.get(creditor, {}).get(debtor, 0.0)

    def net_balances(self) -> Dict[str, float]:
        """Returns each user's overall net position: positive means the
        user is a net creditor (is owed money), negative means net
        debtor.
        """
        with self._lock:
            return {a: sum(row.values()) for a, row in self._balances.items()}


@dataclass(frozen=True)
class Transaction:
    from_id: str
    to_id: str
    amount: float


def simplify_debts(net: Dict[str, float]) -> List[Transaction]:
    """Takes each user's net balance (positive = creditor, negative =
    debtor) and greedily matches the largest creditor with the largest
    debtor, minimizing the number of settling transactions. This is the
    classic greedy algorithm: it does not guarantee the theoretical
    minimum in every adversarial case, but it is optimal for the common
    case and is what real Splitwise-style systems use.
    """
    creditors = sorted(
        ([uid, amt] for uid, amt in net.items() if amt > EPSILON),
        key=lambda b: b[1],
        reverse=True,
    )
    debtors = sorted(
        ([uid, -amt] for uid, amt in net.items() if amt < -EPSILON),
        key=lambda b: b[1],
        reverse=True,
    )

    transactions: List[Transaction] = []
    i, j = 0, 0
    while i < len(creditors) and j < len(debtors):
        c, d = creditors[i], debtors[j]
        amount = min(c[1], d[1])

        transactions.append(Transaction(from_id=d[0], to_id=c[0], amount=amount))

        c[1] -= amount
        d[1] -= amount
        if c[1] <= EPSILON:
            i += 1
        if d[1] <= EPSILON:
            j += 1

    return transactions


class Splitwise:
    """Ties together users, groups and a ledger into the app-level API."""

    def __init__(self):
        self._lock = threading.Lock()
        self.users: Dict[str, User] = {}
        self.groups: Dict[str, Group] = {}
        self.ledger = Ledger()
        self._seq = 0

    def add_user(self, user_id: str, name: str) -> User:
        with self._lock:
            u = User(id=user_id, name=name)
            self.users[user_id] = u
            return u

    def add_group(self, group_id: str, name: str, members: List[User]) -> Group:
        with self._lock:
            g = Group(id=group_id, name=name, members=members)
            self.groups[group_id] = g
            return g

    def add_expense(
        self,
        description: str,
        paid_by: User,
        amount: float,
        participants: List[User],
        strategy: SplitStrategy,
    ) -> Expense:
        """Creates an expense with the given strategy and records it in
        the ledger. Returns the expense (with a generated ID) or raises
        InvalidSplitError if the split is invalid.
        """
        with self._lock:
            self._seq += 1
            expense = Expense(
                id=f"E-{self._seq}",
                description=description,
                paid_by=paid_by,
                amount=amount,
                participants=participants,
                strategy=strategy,
            )

        self.ledger.add_expense(expense)
        return expense


def _demo() -> None:
    sw = Splitwise()
    alice = sw.add_user("u1", "Alice")
    bob = sw.add_user("u2", "Bob")
    carol = sw.add_user("u3", "Carol")

    sw.add_expense("Dinner", alice, 90, [alice, bob, carol], EqualSplit())
    print(f"Bob owes Alice: {sw.ledger.net_balance(bob.id, alice.id):.2f}")
    print(f"Carol owes Alice: {sw.ledger.net_balance(carol.id, alice.id):.2f}")

    net = sw.ledger.net_balances()
    for txn in simplify_debts(net):
        print(f"{txn.from_id} pays {txn.to_id}: {txn.amount:.2f}")


if __name__ == "__main__":
    _demo()
