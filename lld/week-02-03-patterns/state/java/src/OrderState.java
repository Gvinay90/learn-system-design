/**
 * Every concrete state implements this. Each method mutates the order to the
 * next state on success, or throws IllegalTransitionException.
 */
public interface OrderState {
    class IllegalTransitionException extends RuntimeException {
        public IllegalTransitionException(String message) { super(message); }
    }

    String name();
    void pay(Order order);
    void ship(Order order);
    void deliver(Order order);
    void cancel(Order order);
}
