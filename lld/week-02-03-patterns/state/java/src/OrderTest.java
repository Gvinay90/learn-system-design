public class OrderTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testHappyPathTransitions();
        testCannotShipBeforePayment();
        testCancelAllowedBeforeShipping();
        testCannotCancelAfterShipping();
        testDeliveredIsTerminal();
        System.out.println("All OrderTest cases passed.");
    }

    private static void testHappyPathTransitions() {
        Order order = new Order("O-1");
        assertEquals("Created", order.getState().name(), "initial state");

        order.pay();
        assertEquals("Paid", order.getState().name(), "after pay");

        order.ship();
        assertEquals("Shipped", order.getState().name(), "after ship");

        order.deliver();
        assertEquals("Delivered", order.getState().name(), "after deliver");
    }

    private static void testCannotShipBeforePayment() {
        Order order = new Order("O-2");
        try {
            order.ship();
            throw new AssertionError("expected IllegalTransitionException");
        } catch (OrderState.IllegalTransitionException e) {
            // expected
        }
        assertEquals("Created", order.getState().name(), "state unchanged");
    }

    private static void testCancelAllowedBeforeShipping() {
        Order order = new Order("O-3");
        order.pay();
        order.cancel();
        assertEquals("Cancelled", order.getState().name(), "after cancel before shipping");
    }

    private static void testCannotCancelAfterShipping() {
        Order order = new Order("O-4");
        order.pay();
        order.ship();
        try {
            order.cancel();
            throw new AssertionError("expected IllegalTransitionException");
        } catch (OrderState.IllegalTransitionException e) {
            // expected
        }
        assertEquals("Shipped", order.getState().name(), "state unchanged");
    }

    private static void testDeliveredIsTerminal() {
        Order order = new Order("O-5");
        order.pay();
        order.ship();
        order.deliver();

        try {
            order.cancel();
            throw new AssertionError("expected IllegalTransitionException on cancel");
        } catch (OrderState.IllegalTransitionException e) {
            // expected
        }
        try {
            order.pay();
            throw new AssertionError("expected IllegalTransitionException on pay");
        } catch (OrderState.IllegalTransitionException e) {
            // expected
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
