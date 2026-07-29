import java.util.List;

public class Main {
    public static void main(String[] args) {
        Broker broker = new Broker();
        broker.createTopic("orders", 3);

        for (String v : new String[] {"a", "b", "c"}) {
            Broker.ProduceResult r = broker.produce("orders", "k1", v);
            System.out.println("produced " + v + " -> partition " + r.partitionId + " offset " + r.offset);
        }

        List<Message> messages = broker.consume("group-1", "orders", 0, 10);
        System.out.print("group-1 read: ");
        for (Message m : messages) {
            System.out.print("[" + m.getOffset() + ":" + m.getValue() + "] ");
        }
        System.out.println();

        List<Message> again = broker.consume("group-1", "orders", 0, 10);
        System.out.println("group-1 read again (after auto-commit): " + again.size() + " messages");

        BrokerTest.runAll();
    }
}
