import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ConcurrentMap;

public class IdempotencyStore {
    public static class Reservation {
        public final PaymentResult existingResult;
        public final boolean isOwner;

        Reservation(PaymentResult existingResult, boolean isOwner) {
            this.existingResult = existingResult;
            this.isOwner = isOwner;
        }
    }

    private final ConcurrentMap<String, CompletableFuture<PaymentResult>> entries = new ConcurrentHashMap<>();

    // Atomically claims key for the first caller (isOwner = true, who must
    // then call complete(key, result)); every other caller with the same
    // key blocks on the same future and receives the identical result.
    public Reservation reserveOrWait(String key) {
        CompletableFuture<PaymentResult> future = new CompletableFuture<>();
        CompletableFuture<PaymentResult> existing = entries.putIfAbsent(key, future);
        if (existing != null) {
            return new Reservation(existing.join(), false);
        }
        return new Reservation(null, true);
    }

    public void complete(String key, PaymentResult result) {
        entries.get(key).complete(result);
    }
}
