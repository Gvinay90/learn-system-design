/**
 * Enough money has been inserted to cover the selected item's price.
 * Further top-ups are still accepted; Dispense completes the purchase, and
 * Cancel refunds everything inserted.
 */
public class HasMoneyState implements VendingState {
    @Override
    public String name() { return "HasMoney"; }

    @Override
    public void selectItem(VendingMachine vm, String slotId) {
        throw new InvalidStateException("cannot change selection while funds are pending");
    }

    @Override
    public void insertMoney(VendingMachine vm, int amount) {
        if (amount <= 0) {
            throw new InvalidAmountException("amount must be positive");
        }
        vm.addBalance(amount);
    }

    @Override
    public DispenseResult dispense(VendingMachine vm) {
        vm.setState(new DispensingState());
        return vm.getState().dispense(vm);
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
