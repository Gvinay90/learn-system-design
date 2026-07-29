import java.util.List;
import java.util.Map;

public class Main {
    public static void main(String[] args) {
        Splitwise app = new Splitwise();

        User alice = app.addUser("u1", "Alice");
        User bob = app.addUser("u2", "Bob");
        User carol = app.addUser("u3", "Carol");
        User dave = app.addUser("u4", "Dave");
        app.addGroup("g1", "Goa Trip", List.of(alice, bob, carol, dave));

        app.addExpense("Dinner", alice, 800, List.of(alice, bob, carol, dave), new EqualSplit());
        app.addExpense("Cabs", bob, 300,
                List.of(bob, carol),
                new ExactSplit(Map.of(bob.getId(), 100.0, carol.getId(), 200.0)));
        app.addExpense("Hotel", carol, 1000,
                List.of(alice, bob, carol, dave),
                new PercentSplit(Map.of(
                        alice.getId(), 25.0,
                        bob.getId(), 25.0,
                        carol.getId(), 25.0,
                        dave.getId(), 25.0)));

        System.out.println("Net balances:");
        Map<String, Double> net = app.getLedger().netBalances();
        for (Map.Entry<String, Double> entry : net.entrySet()) {
            System.out.println("  " + entry.getKey() + " -> " + entry.getValue());
        }

        System.out.println("Simplified settlement plan:");
        for (Transaction t : DebtSimplifier.simplify(net)) {
            System.out.println("  " + t);
        }

        SplitwiseTest.runAll();
    }
}
