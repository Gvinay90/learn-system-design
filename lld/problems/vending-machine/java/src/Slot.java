/**
 * Inventory for a single selectable position in the machine.
 */
public class Slot {
    public final String item;
    public final int price; // in cents
    public int quantity;

    public Slot(String item, int price, int quantity) {
        this.item = item;
        this.price = price;
        this.quantity = quantity;
    }

    public Slot copy() {
        return new Slot(item, price, quantity);
    }
}
