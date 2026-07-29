import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

// Records one row per successful payment (debit-payer/credit-payee collapsed
// into a single entry to keep the example focused on idempotency/retries).
public class Ledger {
    private final List<LedgerEntry> entries = new ArrayList<>();

    public synchronized void record(LedgerEntry entry) {
        entries.add(entry);
    }

    public synchronized List<LedgerEntry> getEntries() {
        return Collections.unmodifiableList(new ArrayList<>(entries));
    }
}
