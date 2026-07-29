import pytest

from proxy import (
    BankAccountProxy,
    InsufficientFundsError,
    RealBankAccount,
    UnauthorizedError,
)


def test_owner_can_withdraw():
    acc = BankAccountProxy(RealBankAccount("alice", 100))
    acc.withdraw("alice", 40)
    assert acc.balance() == 60


def test_non_owner_cannot_withdraw():
    acc = BankAccountProxy(RealBankAccount("alice", 100))
    with pytest.raises(UnauthorizedError):
        acc.withdraw("mallory", 10)
    assert acc.balance() == 100


def test_non_owner_cannot_deposit():
    acc = BankAccountProxy(RealBankAccount("alice", 100))
    with pytest.raises(UnauthorizedError):
        acc.deposit("mallory", 500)
    assert acc.balance() == 100


def test_insufficient_funds():
    acc = BankAccountProxy(RealBankAccount("alice", 20))
    with pytest.raises(InsufficientFundsError):
        acc.withdraw("alice", 50)


def test_access_log_records_all_attempts():
    acc = BankAccountProxy(RealBankAccount("alice", 100))
    acc.withdraw("alice", 10)
    try:
        acc.deposit("mallory", 5)
    except UnauthorizedError:
        pass
    assert len(acc.access_log) == 2
