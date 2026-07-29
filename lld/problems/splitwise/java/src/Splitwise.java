import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class Splitwise {
    private final Map<String, User> users = new ConcurrentHashMap<>();
    private final Map<String, Group> groups = new ConcurrentHashMap<>();
    private final Ledger ledger = new Ledger();
    private final AtomicInteger seq = new AtomicInteger(0);

    public User addUser(String id, String name) {
        User u = new User(id, name);
        users.put(id, u);
        return u;
    }

    public Group addGroup(String id, String name, List<User> members) {
        Group g = new Group(id, name, members);
        groups.put(id, g);
        return g;
    }

    public Expense addExpense(String description, User paidBy, double amount,
                               List<User> participants, SplitStrategy strategy) {
        String id = "E-" + seq.incrementAndGet();
        Expense e = new Expense(id, description, paidBy, amount, participants, strategy);
        ledger.addExpense(e);
        return e;
    }

    public Ledger getLedger() { return ledger; }
}
