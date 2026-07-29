public class CreatedState implements OrderState {
    @Override
    public String name() { return "Created"; }

    @Override
    public void pay(Order order) { order.setState(new PaidState()); }

    @Override
    public void ship(Order order) {
        throw new IllegalTransitionException("cannot ship an order that has not been paid");
    }

    @Override
    public void deliver(Order order) {
        throw new IllegalTransitionException("cannot deliver an order that has not shipped");
    }

    @Override
    public void cancel(Order order) { order.setState(new CancelledState()); }
}
