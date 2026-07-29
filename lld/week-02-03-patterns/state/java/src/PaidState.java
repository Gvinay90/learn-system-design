public class PaidState implements OrderState {
    @Override
    public String name() { return "Paid"; }

    @Override
    public void pay(Order order) {
        throw new IllegalTransitionException("order already paid");
    }

    @Override
    public void ship(Order order) { order.setState(new ShippedState()); }

    @Override
    public void deliver(Order order) {
        throw new IllegalTransitionException("cannot deliver an order that has not shipped");
    }

    @Override
    public void cancel(Order order) { order.setState(new CancelledState()); }
}
