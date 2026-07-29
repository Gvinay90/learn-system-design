import java.util.ArrayList;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Holds resting orders for a single symbol, matched under price-time
 * priority: buyOrders sorted descending by price then ascending by
 * arrival sequence, sellOrders ascending by price then ascending by
 * arrival sequence.
 *
 * The whole book is locked (via `synchronized`) for the duration of
 * matching: concurrent submissions to the same book must be serialized,
 * otherwise two threads could both match against the same resting order
 * and double-fill it.
 */
public class OrderBook {
    public final String symbol;

    private final List<Order> buyOrders = new ArrayList<>();
    private final List<Order> sellOrders = new ArrayList<>();
    private final Map<String, Order> orders = new HashMap<>();
    private final List<Trade> trades = new ArrayList<>();
    private long seq = 0;
    private long tradeSeq = 0;

    private static final Comparator<Order> BUY_ORDER = Comparator
            .comparingDouble((Order o) -> o.price).reversed()
            .thenComparingLong(o -> o.seq);
    private static final Comparator<Order> SELL_ORDER = Comparator
            .comparingDouble((Order o) -> o.price)
            .thenComparingLong(o -> o.seq);

    public OrderBook(String symbol) {
        this.symbol = symbol;
    }

    /**
     * Matches the incoming order against the resting opposite side while
     * prices cross, returns the resulting trades, and rests any unfilled
     * remainder in the book (LIMIT orders only; MARKET orders never rest —
     * an unfilled remainder is simply cancelled).
     */
    public synchronized List<Trade> submitOrder(Order order) {
        seq++;
        order.seq = seq;
        order.remainingQty = order.quantity;
        order.status = OrderStatus.OPEN;
        orders.put(order.id, order);

        List<Trade> resultTrades;
        if (order.side == Side.BUY) {
            resultTrades = matchIncoming(order, sellOrders,
                    restingPrice -> order.type == OrderType.MARKET || order.price >= restingPrice);
        } else {
            resultTrades = matchIncoming(order, buyOrders,
                    restingPrice -> order.type == OrderType.MARKET || order.price <= restingPrice);
        }

        if (order.remainingQty == 0) {
            order.status = OrderStatus.FILLED;
            orders.remove(order.id);
            return resultTrades;
        }

        if (order.type == OrderType.MARKET) {
            order.status = OrderStatus.CANCELLED;
            orders.remove(order.id);
            return resultTrades;
        }

        if (!resultTrades.isEmpty()) {
            order.status = OrderStatus.PARTIALLY_FILLED;
        }
        if (order.side == Side.BUY) {
            insertSorted(buyOrders, order, BUY_ORDER);
        } else {
            insertSorted(sellOrders, order, SELL_ORDER);
        }
        return resultTrades;
    }

    private interface CrossPredicate {
        boolean crosses(double restingPrice);
    }

    /**
     * Repeatedly crosses order against the front of restingSide while
     * crosses(price) holds, producing trades and shrinking/removing
     * resting orders as they get filled.
     */
    private List<Trade> matchIncoming(Order order, List<Order> restingSide, CrossPredicate crosses) {
        List<Trade> resultTrades = new ArrayList<>();
        while (order.remainingQty > 0 && !restingSide.isEmpty()) {
            Order resting = restingSide.get(0);
            if (!crosses.crosses(resting.price)) {
                break;
            }

            int qty = Math.min(order.remainingQty, resting.remainingQty);

            tradeSeq++;
            String buyOrderId;
            String sellOrderId;
            if (order.side == Side.BUY) {
                buyOrderId = order.id;
                sellOrderId = resting.id;
            } else {
                buyOrderId = resting.id;
                sellOrderId = order.id;
            }
            Trade trade = new Trade("TR-" + tradeSeq, symbol, buyOrderId, sellOrderId, resting.price, qty);
            resultTrades.add(trade);
            trades.add(trade);

            order.remainingQty -= qty;
            resting.remainingQty -= qty;

            if (resting.remainingQty == 0) {
                resting.status = OrderStatus.FILLED;
                orders.remove(resting.id);
                restingSide.remove(0);
            } else {
                resting.status = OrderStatus.PARTIALLY_FILLED;
            }
        }
        return resultTrades;
    }

    /** Removes a still-resting order from the book. */
    public synchronized void cancelOrder(String orderId) throws OrderNotFoundException {
        Order order = orders.get(orderId);
        if (order == null) {
            throw new OrderNotFoundException();
        }
        order.status = OrderStatus.CANCELLED;
        orders.remove(orderId);

        if (order.side == Side.BUY) {
            removeById(buyOrders, orderId);
        } else {
            removeById(sellOrders, orderId);
        }
    }

    /** Snapshot of resting buy orders, best priority first. */
    public synchronized List<Order> buyOrders() {
        return new ArrayList<>(buyOrders);
    }

    /** Snapshot of resting sell orders, best priority first. */
    public synchronized List<Order> sellOrders() {
        return new ArrayList<>(sellOrders);
    }

    public synchronized List<Trade> trades() {
        return new ArrayList<>(trades);
    }

    private static void insertSorted(List<Order> orders, Order o, Comparator<Order> cmp) {
        int i = 0;
        while (i < orders.size() && cmp.compare(orders.get(i), o) < 0) {
            i++;
        }
        orders.add(i, o);
    }

    private static void removeById(List<Order> orders, String id) {
        for (int i = 0; i < orders.size(); i++) {
            if (orders.get(i).id.equals(id)) {
                orders.remove(i);
                return;
            }
        }
    }
}
