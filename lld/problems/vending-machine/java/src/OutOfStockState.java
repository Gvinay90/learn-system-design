/**
 * The last selectItem named a slot with zero quantity. The machine holds
 * no money in this state (nothing was accepted for an out-of-stock item),
 * so cancel simply returns to Idle; selectItem may be retried with a
 * different (in-stock) slot.
 */
public class OutOfStockState implements VendingState {
    @Override
    public String name() { return "OutOfStock"; }

    @Override
    public void selectItem(VendingMachine vm, String slotId) {
        StateHelpers.selectItem(vm, slotId);
    }

    @Override
    public void insertMoney(VendingMachine vm, int amount) {
        throw new InvalidStateException("selected item is out of stock");
    }

    @Override
    public DispenseResult dispense(VendingMachine vm) {
        throw new InvalidStateException("selected item is out of stock");
    }

    @Override
    public int cancel(VendingMachine vm) {
        int refund = vm.getBalance();
        vm.setBalance(0);
        vm.setSelected("");
        vm.setState(new IdleState());
        return refund;
    }
}
