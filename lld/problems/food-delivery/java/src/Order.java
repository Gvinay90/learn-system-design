import java.time.Instant;
import java.util.List;

public class Order {
    private final String id;
    private final Customer customer;
    private final Restaurant restaurant;
    private final List<MenuItem> items;
    private final Instant placedAt;
    private OrderStatus status;
    private DeliveryPartner partner;

    public Order(String id, Customer customer, Restaurant restaurant, List<MenuItem> items, Instant placedAt) {
        this.id = id;
        this.customer = customer;
        this.restaurant = restaurant;
        this.items = items;
        this.placedAt = placedAt;
        this.status = OrderStatus.PLACED;
    }

    public String getId() { return id; }
    public Customer getCustomer() { return customer; }
    public Restaurant getRestaurant() { return restaurant; }
    public List<MenuItem> getItems() { return items; }
    public Instant getPlacedAt() { return placedAt; }
    public OrderStatus getStatus() { return status; }
    public void setStatus(OrderStatus status) { this.status = status; }
    public DeliveryPartner getPartner() { return partner; }
    public void setPartner(DeliveryPartner partner) { this.partner = partner; }
}
