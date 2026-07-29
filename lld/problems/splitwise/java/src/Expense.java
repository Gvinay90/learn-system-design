import java.util.List;

public class Expense {
    private final String id;
    private final String description;
    private final User paidBy;
    private final double amount;
    private final List<User> participants;
    private final SplitStrategy strategy;

    public Expense(String id, String description, User paidBy, double amount,
                    List<User> participants, SplitStrategy strategy) {
        this.id = id;
        this.description = description;
        this.paidBy = paidBy;
        this.amount = amount;
        this.participants = participants;
        this.strategy = strategy;
    }

    public String getId() { return id; }
    public String getDescription() { return description; }
    public User getPaidBy() { return paidBy; }
    public double getAmount() { return amount; }
    public List<User> getParticipants() { return participants; }
    public SplitStrategy getStrategy() { return strategy; }
}
