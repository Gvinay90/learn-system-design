import pytest

from chain_of_responsibility import (
    Director,
    ExpenseRequest,
    Manager,
    NoApproverError,
    VP,
    build_chain,
)


def new_test_chain():
    return build_chain([Manager(1000), Director(5000), VP(20000)])


def test_manager_approves_small_expense():
    chain = new_test_chain()
    assert chain.approve(ExpenseRequest(500, "team lunch")) == "Manager"


def test_director_approves_mid_expense():
    chain = new_test_chain()
    assert chain.approve(ExpenseRequest(3000, "conference")) == "Director"


def test_vp_approves_large_expense():
    chain = new_test_chain()
    assert chain.approve(ExpenseRequest(15000, "new hire equipment")) == "VP"


def test_no_approver_for_excessive_expense():
    chain = new_test_chain()
    with pytest.raises(NoApproverError):
        chain.approve(ExpenseRequest(100000, "private jet"))


def test_chain_order_matters():
    chain = build_chain([Manager(1000)])
    with pytest.raises(NoApproverError):
        chain.approve(ExpenseRequest(2000))
