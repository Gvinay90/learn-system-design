public class ShoppingCartTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testRegularPricingNoDiscount();
        testPercentageDiscountPricing();
        testClearancePricing();
        testClearancePricingFloorsAtZero();
        testSwitchingStrategyAtRuntime();
        System.out.println("All ShoppingCartTest cases passed.");
    }

    private static void testRegularPricingNoDiscount() {
        ShoppingCart cart = new ShoppingCart(new RegularPricing());
        cart.addItem(new Item("book", 20, 2));
        assertEquals(40.0, cart.checkout(), "regular pricing");
    }

    private static void testPercentageDiscountPricing() {
        ShoppingCart cart = new ShoppingCart(new PercentageDiscountPricing(10));
        cart.addItem(new Item("shoes", 100, 1));
        assertEquals(90.0, cart.checkout(), "10% off");
    }

    private static void testClearancePricing() {
        ShoppingCart cart = new ShoppingCart(new ClearancePricing(20, 5));
        cart.addItem(new Item("jacket", 100, 1));
        assertEquals(75.0, cart.checkout(), "clearance pricing");
    }

    private static void testClearancePricingFloorsAtZero() {
        ShoppingCart cart = new ShoppingCart(new ClearancePricing(50, 100));
        cart.addItem(new Item("sticker", 5, 1));
        assertEquals(0.0, cart.checkout(), "floored at zero");
    }

    private static void testSwitchingStrategyAtRuntime() {
        ShoppingCart cart = new ShoppingCart(new RegularPricing());
        cart.addItem(new Item("widget", 50, 2));
        assertEquals(100.0, cart.checkout(), "before switch");

        cart.setStrategy(new PercentageDiscountPricing(25));
        assertEquals(75.0, cart.checkout(), "after switch");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
