/** Charges a base fare plus a per-unit-distance rate. */
public class DistanceBasedPricing implements PricingStrategy {
    private final double baseFare;
    private final double perDistance;

    public DistanceBasedPricing(double baseFare, double perDistance) {
        this.baseFare = baseFare;
        this.perDistance = perDistance;
    }

    @Override
    public double calculateFare(Trip trip) {
        return baseFare + perDistance * trip.pickup.distanceTo(trip.dropoff);
    }
}
