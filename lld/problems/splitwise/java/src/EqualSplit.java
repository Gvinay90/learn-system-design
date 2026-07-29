import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class EqualSplit implements SplitStrategy {
    @Override
    public Map<String, Double> compute(double totalAmount, List<User> participants) {
        if (participants.isEmpty()) {
            throw new InvalidSplitException("equal split requires at least one participant");
        }
        double share = totalAmount / participants.size();
        Map<String, Double> shares = new HashMap<>();
        for (User p : participants) {
            shares.put(p.getId(), share);
        }
        return shares;
    }
}
