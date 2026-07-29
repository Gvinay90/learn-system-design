import java.time.Instant;

public class LedgerEntry {
    private final String paymentId;
    private final String payerId;
    private final String payeeId;
    private final double amount;
    private final Instant at;

    public LedgerEntry(String paymentId, String payerId, String payeeId, double amount, Instant at) {
        this.paymentId = paymentId;
        this.payerId = payerId;
        this.payeeId = payeeId;
        this.amount = amount;
        this.at = at;
    }

    public String getPaymentId() { return paymentId; }
    public String getPayerId() { return payerId; }
    public String getPayeeId() { return payeeId; }
    public double getAmount() { return amount; }
    public Instant getAt() { return at; }
}
