/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out DecoratorTest` directly.
 */
public class DecoratorTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testBaseCoffee();
        testSingleDecorator();
        testStackedDecoratorsCumulativeCostAndDescription();
        testDecoratorOrderIndependenceOfCost();
        System.out.println("All DecoratorTest cases passed.");
    }

    private static void testBaseCoffee() {
        Coffee c = new Espresso();
        assertEquals(2.0, c.cost(), "base cost");
        assertEquals("Espresso", c.description(), "base description");
    }

    private static void testSingleDecorator() {
        Coffee c = new MilkDecorator(new Espresso());
        assertEquals(2.5, c.cost(), "single decorator cost");
        assertEquals("Espresso + Milk", c.description(), "single decorator description");
    }

    private static void testStackedDecoratorsCumulativeCostAndDescription() {
        Coffee c = new Espresso();
        c = new MilkDecorator(c);
        c = new SugarDecorator(c);
        c = new WhipDecorator(c);

        assertEquals(2.0 + 0.5 + 0.25 + 0.75, c.cost(), "stacked cost");
        assertEquals("Espresso + Milk + Sugar + Whip", c.description(), "stacked description");
    }

    private static void testDecoratorOrderIndependenceOfCost() {
        Coffee a = new WhipDecorator(new MilkDecorator(new Espresso()));
        Coffee b = new MilkDecorator(new WhipDecorator(new Espresso()));
        assertEquals(a.cost(), b.cost(), "cost should be order-independent");
        if (a.description().equals(b.description())) {
            throw new AssertionError("descriptions should differ by wrap order");
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
