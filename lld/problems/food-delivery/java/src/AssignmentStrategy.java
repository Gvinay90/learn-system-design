import java.util.List;

public interface AssignmentStrategy {
    DeliveryPartner assign(Restaurant restaurant, List<DeliveryPartner> partners);
}
