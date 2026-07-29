import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ExactSplit implements SplitStrategy {
    private static final double EPSILON = 1e-6;
    private final Map<String, Double> amounts;

    public ExactSplit(Map<String, Double> amounts) {
        this.amounts = amounts;
    }

    @Override
    public Map<String, Double> compute(double totalAmount, List<User> participants) {
        Map<String, Double> shares = new HashMap<>();
        double sum = 0;
        for (User p : participants) {
            Double amt = amounts.get(p.getId());
            if (amt == null) {
                throw new InvalidSplitException("exact split missing amount for " + p.getId());
            }
            shares.put(p.getId(), amt);
            sum += amt;
        }
        if (Math.abs(sum - totalAmount) > EPSILON) {
            throw new InvalidSplitException(
                    String.format("exact split amounts sum to %.2f, want %.2f", sum, totalAmount));
        }
        return shares;
    }
}
