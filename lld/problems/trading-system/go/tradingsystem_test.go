package tradingsystem

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func newOrder(id string, side Side, orderType OrderType, price float64, qty int) *Order {
	return &Order{
		ID:        id,
		Symbol:    "AAPL",
		Side:      side,
		Type:      orderType,
		Price:     price,
		Quantity:  qty,
		Timestamp: time.Now(),
	}
}

func TestRestingOrderNoCrossNoTrade(t *testing.T) {
	ob := NewOrderBook("AAPL")
	trades := ob.SubmitOrder(newOrder("B1", Buy, Limit, 100, 10))
	if len(trades) != 0 {
		t.Fatalf("expected no trades, got %d", len(trades))
	}
	buys := ob.BuyOrders()
	if len(buys) != 1 || buys[0].ID != "B1" {
		t.Fatalf("expected B1 resting on the buy side, got %+v", buys)
	}
	if buys[0].Status != Open {
		t.Fatalf("expected resting order status Open, got %v", buys[0].Status)
	}
}

func TestExactMatchFillsBoth(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("S1", Sell, Limit, 100, 10))
	trades := ob.SubmitOrder(newOrder("B1", Buy, Limit, 100, 10))

	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}
	tr := trades[0]
	if tr.Price != 100 || tr.Quantity != 10 || tr.BuyOrderID != "B1" || tr.SellOrderID != "S1" {
		t.Fatalf("unexpected trade: %+v", tr)
	}
	if len(ob.BuyOrders()) != 0 || len(ob.SellOrders()) != 0 {
		t.Fatalf("expected both sides empty after full match")
	}
}

func TestPartialFillLeavesRemainderResting(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("S1", Sell, Limit, 100, 5))
	trades := ob.SubmitOrder(newOrder("B1", Buy, Limit, 100, 10))

	if len(trades) != 1 || trades[0].Quantity != 5 {
		t.Fatalf("expected 1 trade of qty 5, got %+v", trades)
	}
	buys := ob.BuyOrders()
	if len(buys) != 1 || buys[0].ID != "B1" {
		t.Fatalf("expected B1 resting with remainder, got %+v", buys)
	}
	if buys[0].RemainingQty != 5 {
		t.Fatalf("expected remaining qty 5, got %d", buys[0].RemainingQty)
	}
	if buys[0].Status != PartiallyFilled {
		t.Fatalf("expected PartiallyFilled, got %v", buys[0].Status)
	}
	if len(ob.SellOrders()) != 0 {
		t.Fatalf("expected sell side fully consumed")
	}
}

// TestPriceTimePriority asserts that among resting sell orders at different
// prices, the incoming buy matches the best (lowest) price first, and among
// equal prices, the earliest-submitted order matches first (time priority).
func TestPriceTimePriority(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("S-expensive", Sell, Limit, 105, 10))
	ob.SubmitOrder(newOrder("S-cheap", Sell, Limit, 100, 10))
	ob.SubmitOrder(newOrder("S-cheap-2", Sell, Limit, 100, 10))

	trades := ob.SubmitOrder(newOrder("B1", Buy, Limit, 105, 10))
	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}
	if trades[0].SellOrderID != "S-cheap" {
		t.Fatalf("expected best price (S-cheap) to match first, got %s", trades[0].SellOrderID)
	}

	trades2 := ob.SubmitOrder(newOrder("B2", Buy, Limit, 105, 10))
	if len(trades2) != 1 || trades2[0].SellOrderID != "S-cheap-2" {
		t.Fatalf("expected time priority to match S-cheap-2 next, got %+v", trades2)
	}
}

func TestMarketOrderMatchesAtRestingPriceAndDoesNotRest(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("S1", Sell, Limit, 100, 5))

	trades := ob.SubmitOrder(newOrder("M1", Buy, Market, 0, 10))
	if len(trades) != 1 || trades[0].Quantity != 5 || trades[0].Price != 100 {
		t.Fatalf("expected market order to match available 5 @ 100, got %+v", trades)
	}
	// Unfilled remainder of a market order must not rest in the book.
	if len(ob.BuyOrders()) != 0 {
		t.Fatalf("expected market order remainder to be cancelled, not resting, got %+v", ob.BuyOrders())
	}
}

func TestNonCrossingOrdersBothRest(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("B1", Buy, Limit, 90, 10))
	trades := ob.SubmitOrder(newOrder("S1", Sell, Limit, 100, 10))

	if len(trades) != 0 {
		t.Fatalf("expected no trades for non-crossing prices, got %d", len(trades))
	}
	if len(ob.BuyOrders()) != 1 || len(ob.SellOrders()) != 1 {
		t.Fatalf("expected both orders resting untouched")
	}
}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("B1", Buy, Limit, 100, 10))

	if err := ob.CancelOrder("B1"); err != nil {
		t.Fatalf("unexpected err cancelling: %v", err)
	}
	if len(ob.BuyOrders()) != 0 {
		t.Fatalf("expected order removed from book after cancel")
	}

	if err := ob.CancelOrder("B1"); err != ErrOrderNotFound {
		t.Fatalf("expected ErrOrderNotFound cancelling twice, got %v", err)
	}
	if err := ob.CancelOrder("bogus"); err != ErrOrderNotFound {
		t.Fatalf("expected ErrOrderNotFound for unknown order, got %v", err)
	}
}

func TestCancelledOrderCannotMatch(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("B1", Buy, Limit, 100, 10))
	if err := ob.CancelOrder("B1"); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	trades := ob.SubmitOrder(newOrder("S1", Sell, Limit, 100, 10))
	if len(trades) != 0 {
		t.Fatalf("expected cancelled order to not participate in matching, got %d trades", len(trades))
	}
}

// TestConcurrentSubmissionsNoDoubleFill asserts many goroutines submitting
// marketable orders concurrently against a single resting order never
// double-fill it — total matched quantity across all trades must not
// exceed what was actually available.
func TestConcurrentSubmissionsNoDoubleFill(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("S1", Sell, Limit, 100, 100))

	const n = 20
	var wg sync.WaitGroup
	tradesCh := make(chan []Trade, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("B-%d", i)
			trades := ob.SubmitOrder(newOrder(id, Buy, Limit, 100, 10))
			tradesCh <- trades
		}(i)
	}
	wg.Wait()
	close(tradesCh)

	totalMatched := 0
	for trades := range tradesCh {
		for _, tr := range trades {
			totalMatched += tr.Quantity
		}
	}
	if totalMatched != 100 {
		t.Fatalf("expected exactly 100 total matched quantity (no double-fill), got %d", totalMatched)
	}
	if len(ob.SellOrders()) != 0 {
		t.Fatalf("expected sell order fully consumed")
	}
}

func TestTradesRecorded(t *testing.T) {
	ob := NewOrderBook("AAPL")
	ob.SubmitOrder(newOrder("S1", Sell, Limit, 100, 10))
	ob.SubmitOrder(newOrder("B1", Buy, Limit, 100, 10))

	trades := ob.Trades()
	if len(trades) != 1 {
		t.Fatalf("expected 1 recorded trade, got %d", len(trades))
	}
}
