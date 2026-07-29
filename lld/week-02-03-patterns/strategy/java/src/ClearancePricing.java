public class ClearancePricing implements PricingStrategy {
    private final double percentOff;
    private final double flatOff;

    public ClearancePricing(double percentOff, double flatOff) {
        this.percentOff = percentOff;
        this.flatOff = flatOff;
    }

    @Override
    public double applyDiscount(double subtotal) {
        double discounted = subtotal * (1 - percentOff / 100) - flatOff;
        return Math.max(discounted, 0);
    }
}
