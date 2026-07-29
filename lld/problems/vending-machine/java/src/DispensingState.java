/**
 * The machine is actively completing a purchase. This is entered only
 * transiently from HasMoneyState.dispense, which performs the inventory
 * decrement and change calculation and then falls back to Idle; it rejects
 * every other action as busy.
 */
public class DispensingState implements VendingState {
    @Override
    public String name() { return "Dispensing"; }

    @Override
    public void selectItem(VendingMachine vm, String slotId) {
        throw new InvalidStateException("machine is busy dispensing");
    }

    @Override
    public void insertMoney(VendingMachine vm, int amount) {
        throw new InvalidStateException("machine is busy dispensing");
    }

    @Override
    public DispenseResult dispense(VendingMachine vm) {
        Slot s = vm.slot(vm.getSelected());
        if (vm.getBalance() < s.price) {
            vm.setState(new HasMoneyState());
            throw new NotEnoughMoneyException("not enough money inserted");
        }

        int change = vm.getBalance() - s.price;
        s.quantity--;
        String item = s.item;

        vm.setBalance(0);
        vm.setSelected("");
        vm.setState(new IdleState());

        return new DispenseResult(item, change);
    }

    @Override
    public int cancel(VendingMachine vm) {
        throw new InvalidStateException("machine is busy dispensing");
    }
}
