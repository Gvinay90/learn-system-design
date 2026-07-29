public class RegularPricing implements PricingStrategy {
    @Override
    public double applyDiscount(double subtotal) {
        return subtotal;
    }
}
