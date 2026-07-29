import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out OrderBookTest` directly.
 */
public class OrderBookTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testRestingOrderNoCrossNoTrade();
        testExactMatchFillsBoth();
        testPartialFillLeavesRemainderResting();
        testPriceTimePriority();
        testMarketOrderMatchesAtRestingPriceAndDoesNotRest();
        testNonCrossingOrdersBothRest();
        testCancelOrder();
        testCancelledOrderCannotMatch();
        testConcurrentSubmissionsNoDoubleFill();
        testTradesRecorded();
        System.out.println("All OrderBookTest cases passed.");
    }

    private static Order newOrder(String id, Side side, OrderType type, double price, int qty) {
        return new Order(id, "AAPL", side, type, price, qty);
    }

    private static void testRestingOrderNoCrossNoTrade() {
        OrderBook ob = new OrderBook("AAPL");
        List<Trade> trades = ob.submitOrder(newOrder("B1", Side.BUY, OrderType.LIMIT, 100, 10));
        assertEquals(0, trades.size(), "expected no trades");

        List<Order> buys = ob.buyOrders();
        assertEquals(1, buys.size(), "expected B1 resting on the buy side");
        assertEquals("B1", buys.get(0).id, "expected B1 resting on the buy side");
        assertEquals(OrderStatus.OPEN, buys.get(0).status, "expected resting order status OPEN");
    }

    private static void testExactMatchFillsBoth() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("S1", Side.SELL, OrderType.LIMIT, 100, 10));
        List<Trade> trades = ob.submitOrder(newOrder("B1", Side.BUY, OrderType.LIMIT, 100, 10));

        assertEquals(1, trades.size(), "expected 1 trade");
        Trade tr = trades.get(0);
        assertTrue(tr.price == 100 && tr.quantity == 10 && tr.buyOrderId.equals("B1") && tr.sellOrderId.equals("S1"),
                "unexpected trade: " + tr);
        assertTrue(ob.buyOrders().isEmpty() && ob.sellOrders().isEmpty(), "expected both sides empty after full match");
    }

    private static void testPartialFillLeavesRemainderResting() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("S1", Side.SELL, OrderType.LIMIT, 100, 5));
        List<Trade> trades = ob.submitOrder(newOrder("B1", Side.BUY, OrderType.LIMIT, 100, 10));

        assertTrue(trades.size() == 1 && trades.get(0).quantity == 5, "expected 1 trade of qty 5, got " + trades);
        List<Order> buys = ob.buyOrders();
        assertTrue(buys.size() == 1 && buys.get(0).id.equals("B1"), "expected B1 resting with remainder, got " + buys);
        assertEquals(5, buys.get(0).remainingQty, "expected remaining qty 5");
        assertEquals(OrderStatus.PARTIALLY_FILLED, buys.get(0).status, "expected PARTIALLY_FILLED");
        assertTrue(ob.sellOrders().isEmpty(), "expected sell side fully consumed");
    }

    // Asserts that among resting sell orders at different prices, the incoming
    // buy matches the best (lowest) price first, and among equal prices, the
    // earliest-submitted order matches first (time priority).
    private static void testPriceTimePriority() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("S-expensive", Side.SELL, OrderType.LIMIT, 105, 10));
        ob.submitOrder(newOrder("S-cheap", Side.SELL, OrderType.LIMIT, 100, 10));
        ob.submitOrder(newOrder("S-cheap-2", Side.SELL, OrderType.LIMIT, 100, 10));

        List<Trade> trades = ob.submitOrder(newOrder("B1", Side.BUY, OrderType.LIMIT, 105, 10));
        assertEquals(1, trades.size(), "expected 1 trade");
        assertEquals("S-cheap", trades.get(0).sellOrderId, "expected best price (S-cheap) to match first");

        List<Trade> trades2 = ob.submitOrder(newOrder("B2", Side.BUY, OrderType.LIMIT, 105, 10));
        assertTrue(trades2.size() == 1 && trades2.get(0).sellOrderId.equals("S-cheap-2"),
                "expected time priority to match S-cheap-2 next, got " + trades2);
    }

    private static void testMarketOrderMatchesAtRestingPriceAndDoesNotRest() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("S1", Side.SELL, OrderType.LIMIT, 100, 5));

        List<Trade> trades = ob.submitOrder(newOrder("M1", Side.BUY, OrderType.MARKET, 0, 10));
        assertTrue(trades.size() == 1 && trades.get(0).quantity == 5 && trades.get(0).price == 100,
                "expected market order to match available 5 @ 100, got " + trades);
        // Unfilled remainder of a market order must not rest in the book.
        assertTrue(ob.buyOrders().isEmpty(), "expected market order remainder to be cancelled, not resting, got " + ob.buyOrders());
    }

    private static void testNonCrossingOrdersBothRest() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("B1", Side.BUY, OrderType.LIMIT, 90, 10));
        List<Trade> trades = ob.submitOrder(newOrder("S1", Side.SELL, OrderType.LIMIT, 100, 10));

        assertEquals(0, trades.size(), "expected no trades for non-crossing prices");
        assertTrue(ob.buyOrders().size() == 1 && ob.sellOrders().size() == 1, "expected both orders resting untouched");
    }

    private static void testCancelOrder() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("B1", Side.BUY, OrderType.LIMIT, 100, 10));

        try {
            ob.cancelOrder("B1");
        } catch (OrderNotFoundException e) {
            throw new AssertionError("unexpected err cancelling: " + e);
        }
        assertTrue(ob.buyOrders().isEmpty(), "expected order removed from book after cancel");

        boolean threw = false;
        try {
            ob.cancelOrder("B1");
        } catch (OrderNotFoundException e) {
            threw = true;
        }
        assertTrue(threw, "expected OrderNotFoundException cancelling twice");

        threw = false;
        try {
            ob.cancelOrder("bogus");
        } catch (OrderNotFoundException e) {
            threw = true;
        }
        assertTrue(threw, "expected OrderNotFoundException for unknown order");
    }

    private static void testCancelledOrderCannotMatch() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("B1", Side.BUY, OrderType.LIMIT, 100, 10));
        try {
            ob.cancelOrder("B1");
        } catch (OrderNotFoundException e) {
            throw new AssertionError("unexpected err: " + e);
        }

        List<Trade> trades = ob.submitOrder(newOrder("S1", Side.SELL, OrderType.LIMIT, 100, 10));
        assertEquals(0, trades.size(), "expected cancelled order to not participate in matching");
    }

    // Asserts many threads submitting marketable orders concurrently against a
    // single resting order never double-fill it — total matched quantity
    // across all trades must not exceed what was actually available.
    private static void testConcurrentSubmissionsNoDoubleFill() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("S1", Side.SELL, OrderType.LIMIT, 100, 100));

        final int n = 20;
        AtomicInteger totalMatched = new AtomicInteger(0);
        CountDownLatch latch = new CountDownLatch(n);
        for (int i = 0; i < n; i++) {
            final int idx = i;
            Thread t = new Thread(() -> {
                try {
                    List<Trade> trades = ob.submitOrder(newOrder("B-" + idx, Side.BUY, OrderType.LIMIT, 100, 10));
                    int matched = 0;
                    for (Trade tr : trades) {
                        matched += tr.quantity;
                    }
                    totalMatched.addAndGet(matched);
                } finally {
                    latch.countDown();
                }
            });
            t.start();
        }
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        assertEquals(100, totalMatched.get(), "expected exactly 100 total matched quantity (no double-fill)");
        assertTrue(ob.sellOrders().isEmpty(), "expected sell order fully consumed");
    }

    private static void testTradesRecorded() {
        OrderBook ob = new OrderBook("AAPL");
        ob.submitOrder(newOrder("S1", Side.SELL, OrderType.LIMIT, 100, 10));
        ob.submitOrder(newOrder("B1", Side.BUY, OrderType.LIMIT, 100, 10));

        List<Trade> trades = ob.trades();
        assertEquals(1, trades.size(), "expected 1 recorded trade");
    }

    private static void assertTrue(boolean condition, String label) {
        if (!condition) {
            throw new AssertionError(label);
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
