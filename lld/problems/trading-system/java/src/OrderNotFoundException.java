public class OrderNotFoundException extends Exception {
    public OrderNotFoundException() {
        super("order not found or already filled/cancelled");
    }
}
