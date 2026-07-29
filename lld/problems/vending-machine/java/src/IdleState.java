/**
 * No completed transaction in progress. An item may or may not be selected
 * yet, and partial money may have been inserted; once inserted money
 * covers the selected item's price the machine moves to HasMoney.
 */
public class IdleState implements VendingState {
    @Override
    public String name() { return "Idle"; }

    @Override
    public void selectItem(VendingMachine vm, String slotId) {
        StateHelpers.selectItem(vm, slotId);
    }

    @Override
    public void insertMoney(VendingMachine vm, int amount) {
        if (amount <= 0) {
            throw new InvalidAmountException("amount must be positive");
        }
        if (vm.getSelected().isEmpty()) {
            throw new NoSelectionException("no item selected");
        }
        Slot s = vm.slot(vm.getSelected());
        vm.addBalance(amount);
        if (vm.getBalance() >= s.price) {
            vm.setState(new HasMoneyState());
        }
    }

    @Override
    public DispenseResult dispense(VendingMachine vm) {
        throw new NotEnoughMoneyException("not enough money inserted");
    }

    @Override
    public int cancel(VendingMachine vm) {
        int refund = vm.getBalance();
        vm.setBalance(0);
        vm.setSelected("");
        return refund;
    }
}
