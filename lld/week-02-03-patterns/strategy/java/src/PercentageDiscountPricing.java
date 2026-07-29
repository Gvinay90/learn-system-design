public class PercentageDiscountPricing implements PricingStrategy {
    private final double percentOff;

    public PercentageDiscountPricing(double percentOff) {
        this.percentOff = percentOff;
    }

    @Override
    public double applyDiscount(double subtotal) {
        return subtotal * (1 - percentOff / 100);
    }
}
