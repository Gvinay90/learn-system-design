import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out SplitwiseTest` directly.
 */
public class SplitwiseTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testEqualSplitCreatesCorrectBalances();
        testExactSplitValidation();
        testPercentSplitCorrectness();
        testSimplifyDebts();
        System.out.println("All SplitwiseTest cases passed.");
    }

    private static void testEqualSplitCreatesCorrectBalances() {
        Splitwise app = new Splitwise();
        User alice = app.addUser("u1", "Alice");
        User bob = app.addUser("u2", "Bob");
        User carol = app.addUser("u3", "Carol");

        app.addExpense("Dinner", alice, 90, List.of(alice, bob, carol), new EqualSplit());

        assertEquals(30.0, app.getLedger().netBalance(bob.getId(), alice.getId()), "bob owes alice");
        assertEquals(30.0, app.getLedger().netBalance(carol.getId(), alice.getId()), "carol owes alice");
        assertEquals(-30.0, app.getLedger().netBalance(alice.getId(), bob.getId()), "anti-symmetric balance");
    }

    private static void testExactSplitValidation() {
        Splitwise app = new Splitwise();
        User alice = app.addUser("u1", "Alice");
        User bob = app.addUser("u2", "Bob");

        try {
            app.addExpense("Rent", alice, 100, List.of(alice, bob),
                    new ExactSplit(Map.of(alice.getId(), 40.0, bob.getId(), 50.0)));
            throw new AssertionError("expected InvalidSplitException for amounts not summing to total");
        } catch (InvalidSplitException expected) {
            // expected
        }

        app.addExpense("Rent", alice, 100, List.of(alice, bob),
                new ExactSplit(Map.of(alice.getId(), 40.0, bob.getId(), 60.0)));
        assertEquals(60.0, app.getLedger().netBalance(bob.getId(), alice.getId()), "bob owes alice exact 60");
    }

    private static void testPercentSplitCorrectness() {
        Splitwise app = new Splitwise();
        User alice = app.addUser("u1", "Alice");
        User bob = app.addUser("u2", "Bob");
        User carol = app.addUser("u3", "Carol");

        app.addExpense("Trip", alice, 200, List.of(alice, bob, carol),
                new PercentSplit(Map.of(alice.getId(), 50.0, bob.getId(), 25.0, carol.getId(), 25.0)));

        assertEquals(50.0, app.getLedger().netBalance(bob.getId(), alice.getId()), "bob owes alice 50");
        assertEquals(50.0, app.getLedger().netBalance(carol.getId(), alice.getId()), "carol owes alice 50");

        try {
            app.addExpense("Bad", alice, 200, List.of(alice, bob, carol),
                    new PercentSplit(Map.of(alice.getId(), 50.0, bob.getId(), 25.0, carol.getId(), 20.0)));
            throw new AssertionError("expected InvalidSplitException for percentages not summing to 100");
        } catch (InvalidSplitException expected) {
            // expected
        }
    }

    private static void testSimplifyDebts() {
        Map<String, Double> net = new HashMap<>();
        net.put("A", -300.0);
        net.put("B", -100.0);
        net.put("C", 200.0);
        net.put("D", 200.0);

        List<Transaction> txns = DebtSimplifier.simplify(net);
        if (txns.size() > 3) {
            throw new AssertionError("expected at most 3 transactions, got " + txns.size());
        }

        Map<String, Double> settled = new HashMap<>();
        for (Transaction t : txns) {
            settled.merge(t.getFrom(), -t.getAmount(), Double::sum);
            settled.merge(t.getTo(), t.getAmount(), Double::sum);
        }
        for (Map.Entry<String, Double> entry : net.entrySet()) {
            assertEquals(entry.getValue(), settled.getOrDefault(entry.getKey(), 0.0),
                    "settlement for " + entry.getKey());
        }
    }

    private static void assertEquals(double expected, double actual, String label) {
        if (Math.abs(expected - actual) > 1e-6) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
