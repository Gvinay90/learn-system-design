public class ShippedState implements OrderState {
    @Override
    public String name() { return "Shipped"; }

    @Override
    public void pay(Order order) {
        throw new IllegalTransitionException("order already paid");
    }

    @Override
    public void ship(Order order) {
        throw new IllegalTransitionException("order already shipped");
    }

    @Override
    public void deliver(Order order) { order.setState(new DeliveredState()); }

    @Override
    public void cancel(Order order) {
        throw new IllegalTransitionException("cannot cancel an order that has already shipped");
    }
}
