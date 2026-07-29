public class DeliveredState implements OrderState {
    @Override
    public String name() { return "Delivered"; }

    @Override
    public void pay(Order order) {
        throw new IllegalTransitionException("order already delivered");
    }

    @Override
    public void ship(Order order) {
        throw new IllegalTransitionException("order already delivered");
    }

    @Override
    public void deliver(Order order) {
        throw new IllegalTransitionException("order already delivered");
    }

    @Override
    public void cancel(Order order) {
        throw new IllegalTransitionException("cannot cancel a delivered order");
    }
}
