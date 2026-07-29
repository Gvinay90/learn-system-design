import java.util.HashMap;
import java.util.Map;

/**
 * Plain-assertion test suite (no JUnit dependency) mirroring the Go test
 * cases in vending_machine_test.go. Run via {@code java -cp out Main}.
 */
public class VendingMachineTest {

    private static VendingMachine newTestMachine() {
        Map<String, Slot> slots = new HashMap<>();
        slots.put("A1", new Slot("Soda", 150, 2));
        slots.put("A2", new Slot("Chips", 200, 1));
        slots.put("B1", new Slot("Water", 100, 0));
        return new VendingMachine(slots);
    }

    private static void check(boolean condition, String message) {
        if (!condition) {
            throw new AssertionError(message);
        }
    }

    static void testSelectAndInsertExactMoneyThenDispense() {
        VendingMachine vm = newTestMachine();
        vm.selectItem("A1");
        vm.insertMoney(150);
        check(vm.getState().name().equals("HasMoney"), "expected HasMoney state");

        DispenseResult result = vm.dispense();
        check(result.item.equals("Soda"), "expected Soda");
        check(result.change == 0, "expected no change");
        check(vm.getState().name().equals("Idle"), "expected Idle state after dispense");
    }

    static void testInsertMoneyGivesChangeOnOverpay() {
        VendingMachine vm = newTestMachine();
        vm.selectItem("A1");
        vm.insertMoney(200);
        DispenseResult result = vm.dispense();
        check(result.change == 50, "expected 50 change, got " + result.change);
    }

    static void testInsufficientMoneyThenTopUpThenDispense() {
        VendingMachine vm = newTestMachine();
        vm.selectItem("A2"); // Chips, price 200
        vm.insertMoney(100);
        check(vm.getState().name().equals("Idle"), "expected still Idle after partial payment");

        try {
            vm.dispense();
            throw new AssertionError("expected NotEnoughMoneyException");
        } catch (VendingState.NotEnoughMoneyException expected) {
            // expected
        }

        vm.insertMoney(100);
        check(vm.getState().name().equals("HasMoney"), "expected HasMoney after top-up");

        DispenseResult result = vm.dispense();
        check(result.item.equals("Chips") && result.change == 0, "unexpected dispense result");
    }

    static void testCancelRefundsInsertedMoney() {
        VendingMachine vm = newTestMachine();
        vm.selectItem("A1");
        vm.insertMoney(75);

        int refund = vm.cancel();
        check(refund == 75, "expected refund of 75, got " + refund);
        check(vm.getState().name().equals("Idle"), "expected Idle state");
    }

    static void testSelectOutOfStockSlot() {
        VendingMachine vm = newTestMachine();
        try {
            vm.selectItem("B1");
            throw new AssertionError("expected OutOfStockException");
        } catch (VendingState.OutOfStockException expected) {
            // expected
        }
        check(vm.getState().name().equals("OutOfStock"), "expected OutOfStock state");

        vm.selectItem("A1");
        check(vm.getState().name().equals("Idle"), "expected recovery to Idle");
    }

    static void testInsertMoneyWithoutSelectionThrows() {
        VendingMachine vm = newTestMachine();
        try {
            vm.insertMoney(100);
            throw new AssertionError("expected NoSelectionException");
        } catch (VendingState.NoSelectionException expected) {
            // expected
        }
    }

    static void testSelectItemRejectedInHasMoneyState() {
        VendingMachine vm = newTestMachine();
        vm.selectItem("A1");
        vm.insertMoney(150);
        try {
            vm.selectItem("A2");
            throw new AssertionError("expected InvalidStateException");
        } catch (VendingState.InvalidStateException expected) {
            // expected
        }
    }

    static void testDispenseDecrementsInventory() {
        VendingMachine vm = newTestMachine();
        vm.selectItem("A1");
        vm.insertMoney(150);
        vm.dispense();
        check(vm.getInventory("A1").quantity == 1, "expected quantity decremented to 1");
    }

    public static void runAll() {
        testSelectAndInsertExactMoneyThenDispense();
        testInsertMoneyGivesChangeOnOverpay();
        testInsufficientMoneyThenTopUpThenDispense();
        testCancelRefundsInsertedMoney();
        testSelectOutOfStockSlot();
        testInsertMoneyWithoutSelectionThrows();
        testSelectItemRejectedInHasMoneyState();
        testDispenseDecrementsInventory();
        System.out.println("All VendingMachineTest cases passed.");
    }

    public static void main(String[] args) {
        runAll();
    }
}
