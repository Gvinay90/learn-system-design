import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.Map;

/**
 * Greedy debt simplification: repeatedly matches the largest creditor with
 * the largest debtor. It does not guarantee the theoretical minimum in every
 * adversarial case, but it is optimal for the common case and is what real
 * Splitwise-style systems use.
 */
public class DebtSimplifier {
    private static final double EPSILON = 1e-6;

    private static class Balance {
        String id;
        double amount;
        Balance(String id, double amount) { this.id = id; this.amount = amount; }
    }

    public static List<Transaction> simplify(Map<String, Double> net) {
        List<Balance> creditors = new ArrayList<>();
        List<Balance> debtors = new ArrayList<>();
        for (Map.Entry<String, Double> entry : net.entrySet()) {
            double amt = entry.getValue();
            if (amt > EPSILON) {
                creditors.add(new Balance(entry.getKey(), amt));
            } else if (amt < -EPSILON) {
                debtors.add(new Balance(entry.getKey(), -amt));
            }
        }

        creditors.sort(Comparator.comparingDouble((Balance b) -> b.amount).reversed());
        debtors.sort(Comparator.comparingDouble((Balance b) -> b.amount).reversed());

        List<Transaction> transactions = new ArrayList<>();
        int i = 0, j = 0;
        while (i < creditors.size() && j < debtors.size()) {
            Balance c = creditors.get(i);
            Balance d = debtors.get(j);
            double amount = Math.min(c.amount, d.amount);

            transactions.add(new Transaction(d.id, c.id, amount));

            c.amount -= amount;
            d.amount -= amount;
            if (c.amount <= EPSILON) {
                i++;
            }
            if (d.amount <= EPSILON) {
                j++;
            }
        }
        return transactions;
    }
}
