"""Proxy pattern — BankAccountProxy controls access to a RealBankAccount.

See ../README.md for the design writeup.
"""
from __future__ import annotations

from typing import List, Protocol


class UnauthorizedError(Exception):
    pass


class InsufficientFundsError(Exception):
    pass


class Account(Protocol):
    def deposit(self, requester: str, amount: float) -> None: ...
    def withdraw(self, requester: str, amount: float) -> None: ...
    def balance(self) -> float: ...


class RealBankAccount:
    """The expensive/sensitive real subject. Performs no authorization
    checks of its own — that's the proxy's job."""

    def __init__(self, owner: str, initial_balance: float):
        self.owner = owner
        self._balance = initial_balance

    def deposit(self, requester: str, amount: float) -> None:
        self._balance += amount

    def withdraw(self, requester: str, amount: float) -> None:
        if amount > self._balance:
            raise InsufficientFundsError("insufficient funds")
        self._balance -= amount

    def balance(self) -> float:
        return self._balance


class BankAccountProxy:
    """Protection proxy: enforces that only the account owner may deposit or
    withdraw, and logs every access attempt."""

    def __init__(self, real: RealBankAccount):
        self._real = real
        self.access_log: List[str] = []

    def deposit(self, requester: str, amount: float) -> None:
        self.access_log.append(f"deposit attempt by {requester}")
        if requester != self._real.owner:
            raise UnauthorizedError("requester is not the account owner")
        self._real.deposit(requester, amount)

    def withdraw(self, requester: str, amount: float) -> None:
        self.access_log.append(f"withdraw attempt by {requester}")
        if requester != self._real.owner:
            raise UnauthorizedError("requester is not the account owner")
        self._real.withdraw(requester, amount)

    def balance(self) -> float:
        return self._real.balance()


def _demo() -> None:
    real = RealBankAccount("alice", 100)
    acc = BankAccountProxy(real)

    acc.withdraw("alice", 40)
    print(f"Balance after owner withdraw: {acc.balance()}")

    try:
        acc.withdraw("mallory", 10)
    except UnauthorizedError as e:
        print(f"Blocked: {e}")


if __name__ == "__main__":
    _demo()
