/**
 * Shared selectItem logic used by both IdleState and OutOfStockState:
 * validate the slot, and either settle on it (moving/staying in Idle,
 * awaiting money) or bounce to OutOfStockState if it's empty.
 */
final class StateHelpers {
    private StateHelpers() {}

    static void selectItem(VendingMachine vm, String slotId) {
        Slot s = vm.slot(slotId);
        if (s.quantity <= 0) {
            vm.setSelected(slotId);
            vm.setState(new OutOfStockState());
            throw new VendingState.OutOfStockException("item out of stock: " + slotId);
        }
        vm.setSelected(slotId);
        vm.setState(new IdleState());
    }
}
