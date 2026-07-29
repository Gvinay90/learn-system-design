"""Trading system LLD — Python reference implementation.

A per-symbol order book that matches buy/sell limit (and market) orders
under price-time priority, supports partial fills and cancellation. See
../go/tradingsystem.go for the original design writeup.
"""
from __future__ import annotations

import itertools
import threading
import time
from dataclasses import dataclass, field
from enum import Enum, auto
from typing import Callable, List


class Side(Enum):
    BUY = auto()
    SELL = auto()


class OrderType(Enum):
    LIMIT = auto()
    MARKET = auto()


class OrderStatus(Enum):
    OPEN = auto()
    PARTIALLY_FILLED = auto()
    FILLED = auto()
    CANCELLED = auto()


class OrderNotFoundError(Exception):
    """Raised when cancelling an order that doesn't exist or already
    filled/cancelled."""

    def __init__(self) -> None:
        super().__init__("order not found or already filled/cancelled")


@dataclass
class Order:
    id: str
    symbol: str
    side: Side
    type: OrderType
    price: float
    quantity: int
    remaining_qty: int = field(init=False)
    timestamp: float = field(default_factory=time.time)
    status: OrderStatus = field(default=OrderStatus.OPEN)
    seq: int = field(default=0, repr=False)

    def __post_init__(self) -> None:
        self.remaining_qty = self.quantity


@dataclass
class Trade:
    id: str
    symbol: str
    buy_order_id: str
    sell_order_id: str
    price: float
    quantity: int
    timestamp: float = field(default_factory=time.time)


class OrderBook:
    """Holds resting orders for a single symbol, sorted by price-time
    priority: buy orders descending by price then ascending by arrival
    sequence, sell orders ascending by price then ascending by arrival
    sequence.

    The whole book is locked for the duration of matching: concurrent
    submissions to the same book must be serialized, otherwise two threads
    could both match against the same resting order and double-fill it.
    """

    def __init__(self, symbol: str):
        self.symbol = symbol
        self._lock = threading.Lock()
        self._buy_orders: List[Order] = []
        self._sell_orders: List[Order] = []
        self._orders: dict[str, Order] = {}
        self._trades: List[Trade] = []
        self._seq = itertools.count(1)
        self._trade_seq = itertools.count(1)

    def submit_order(self, order: Order) -> List[Trade]:
        """Matches the incoming order against the resting opposite side
        while prices cross, returns the resulting trades, and rests any
        unfilled remainder in the book (LIMIT orders only; MARKET orders
        never rest — an unfilled remainder is simply cancelled)."""
        with self._lock:
            order.seq = next(self._seq)
            order.remaining_qty = order.quantity
            order.status = OrderStatus.OPEN
            self._orders[order.id] = order

            if order.side == Side.BUY:
                trades = self._match_incoming(
                    order,
                    self._sell_orders,
                    lambda resting_price: order.type == OrderType.MARKET or order.price >= resting_price,
                )
            else:
                trades = self._match_incoming(
                    order,
                    self._buy_orders,
                    lambda resting_price: order.type == OrderType.MARKET or order.price <= resting_price,
                )

            if order.remaining_qty == 0:
                order.status = OrderStatus.FILLED
                self._orders.pop(order.id, None)
                return trades

            if order.type == OrderType.MARKET:
                order.status = OrderStatus.CANCELLED
                self._orders.pop(order.id, None)
                return trades

            if trades:
                order.status = OrderStatus.PARTIALLY_FILLED

            if order.side == Side.BUY:
                self._insert_sorted(self._buy_orders, order, buy=True)
            else:
                self._insert_sorted(self._sell_orders, order, buy=False)
            return trades

    def _match_incoming(
        self,
        order: Order,
        resting_side: List[Order],
        crosses: Callable[[float], bool],
    ) -> List[Trade]:
        """Repeatedly crosses order against the front of resting_side while
        crosses(price) holds, producing trades and shrinking/removing
        resting orders as they get filled."""
        trades: List[Trade] = []
        while order.remaining_qty > 0 and resting_side:
            resting = resting_side[0]
            if not crosses(resting.price):
                break

            qty = min(order.remaining_qty, resting.remaining_qty)

            if order.side == Side.BUY:
                buy_order_id, sell_order_id = order.id, resting.id
            else:
                buy_order_id, sell_order_id = resting.id, order.id

            trade = Trade(
                id=f"TR-{next(self._trade_seq)}",
                symbol=self.symbol,
                buy_order_id=buy_order_id,
                sell_order_id=sell_order_id,
                price=resting.price,
                quantity=qty,
            )
            trades.append(trade)
            self._trades.append(trade)

            order.remaining_qty -= qty
            resting.remaining_qty -= qty

            if resting.remaining_qty == 0:
                resting.status = OrderStatus.FILLED
                self._orders.pop(resting.id, None)
                resting_side.pop(0)
            else:
                resting.status = OrderStatus.PARTIALLY_FILLED

        return trades

    def cancel_order(self, order_id: str) -> None:
        """Removes a still-resting order from the book. Raises
        OrderNotFoundError if the order does not exist or has already been
        filled/cancelled."""
        with self._lock:
            order = self._orders.get(order_id)
            if order is None:
                raise OrderNotFoundError()
            order.status = OrderStatus.CANCELLED
            del self._orders[order_id]

            side = self._buy_orders if order.side == Side.BUY else self._sell_orders
            for i, o in enumerate(side):
                if o.id == order_id:
                    del side[i]
                    break

    def buy_orders(self) -> List[Order]:
        """Snapshot of resting buy orders, best priority first."""
        with self._lock:
            return list(self._buy_orders)

    def sell_orders(self) -> List[Order]:
        """Snapshot of resting sell orders, best priority first."""
        with self._lock:
            return list(self._sell_orders)

    def trades(self) -> List[Trade]:
        with self._lock:
            return list(self._trades)

    @staticmethod
    def _insert_sorted(orders: List[Order], order: Order, buy: bool) -> None:
        # Buy side: best price first = highest price, then earliest seq.
        # Sell side: best price first = lowest price, then earliest seq.
        price_key = -order.price if buy else order.price
        i = 0
        while i < len(orders):
            existing_key = -orders[i].price if buy else orders[i].price
            if (existing_key, orders[i].seq) <= (price_key, order.seq):
                i += 1
            else:
                break
        orders.insert(i, order)


def _demo() -> None:
    book = OrderBook("AAPL")
    book.submit_order(Order("S1", "AAPL", Side.SELL, OrderType.LIMIT, 100, 10))
    trades = book.submit_order(Order("B1", "AAPL", Side.BUY, OrderType.LIMIT, 100, 10))

    print("Trades:", trades)
    print("Resting buys:", book.buy_orders())
    print("Resting sells:", book.sell_orders())
    print("All trades:", book.trades())


if __name__ == "__main__":
    _demo()
