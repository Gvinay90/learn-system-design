public class ExpenseRequest {
    private final double amount;
    private final String description;

    public ExpenseRequest(double amount, String description) {
        this.amount = amount;
        this.description = description;
    }

    public double getAmount() { return amount; }
    public String getDescription() { return description; }
}
