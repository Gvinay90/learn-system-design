import math

import pytest

from splitwise import (
    EqualSplit,
    ExactSplit,
    PercentSplit,
    Splitwise,
    simplify_debts,
)


def approx_equal(a: float, b: float) -> bool:
    return math.isclose(a, b, abs_tol=1e-6)


def test_equal_split_creates_correct_balances():
    sw = Splitwise()
    alice = sw.add_user("u1", "Alice")
    bob = sw.add_user("u2", "Bob")
    carol = sw.add_user("u3", "Carol")

    sw.add_expense("Dinner", alice, 90, [alice, bob, carol], EqualSplit())

    assert approx_equal(sw.ledger.net_balance(bob.id, alice.id), 30)
    assert approx_equal(sw.ledger.net_balance(carol.id, alice.id), 30)
    assert approx_equal(sw.ledger.net_balance(alice.id, bob.id), -30)


def test_exact_split_validation():
    sw = Splitwise()
    alice = sw.add_user("u1", "Alice")
    bob = sw.add_user("u2", "Bob")

    with pytest.raises(Exception):
        sw.add_expense(
            "Rent", alice, 100, [alice, bob],
            ExactSplit(amounts={alice.id: 40, bob.id: 50}),
        )

    sw.add_expense(
        "Rent", alice, 100, [alice, bob],
        ExactSplit(amounts={alice.id: 40, bob.id: 60}),
    )
    assert approx_equal(sw.ledger.net_balance(bob.id, alice.id), 60)


def test_percent_split_correctness():
    sw = Splitwise()
    alice = sw.add_user("u1", "Alice")
    bob = sw.add_user("u2", "Bob")
    carol = sw.add_user("u3", "Carol")

    sw.add_expense(
        "Trip", alice, 200, [alice, bob, carol],
        PercentSplit(percentages={alice.id: 50, bob.id: 25, carol.id: 25}),
    )

    assert approx_equal(sw.ledger.net_balance(bob.id, alice.id), 50)
    assert approx_equal(sw.ledger.net_balance(carol.id, alice.id), 50)

    with pytest.raises(Exception):
        sw.add_expense(
            "Bad", alice, 200, [alice, bob, carol],
            PercentSplit(percentages={alice.id: 50, bob.id: 25, carol.id: 20}),
        )


def test_simplify_debts():
    # A owes 300 net, B owes 100 net, C is owed 200, D is owed 200.
    net = {"A": -300, "B": -100, "C": 200, "D": 200}

    txns = simplify_debts(net)

    assert len(txns) <= 3

    settled: dict = {}
    for tx in txns:
        settled[tx.from_id] = settled.get(tx.from_id, 0.0) - tx.amount
        settled[tx.to_id] = settled.get(tx.to_id, 0.0) + tx.amount

    for user, expected in net.items():
        assert approx_equal(settled.get(user, 0.0), expected)


def test_simplify_debts_minimizes_transaction_count():
    # Three-way cycle: A owes B 10, B owes C 10, C owes A 10 -> net balances
    # are all zero, so zero transactions should be needed.
    net = {"A": 0, "B": 0, "C": 0}
    txns = simplify_debts(net)
    assert txns == []


def test_equal_split_rejects_no_participants():
    sw = Splitwise()
    alice = sw.add_user("u1", "Alice")
    with pytest.raises(Exception):
        sw.add_expense("Nothing", alice, 10, [], EqualSplit())
