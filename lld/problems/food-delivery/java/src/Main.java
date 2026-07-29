import java.util.List;

public class Main {
    public static void main(String[] args) {
        Restaurant restaurant = new Restaurant("R1", "Tasty Bites", new Location(0, 0), List.of(
                new MenuItem("I1", "Burger", 5),
                new MenuItem("I2", "Fries", 2)
        ), true);
        DeliveryPartner partner = new DeliveryPartner("P1", "Alex", new Location(1, 1), true);
        FoodDeliverySystem sys = new FoodDeliverySystem(List.of(partner), new NearestAvailablePartnerStrategy());
        Customer customer = new Customer("C1", "Sam");

        Order order = sys.placeOrder(customer, restaurant, List.of("I1", "I2"));
        System.out.println("Placed order " + order.getId());

        DeliveryPartner assigned = sys.assignDeliveryPartner(order.getId());
        System.out.println("Assigned partner " + assigned.getId());

        sys.updateOrderStatus(order.getId(), OrderStatus.ACCEPTED);
        sys.updateOrderStatus(order.getId(), OrderStatus.PREPARING);
        sys.updateOrderStatus(order.getId(), OrderStatus.OUT_FOR_DELIVERY);
        sys.updateOrderStatus(order.getId(), OrderStatus.DELIVERED);
        System.out.println("Final status: " + sys.getOrder(order.getId()).getStatus());

        FoodDeliverySystemTest.runAll();
    }
}
