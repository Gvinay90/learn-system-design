import java.util.List;
import java.util.Map;

public interface SplitStrategy {
    Map<String, Double> compute(double totalAmount, List<User> participants);
}
