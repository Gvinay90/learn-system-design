import java.util.HashMap;
import java.util.Map;

public class Main {
    public static void main(String[] args) {
        Map<String, Slot> slots = new HashMap<>();
        slots.put("A1", new Slot("Soda", 150, 2));
        slots.put("A2", new Slot("Chips", 200, 1));

        VendingMachine vm = new VendingMachine(slots);
        System.out.println("Initial state: " + vm.getState().name());

        vm.selectItem("A1");
        vm.insertMoney(200);
        System.out.println("State after payment: " + vm.getState().name());

        DispenseResult result = vm.dispense();
        System.out.println("Dispensed " + result.item + ", change: " + result.change);
        System.out.println("Final state: " + vm.getState().name());

        VendingMachineTest.runAll();
    }
}
