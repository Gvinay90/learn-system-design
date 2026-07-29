public class Main {
    public static void main(String[] args) {
        Order order = new Order("O-1");
        System.out.println("Initial state: " + order.getState().name());
        order.pay();
        order.ship();
        order.deliver();
        System.out.println("Final state: " + order.getState().name());

        OrderTest.runAll();
    }
}
