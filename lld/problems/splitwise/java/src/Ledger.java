import java.util.HashMap;
import java.util.Map;

/**
 * Tracks net pairwise balances between users. balances[a][b] > 0 means b
 * owes a that amount; the map is always kept anti-symmetric
 * (balances[a][b] == -balances[b][a]).
 */
public class Ledger {
    private final Map<String, Map<String, Double>> balances = new HashMap<>();

    private synchronized void adjust(String a, String b, double amount) {
        balances.computeIfAbsent(a, k -> new HashMap<>()).merge(b, amount, Double::sum);
        balances.computeIfAbsent(b, k -> new HashMap<>()).merge(a, -amount, Double::sum);
    }

    /** Splits the expense per its strategy and credits the payer for every other participant's share. */
    public synchronized void addExpense(Expense e) {
        Map<String, Double> shares = e.getStrategy().compute(e.getAmount(), e.getParticipants());
        for (User p : e.getParticipants()) {
            if (p.getId().equals(e.getPaidBy().getId())) {
                continue;
            }
            adjust(e.getPaidBy().getId(), p.getId(), shares.get(p.getId()));
        }
    }

    /** Returns how much debtor owes creditor (may be negative if the balance runs the other way). */
    public synchronized double netBalance(String debtor, String creditor) {
        Map<String, Double> row = balances.get(creditor);
        if (row == null) {
            return 0;
        }
        return row.getOrDefault(debtor, 0.0);
    }

    /** Returns each user's overall net position: positive = net creditor, negative = net debtor. */
    public synchronized Map<String, Double> netBalances() {
        Map<String, Double> net = new HashMap<>();
        for (Map.Entry<String, Map<String, Double>> entry : balances.entrySet()) {
            double total = 0;
            for (double amt : entry.getValue().values()) {
                total += amt;
            }
            net.put(entry.getKey(), total);
        }
        return net;
    }
}
