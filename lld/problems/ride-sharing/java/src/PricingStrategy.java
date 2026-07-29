/** Computes the fare for a completed trip. */
public interface PricingStrategy {
    double calculateFare(Trip trip);
}
