public class Main {
    public static void main(String[] args) {
        ShoppingCart cart = new ShoppingCart(new RegularPricing());
        cart.addItem(new Item("book", 20, 2));
        System.out.println("Regular checkout: " + cart.checkout());

        cart.setStrategy(new PercentageDiscountPricing(10));
        System.out.println("10% off checkout: " + cart.checkout());

        ShoppingCartTest.runAll();
    }
}
