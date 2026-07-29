public class Order {
    public final String id;
    public final String symbol;
    public final Side side;
    public final OrderType type;
    public final double price;
    public final int quantity;
    public int remainingQty;
    public final long timestamp;
    public OrderStatus status;
    long seq;

    public Order(String id, String symbol, Side side, OrderType type, double price, int quantity) {
        this.id = id;
        this.symbol = symbol;
        this.side = side;
        this.type = type;
        this.price = price;
        this.quantity = quantity;
        this.remainingQty = quantity;
        this.timestamp = System.nanoTime();
        this.status = OrderStatus.OPEN;
    }

    @Override
    public String toString() {
        return "Order{id=" + id + ", side=" + side + ", type=" + type + ", price=" + price
                + ", qty=" + quantity + ", remaining=" + remainingQty + ", status=" + status + "}";
    }
}
