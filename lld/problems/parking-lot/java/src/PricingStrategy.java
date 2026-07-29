import java.time.Instant;

public interface PricingStrategy {
    double calculateFee(Ticket ticket, Instant exitTime);
}
