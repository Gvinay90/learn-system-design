import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

public class PaymentResult {
    private final String id;
    private final PaymentRequest request;
    private final Instant createdAt;
    private final List<Attempt> attempts = new ArrayList<>();
    private volatile PaymentStatus status;

    public PaymentResult(String id, PaymentRequest request) {
        this.id = id;
        this.request = request;
        this.status = PaymentStatus.PENDING;
        this.createdAt = Instant.now();
    }

    public String getId() { return id; }
    public PaymentRequest getRequest() { return request; }
    public Instant getCreatedAt() { return createdAt; }
    public List<Attempt> getAttempts() { return attempts; }
    public PaymentStatus getStatus() { return status; }
    public void setStatus(PaymentStatus status) { this.status = status; }
    public void addAttempt(Attempt attempt) { attempts.add(attempt); }
}
