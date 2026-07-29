import java.util.List;

public class Main {
    public static void main(String[] args) throws Exception {
        OrderBook book = new OrderBook("AAPL");

        book.submitOrder(new Order("S1", "AAPL", Side.SELL, OrderType.LIMIT, 100, 10));
        List<Trade> trades = book.submitOrder(new Order("B1", "AAPL", Side.BUY, OrderType.LIMIT, 100, 10));

        System.out.println("Trades from crossing order: " + trades);
        System.out.println("Resting buy orders: " + book.buyOrders());
        System.out.println("Resting sell orders: " + book.sellOrders());
        System.out.println("All trades recorded: " + book.trades());

        OrderBookTest.runAll();
    }
}
