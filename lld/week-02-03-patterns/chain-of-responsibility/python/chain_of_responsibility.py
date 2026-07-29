"""Chain of Responsibility pattern — an expense request travels along a chain
of approvers, each approving within its limit or forwarding to the next.

See ../README.md for the design writeup.
"""
from __future__ import annotations

from dataclasses import dataclass
from typing import List, Optional


class NoApproverError(Exception):
    pass


@dataclass
class ExpenseRequest:
    amount: float
    description: str = ""


class Approver:
    def __init__(self, name: str, limit: float):
        self.name = name
        self.limit = limit
        self._next: Optional["Approver"] = None

    def set_next(self, nxt: "Approver") -> None:
        self._next = nxt

    def approve(self, req: ExpenseRequest) -> str:
        if req.amount <= self.limit:
            return self.name
        if self._next is not None:
            return self._next.approve(req)
        raise NoApproverError("no approver in the chain can approve this amount")


class Manager(Approver):
    def __init__(self, limit: float):
        super().__init__("Manager", limit)


class Director(Approver):
    def __init__(self, limit: float):
        super().__init__("Director", limit)


class VP(Approver):
    def __init__(self, limit: float):
        super().__init__("VP", limit)


def build_chain(approvers: List[Approver]) -> Approver:
    for i in range(len(approvers) - 1):
        approvers[i].set_next(approvers[i + 1])
    return approvers[0]


def _demo() -> None:
    chain = build_chain([Manager(1000), Director(5000), VP(20000)])
    who = chain.approve(ExpenseRequest(3000, "conference"))
    print(f"Approved by: {who}")


if __name__ == "__main__":
    _demo()
