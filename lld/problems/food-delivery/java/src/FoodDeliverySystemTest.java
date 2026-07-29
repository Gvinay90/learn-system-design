import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out FoodDeliverySystemTest` directly.
 */
public class FoodDeliverySystemTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testHappyPathPlaceAssignDeliver();
        testInvalidStatusTransitionRejected();
        testItemNotOnMenuAndNoPartnerAvailable();
        testConcurrentAssignment();
        System.out.println("All FoodDeliverySystemTest cases passed.");
    }

    private static Restaurant testRestaurant() {
        return new Restaurant("R1", "Tasty Bites", new Location(0, 0), List.of(
                new MenuItem("I1", "Burger", 5),
                new MenuItem("I2", "Fries", 2)
        ), true);
    }

    private static void testHappyPathPlaceAssignDeliver() {
        Restaurant restaurant = testRestaurant();
        DeliveryPartner partner = new DeliveryPartner("P1", "Alex", new Location(1, 1), true);
        FoodDeliverySystem sys = new FoodDeliverySystem(List.of(partner), new NearestAvailablePartnerStrategy());
        Customer customer = new Customer("C1", "Sam");

        Order order = sys.placeOrder(customer, restaurant, List.of("I1", "I2"));

        DeliveryPartner assigned = sys.assignDeliveryPartner(order.getId());
        assertEquals("P1", assigned.getId(), "assigned partner");
        assertEquals(false, partner.isAvailable(), "partner marked unavailable after assignment");

        for (OrderStatus next : List.of(OrderStatus.ACCEPTED, OrderStatus.PREPARING, OrderStatus.OUT_FOR_DELIVERY, OrderStatus.DELIVERED)) {
            sys.updateOrderStatus(order.getId(), next);
        }

        assertEquals(OrderStatus.DELIVERED, sys.getOrder(order.getId()).getStatus(), "order delivered");
        assertEquals(true, partner.isAvailable(), "partner freed after delivery");
    }

    private static void testInvalidStatusTransitionRejected() {
        Restaurant restaurant = testRestaurant();
        DeliveryPartner partner = new DeliveryPartner("P1", "Alex", new Location(1, 1), true);
        FoodDeliverySystem sys = new FoodDeliverySystem(List.of(partner), new NearestAvailablePartnerStrategy());
        Customer customer = new Customer("C1", "Sam");

        Order order = sys.placeOrder(customer, restaurant, List.of("I1"));

        try {
            sys.updateOrderStatus(order.getId(), OrderStatus.DELIVERED);
            throw new AssertionError("expected InvalidTransitionException jumping straight to Delivered");
        } catch (FoodDeliverySystem.InvalidTransitionException e) {
            // expected
        }
    }

    private static void testItemNotOnMenuAndNoPartnerAvailable() {
        Restaurant restaurant = testRestaurant();
        FoodDeliverySystem sys = new FoodDeliverySystem(List.of(), new NearestAvailablePartnerStrategy());
        Customer customer = new Customer("C1", "Sam");

        try {
            sys.placeOrder(customer, restaurant, List.of("BOGUS"));
            throw new AssertionError("expected ItemNotOnMenuException");
        } catch (FoodDeliverySystem.ItemNotOnMenuException e) {
            // expected
        }

        Order order = sys.placeOrder(customer, restaurant, List.of("I1"));
        try {
            sys.assignDeliveryPartner(order.getId());
            throw new AssertionError("expected NoPartnerAvailableException");
        } catch (FoodDeliverySystem.NoPartnerAvailableException e) {
            // expected
        }
    }

    private static void testConcurrentAssignment() {
        Restaurant restaurant = testRestaurant();
        DeliveryPartner partner = new DeliveryPartner("P1", "Alex", new Location(1, 1), true);
        FoodDeliverySystem sys = new FoodDeliverySystem(List.of(partner), new NearestAvailablePartnerStrategy());
        Customer customer = new Customer("C1", "Sam");

        Order order1 = sys.placeOrder(customer, restaurant, List.of("I1"));
        Order order2 = sys.placeOrder(customer, restaurant, List.of("I2"));

        AtomicInteger successCount = new AtomicInteger(0);
        CountDownLatch latch = new CountDownLatch(2);

        Runnable makeTask = () -> {};
        Runnable task1 = () -> {
            try {
                sys.assignDeliveryPartner(order1.getId());
                successCount.incrementAndGet();
            } catch (FoodDeliverySystem.NoPartnerAvailableException | FoodDeliverySystem.OrderAlreadyAssignedException ignored) {
            } finally {
                latch.countDown();
            }
        };
        Runnable task2 = () -> {
            try {
                sys.assignDeliveryPartner(order2.getId());
                successCount.incrementAndGet();
            } catch (FoodDeliverySystem.NoPartnerAvailableException | FoodDeliverySystem.OrderAlreadyAssignedException ignored) {
            } finally {
                latch.countDown();
            }
        };
        new Thread(task1).start();
        new Thread(task2).start();
        try {
            latch.await();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        assertEquals(1, successCount.get(), "exactly one order should win the single available partner");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
