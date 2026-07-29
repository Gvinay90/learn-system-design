import java.util.ArrayList;
import java.util.List;

/**
 * Delegates final price computation to whatever PricingStrategy it's
 * configured with, so pricing schemes can change at runtime without
 * touching cart logic.
 */
public class ShoppingCart {
    private final List<Item> items = new ArrayList<>();
    private PricingStrategy strategy;

    public ShoppingCart(PricingStrategy strategy) {
        this.strategy = strategy;
    }

    public void addItem(Item item) {
        items.add(item);
    }

    public void setStrategy(PricingStrategy strategy) {
        this.strategy = strategy;
    }

    public double subtotal() {
        double total = 0;
        for (Item item : items) {
            total += item.getPrice() * item.getQuantity();
        }
        return total;
    }

    public double checkout() {
        return strategy.applyDiscount(subtotal());
    }
}
