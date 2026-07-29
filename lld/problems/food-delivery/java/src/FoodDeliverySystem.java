import java.time.Instant;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class FoodDeliverySystem {
    public static class RestaurantClosedException extends RuntimeException {
        public RestaurantClosedException() { super("restaurant is closed"); }
    }
    public static class ItemNotOnMenuException extends RuntimeException {
        public ItemNotOnMenuException() { super("item not on restaurant menu"); }
    }
    public static class OrderNotFoundException extends RuntimeException {
        public OrderNotFoundException() { super("order not found"); }
    }
    public static class NoPartnerAvailableException extends RuntimeException {
        public NoPartnerAvailableException() { super("no delivery partner available"); }
    }
    public static class InvalidTransitionException extends RuntimeException {
        public InvalidTransitionException() { super("invalid order status transition"); }
    }
    public static class OrderAlreadyAssignedException extends RuntimeException {
        public OrderAlreadyAssignedException() { super("order already has an assigned delivery partner"); }
    }

    private static final Map<OrderStatus, Map<OrderStatus, Boolean>> VALID_NEXT = Map.of(
            OrderStatus.PLACED, Map.of(OrderStatus.ACCEPTED, true, OrderStatus.CANCELLED, true),
            OrderStatus.ACCEPTED, Map.of(OrderStatus.PREPARING, true, OrderStatus.CANCELLED, true),
            OrderStatus.PREPARING, Map.of(OrderStatus.OUT_FOR_DELIVERY, true),
            OrderStatus.OUT_FOR_DELIVERY, Map.of(OrderStatus.DELIVERED, true),
            OrderStatus.DELIVERED, Map.of(),
            OrderStatus.CANCELLED, Map.of()
    );

    private final AssignmentStrategy strategy;
    private final List<DeliveryPartner> partners;
    private final Map<String, Order> orders = new ConcurrentHashMap<>();
    private final AtomicInteger seq = new AtomicInteger(0);
    private final Object lock = new Object();

    public FoodDeliverySystem(List<DeliveryPartner> partners, AssignmentStrategy strategy) {
        this.partners = partners;
        this.strategy = strategy;
    }

    public Order placeOrder(Customer customer, Restaurant restaurant, List<String> itemIds) {
        if (!restaurant.isOpen()) {
            throw new RestaurantClosedException();
        }
        List<MenuItem> items = itemIds.stream()
                .map(id -> restaurant.findItem(id).orElseThrow(ItemNotOnMenuException::new))
                .toList();

        String id = "O-" + seq.incrementAndGet();
        Order order = new Order(id, customer, restaurant, items, Instant.now());
        orders.put(id, order);
        return order;
    }

    public DeliveryPartner assignDeliveryPartner(String orderId) {
        synchronized (lock) {
            Order order = orders.get(orderId);
            if (order == null) {
                throw new OrderNotFoundException();
            }
            if (order.getPartner() != null) {
                throw new OrderAlreadyAssignedException();
            }
            DeliveryPartner partner = strategy.assign(order.getRestaurant(), partners);
            if (partner == null) {
                throw new NoPartnerAvailableException();
            }
            partner.setAvailable(false);
            order.setPartner(partner);
            return partner;
        }
    }

    public void updateOrderStatus(String orderId, OrderStatus next) {
        synchronized (lock) {
            Order order = orders.get(orderId);
            if (order == null) {
                throw new OrderNotFoundException();
            }
            if (!VALID_NEXT.get(order.getStatus()).getOrDefault(next, false)) {
                throw new InvalidTransitionException();
            }
            order.setStatus(next);
            if ((next == OrderStatus.DELIVERED || next == OrderStatus.CANCELLED) && order.getPartner() != null) {
                order.getPartner().setAvailable(true);
                order.setPartner(null);
            }
        }
    }

    public Order getOrder(String orderId) {
        Order order = orders.get(orderId);
        if (order == null) {
            throw new OrderNotFoundException();
        }
        return order;
    }
}
