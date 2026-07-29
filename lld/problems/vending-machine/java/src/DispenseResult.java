/**
 * Result of a successful Dispense call.
 */
public class DispenseResult {
    public final String item;
    public final int change; // in cents

    public DispenseResult(String item, int change) {
        this.item = item;
        this.change = change;
    }
}
