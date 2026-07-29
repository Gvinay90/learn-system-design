/**
 * Every concrete state implements this. Each method mutates the machine to
 * the next state on success, or throws a specific unchecked exception
 * (without changing state) if the action isn't legal from that state.
 */
public interface VendingState {
    class InvalidSlotException extends RuntimeException {
        public InvalidSlotException(String message) { super(message); }
    }

    class OutOfStockException extends RuntimeException {
        public OutOfStockException(String message) { super(message); }
    }

    class NoSelectionException extends RuntimeException {
        public NoSelectionException(String message) { super(message); }
    }

    class InvalidAmountException extends RuntimeException {
        public InvalidAmountException(String message) { super(message); }
    }

    class NotEnoughMoneyException extends RuntimeException {
        public NotEnoughMoneyException(String message) { super(message); }
    }

    class InvalidStateException extends RuntimeException {
        public InvalidStateException(String message) { super(message); }
    }

    String name();
    void selectItem(VendingMachine vm, String slotId);
    void insertMoney(VendingMachine vm, int amount);
    DispenseResult dispense(VendingMachine vm);
    int cancel(VendingMachine vm);
}
