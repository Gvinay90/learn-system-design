import java.time.Instant;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class ParkingLot {
    public static class LotFullException extends RuntimeException {
        public LotFullException() { super("no available spot for vehicle type"); }
    }
    public static class TicketNotFoundException extends RuntimeException {
        public TicketNotFoundException() { super("ticket not recognized"); }
    }

    private final List<Floor> floors;
    private final PricingStrategy pricing;
    private final Map<String, Ticket> tickets = new ConcurrentHashMap<>();
    private final AtomicInteger seq = new AtomicInteger(0);

    public ParkingLot(List<Floor> floors, PricingStrategy pricing) {
        this.floors = floors;
        this.pricing = pricing;
    }

    public Ticket parkVehicle(Vehicle v) {
        for (Floor floor : floors) {
            ParkingSpot spot = floor.findAndAssign(v);
            if (spot != null) {
                String id = "T-" + seq.incrementAndGet();
                Ticket ticket = new Ticket(id, v, spot, floor.getLevel(), Instant.now());
                tickets.put(id, ticket);
                return ticket;
            }
        }
        throw new LotFullException();
    }

    public double unparkVehicle(String ticketId, Instant exitTime) {
        Ticket ticket = tickets.remove(ticketId);
        if (ticket == null) {
            throw new TicketNotFoundException();
        }
        for (Floor floor : floors) {
            if (floor.getLevel() == ticket.getFloorLevel()) {
                floor.free(ticket.getSpot().getId());
                break;
            }
        }
        return pricing.calculateFee(ticket, exitTime);
    }
}
