public class CancelledState implements OrderState {
    @Override
    public String name() { return "Cancelled"; }

    @Override
    public void pay(Order order) {
        throw new IllegalTransitionException("order is cancelled");
    }

    @Override
    public void ship(Order order) {
        throw new IllegalTransitionException("order is cancelled");
    }

    @Override
    public void deliver(Order order) {
        throw new IllegalTransitionException("order is cancelled");
    }

    @Override
    public void cancel(Order order) {
        throw new IllegalTransitionException("order is already cancelled");
    }
}
