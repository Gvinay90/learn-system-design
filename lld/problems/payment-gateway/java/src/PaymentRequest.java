public class PaymentRequest {
    private final String idempotencyKey;
    private final String payerId;
    private final String payeeId;
    private final double amount;
    private final String currency;

    public PaymentRequest(String idempotencyKey, String payerId, String payeeId, double amount, String currency) {
        this.idempotencyKey = idempotencyKey;
        this.payerId = payerId;
        this.payeeId = payeeId;
        this.amount = amount;
        this.currency = currency;
    }

    public String getIdempotencyKey() { return idempotencyKey; }
    public String getPayerId() { return payerId; }
    public String getPayeeId() { return payeeId; }
    public double getAmount() { return amount; }
    public String getCurrency() { return currency; }
}
