import threading

import pytest

from trading_system import (
    OrderBook,
    OrderNotFoundError,
    OrderStatus,
    OrderType,
    Side,
)
from trading_system import Order


def new_order(order_id, side, order_type, price, qty):
    return Order(order_id, "AAPL", side, order_type, price, qty)


def test_resting_order_no_cross_no_trade():
    ob = OrderBook("AAPL")
    trades = ob.submit_order(new_order("B1", Side.BUY, OrderType.LIMIT, 100, 10))
    assert len(trades) == 0

    buys = ob.buy_orders()
    assert len(buys) == 1
    assert buys[0].id == "B1"
    assert buys[0].status == OrderStatus.OPEN


def test_exact_match_fills_both():
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("S1", Side.SELL, OrderType.LIMIT, 100, 10))
    trades = ob.submit_order(new_order("B1", Side.BUY, OrderType.LIMIT, 100, 10))

    assert len(trades) == 1
    tr = trades[0]
    assert tr.price == 100 and tr.quantity == 10
    assert tr.buy_order_id == "B1" and tr.sell_order_id == "S1"
    assert ob.buy_orders() == []
    assert ob.sell_orders() == []


def test_partial_fill_leaves_remainder_resting():
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("S1", Side.SELL, OrderType.LIMIT, 100, 5))
    trades = ob.submit_order(new_order("B1", Side.BUY, OrderType.LIMIT, 100, 10))

    assert len(trades) == 1
    assert trades[0].quantity == 5

    buys = ob.buy_orders()
    assert len(buys) == 1
    assert buys[0].id == "B1"
    assert buys[0].remaining_qty == 5
    assert buys[0].status == OrderStatus.PARTIALLY_FILLED
    assert ob.sell_orders() == []


def test_price_time_priority():
    """Among resting sell orders at different prices, the incoming buy
    matches the best (lowest) price first, and among equal prices, the
    earliest-submitted order matches first (time priority)."""
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("S-expensive", Side.SELL, OrderType.LIMIT, 105, 10))
    ob.submit_order(new_order("S-cheap", Side.SELL, OrderType.LIMIT, 100, 10))
    ob.submit_order(new_order("S-cheap-2", Side.SELL, OrderType.LIMIT, 100, 10))

    trades = ob.submit_order(new_order("B1", Side.BUY, OrderType.LIMIT, 105, 10))
    assert len(trades) == 1
    assert trades[0].sell_order_id == "S-cheap"

    trades2 = ob.submit_order(new_order("B2", Side.BUY, OrderType.LIMIT, 105, 10))
    assert len(trades2) == 1
    assert trades2[0].sell_order_id == "S-cheap-2"


def test_market_order_matches_at_resting_price_and_does_not_rest():
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("S1", Side.SELL, OrderType.LIMIT, 100, 5))

    trades = ob.submit_order(new_order("M1", Side.BUY, OrderType.MARKET, 0, 10))
    assert len(trades) == 1
    assert trades[0].quantity == 5
    assert trades[0].price == 100

    # Unfilled remainder of a market order must not rest in the book.
    assert ob.buy_orders() == []


def test_non_crossing_orders_both_rest():
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("B1", Side.BUY, OrderType.LIMIT, 90, 10))
    trades = ob.submit_order(new_order("S1", Side.SELL, OrderType.LIMIT, 100, 10))

    assert len(trades) == 0
    assert len(ob.buy_orders()) == 1
    assert len(ob.sell_orders()) == 1


def test_cancel_order():
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("B1", Side.BUY, OrderType.LIMIT, 100, 10))

    ob.cancel_order("B1")
    assert ob.buy_orders() == []

    with pytest.raises(OrderNotFoundError):
        ob.cancel_order("B1")
    with pytest.raises(OrderNotFoundError):
        ob.cancel_order("bogus")


def test_cancelled_order_cannot_match():
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("B1", Side.BUY, OrderType.LIMIT, 100, 10))
    ob.cancel_order("B1")

    trades = ob.submit_order(new_order("S1", Side.SELL, OrderType.LIMIT, 100, 10))
    assert len(trades) == 0


def test_concurrent_submissions_no_double_fill():
    """Many threads submitting marketable orders concurrently against a
    single resting order must never double-fill it — total matched
    quantity across all trades must not exceed what was actually
    available."""
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("S1", Side.SELL, OrderType.LIMIT, 100, 100))

    n = 20
    results = [None] * n

    def worker(i):
        results[i] = ob.submit_order(new_order(f"B-{i}", Side.BUY, OrderType.LIMIT, 100, 10))

    threads = [threading.Thread(target=worker, args=(i,)) for i in range(n)]
    for t in threads:
        t.start()
    for t in threads:
        t.join()

    total_matched = sum(tr.quantity for trades in results for tr in trades)
    assert total_matched == 100
    assert ob.sell_orders() == []


def test_trades_recorded():
    ob = OrderBook("AAPL")
    ob.submit_order(new_order("S1", Side.SELL, OrderType.LIMIT, 100, 10))
    ob.submit_order(new_order("B1", Side.BUY, OrderType.LIMIT, 100, 10))

    trades = ob.trades()
    assert len(trades) == 1
