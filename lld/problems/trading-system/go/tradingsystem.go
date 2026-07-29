// Package tradingsystem implements a simplified order matching engine LLD
// problem: a per-symbol order book that matches buy/sell limit orders under
// price-time priority, supports partial fills and cancellation.
package tradingsystem

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Side int

const (
	Buy Side = iota
	Sell
)

type OrderType int

const (
	Limit OrderType = iota
	Market
)

type OrderStatus int

const (
	Open OrderStatus = iota
	PartiallyFilled
	Filled
	Cancelled
)

type Order struct {
	ID            string
	Symbol        string
	Side          Side
	Type          OrderType
	Price         float64
	Quantity      int
	RemainingQty  int
	Timestamp     time.Time
	Status        OrderStatus
	seq           int64
}

type Trade struct {
	ID          string
	Symbol      string
	BuyOrderID  string
	SellOrderID string
	Price       float64
	Quantity    int
	Timestamp   time.Time
}

var ErrOrderNotFound = errors.New("order not found or already filled/cancelled")

// OrderBook holds resting orders for a single symbol, sorted by price-time
// priority: buyOrders descending by price then ascending by seq (arrival
// order), sellOrders ascending by price then ascending by seq.
type OrderBook struct {
	Symbol string

	mu         sync.Mutex
	buyOrders  []*Order
	sellOrders []*Order
	orders     map[string]*Order
	trades     []Trade
	seq        int64
	tradeSeq   int64
}

func NewOrderBook(symbol string) *OrderBook {
	return &OrderBook{
		Symbol: symbol,
		orders: make(map[string]*Order),
	}
}

// SubmitOrder matches the incoming order against the resting opposite side
// while prices cross, returns the resulting trades, and rests any
// unfilled remainder in the book (LIMIT orders only; MARKET orders never
// rest — an unfilled remainder is simply cancelled).
//
// The whole book is locked for the duration of matching: concurrent
// submissions to the same book must be serialized, otherwise two
// goroutines could both match against the same resting order and
// double-fill it.
func (ob *OrderBook) SubmitOrder(order *Order) []Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.seq++
	order.seq = ob.seq
	order.RemainingQty = order.Quantity
	order.Status = Open
	ob.orders[order.ID] = order

	var trades []Trade
	if order.Side == Buy {
		trades = ob.matchIncoming(order, &ob.sellOrders, func(restingPrice float64) bool {
			return order.Type == Market || order.Price >= restingPrice
		})
	} else {
		trades = ob.matchIncoming(order, &ob.buyOrders, func(restingPrice float64) bool {
			return order.Type == Market || order.Price <= restingPrice
		})
	}

	if order.RemainingQty == 0 {
		order.Status = Filled
		delete(ob.orders, order.ID)
		return trades
	}

	if order.Type == Market {
		order.Status = Cancelled
		delete(ob.orders, order.ID)
		return trades
	}

	if len(trades) > 0 {
		order.Status = PartiallyFilled
	}
	if order.Side == Buy {
		ob.buyOrders = insertSorted(ob.buyOrders, order, buyLess)
	} else {
		ob.sellOrders = insertSorted(ob.sellOrders, order, sellLess)
	}
	return trades
}

// matchIncoming repeatedly crosses order against the front of restingSide
// while crosses(price) holds, producing trades and shrinking/removing
// resting orders as they get filled.
func (ob *OrderBook) matchIncoming(order *Order, restingSide *[]*Order, crosses func(restingPrice float64) bool) []Trade {
	var trades []Trade
	for order.RemainingQty > 0 && len(*restingSide) > 0 {
		resting := (*restingSide)[0]
		if !crosses(resting.Price) {
			break
		}

		qty := order.RemainingQty
		if resting.RemainingQty < qty {
			qty = resting.RemainingQty
		}

		ob.tradeSeq++
		trade := Trade{
			Symbol:    ob.Symbol,
			Price:     resting.Price,
			Quantity:  qty,
			Timestamp: time.Now(),
		}
		if order.Side == Buy {
			trade.BuyOrderID, trade.SellOrderID = order.ID, resting.ID
		} else {
			trade.BuyOrderID, trade.SellOrderID = resting.ID, order.ID
		}
		trade.ID = fmt.Sprintf("TR-%d", ob.tradeSeq)
		trades = append(trades, trade)
		ob.trades = append(ob.trades, trade)

		order.RemainingQty -= qty
		resting.RemainingQty -= qty

		if resting.RemainingQty == 0 {
			resting.Status = Filled
			delete(ob.orders, resting.ID)
			*restingSide = (*restingSide)[1:]
		} else {
			resting.Status = PartiallyFilled
		}
	}
	return trades
}

// CancelOrder removes a still-resting order from the book. It fails if the
// order does not exist or has already been filled/cancelled.
func (ob *OrderBook) CancelOrder(orderID string) error {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	order, ok := ob.orders[orderID]
	if !ok {
		return ErrOrderNotFound
	}
	order.Status = Cancelled
	delete(ob.orders, orderID)

	if order.Side == Buy {
		ob.buyOrders = removeByID(ob.buyOrders, orderID)
	} else {
		ob.sellOrders = removeByID(ob.sellOrders, orderID)
	}
	return nil
}

// BuyOrders and SellOrders return snapshots of the resting book, best
// priority first, for inspection/testing.
func (ob *OrderBook) BuyOrders() []*Order {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	out := make([]*Order, len(ob.buyOrders))
	copy(out, ob.buyOrders)
	return out
}

func (ob *OrderBook) SellOrders() []*Order {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	out := make([]*Order, len(ob.sellOrders))
	copy(out, ob.sellOrders)
	return out
}

func (ob *OrderBook) Trades() []Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	out := make([]Trade, len(ob.trades))
	copy(out, ob.trades)
	return out
}

func buyLess(a, b *Order) bool {
	if a.Price != b.Price {
		return a.Price > b.Price
	}
	return a.seq < b.seq
}

func sellLess(a, b *Order) bool {
	if a.Price != b.Price {
		return a.Price < b.Price
	}
	return a.seq < b.seq
}

// insertSorted inserts o into the already-sorted slice, keeping the
// less(a,b) ordering. Linear scan is fine at interview scope; a real
// exchange would index price levels with a balanced tree or heap.
func insertSorted(orders []*Order, o *Order, less func(a, b *Order) bool) []*Order {
	i := 0
	for i < len(orders) && less(orders[i], o) {
		i++
	}
	orders = append(orders, nil)
	copy(orders[i+1:], orders[i:])
	orders[i] = o
	return orders
}

func removeByID(orders []*Order, id string) []*Order {
	for i, o := range orders {
		if o.ID == id {
			return append(orders[:i], orders[i+1:]...)
		}
	}
	return orders
}
