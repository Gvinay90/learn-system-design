import java.util.HashMap;
import java.util.Map;

/**
 * The context: holds the current state, the inventory, the currently
 * selected slot, and the money collected so far, and forwards every action
 * to its current state.
 */
public class VendingMachine {
    private final Map<String, Slot> slots = new HashMap<>();
    private VendingState state = new IdleState();
    private String selected = "";
    private int balance = 0; // in cents

    public VendingMachine(Map<String, Slot> initialSlots) {
        for (Map.Entry<String, Slot> e : initialSlots.entrySet()) {
            slots.put(e.getKey(), e.getValue().copy());
        }
    }

    public VendingState getState() { return state; }
    public int getBalance() { return balance; }
    public String getSelected() { return selected; }

    public Slot getInventory(String slotId) {
        Slot s = slots.get(slotId);
        return s == null ? null : s.copy();
    }

    public void selectItem(String slotId) { state.selectItem(this, slotId); }
    public void insertMoney(int amount) { state.insertMoney(this, amount); }
    public DispenseResult dispense() { return state.dispense(this); }
    public int cancel() { return state.cancel(this); }

    void setState(VendingState state) { this.state = state; }
    void setSelected(String selected) { this.selected = selected; }
    void setBalance(int balance) { this.balance = balance; }
    void addBalance(int amount) { this.balance += amount; }

    Slot slot(String slotId) {
        Slot s = slots.get(slotId);
        if (s == null) {
            throw new VendingState.InvalidSlotException("no such slot: " + slotId);
        }
        return s;
    }
}
