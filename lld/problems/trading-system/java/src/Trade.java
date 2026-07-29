public class Trade {
    public final String id;
    public final String symbol;
    public final String buyOrderId;
    public final String sellOrderId;
    public final double price;
    public final int quantity;
    public final long timestamp;

    public Trade(String id, String symbol, String buyOrderId, String sellOrderId, double price, int quantity) {
        this.id = id;
        this.symbol = symbol;
        this.buyOrderId = buyOrderId;
        this.sellOrderId = sellOrderId;
        this.price = price;
        this.quantity = quantity;
        this.timestamp = System.nanoTime();
    }

    @Override
    public String toString() {
        return "Trade{id=" + id + ", buy=" + buyOrderId + ", sell=" + sellOrderId
                + ", price=" + price + ", qty=" + quantity + "}";
    }
}
