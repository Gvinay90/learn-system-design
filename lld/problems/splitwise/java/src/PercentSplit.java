import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class PercentSplit implements SplitStrategy {
    private static final double EPSILON = 1e-6;
    private final Map<String, Double> percentages;

    public PercentSplit(Map<String, Double> percentages) {
        this.percentages = percentages;
    }

    @Override
    public Map<String, Double> compute(double totalAmount, List<User> participants) {
        Map<String, Double> shares = new HashMap<>();
        double sum = 0;
        for (User p : participants) {
            Double pct = percentages.get(p.getId());
            if (pct == null) {
                throw new InvalidSplitException("percent split missing percentage for " + p.getId());
            }
            shares.put(p.getId(), totalAmount * pct / 100.0);
            sum += pct;
        }
        if (Math.abs(sum - 100.0) > EPSILON) {
            throw new InvalidSplitException(
                    String.format("percent split percentages sum to %.2f, want 100", sum));
        }
        return shares;
    }
}
