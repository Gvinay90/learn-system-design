/**
 * The context: holds a reference to its current state and forwards every
 * action to it.
 */
public class Order {
    private final String id;
    private OrderState state;

    public Order(String id) {
        this.id = id;
        this.state = new CreatedState();
    }

    public String getId() { return id; }
    public OrderState getState() { return state; }

    void setState(OrderState state) { this.state = state; }

    public void pay() { state.pay(this); }
    public void ship() { state.ship(this); }
    public void deliver() { state.deliver(this); }
    public void cancel() { state.cancel(this); }
}
