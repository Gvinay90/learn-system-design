import java.util.List;

public class NearestAvailablePartnerStrategy implements AssignmentStrategy {
    @Override
    public DeliveryPartner assign(Restaurant restaurant, List<DeliveryPartner> partners) {
        DeliveryPartner nearest = null;
        double best = Double.MAX_VALUE;
        for (DeliveryPartner p : partners) {
            if (!p.isAvailable()) {
                continue;
            }
            double d = restaurant.getLocation().distanceTo(p.getLocation());
            if (d < best) {
                best = d;
                nearest = p;
            }
        }
        return nearest;
    }
}
