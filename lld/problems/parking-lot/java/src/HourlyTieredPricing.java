import java.time.Duration;
import java.time.Instant;

public class HourlyTieredPricing implements PricingStrategy {
    private final double firstHourRate;
    private final double subsequentRate;

    public HourlyTieredPricing(double firstHourRate, double subsequentRate) {
        this.firstHourRate = firstHourRate;
        this.subsequentRate = subsequentRate;
    }

    @Override
    public double calculateFee(Ticket ticket, Instant exitTime) {
        double hours = Duration.between(ticket.getEntryTime(), exitTime).toMinutes() / 60.0;
        if (hours <= 1) {
            return firstHourRate;
        }
        return firstHourRate + (hours - 1) * subsequentRate;
    }
}
