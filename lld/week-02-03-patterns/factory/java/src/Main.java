public class Main {
    public static void main(String[] args) {
        Notification email = NotificationFactory.create(NotificationType.EMAIL);
        System.out.println(email.send("alice@example.com", "your order has shipped"));

        Notification sms = NotificationFactory.create(NotificationType.SMS);
        System.out.println(sms.send("+1-555-0100", "your order has shipped"));

        Notification push = NotificationFactory.create(NotificationType.PUSH);
        System.out.println(push.send("device-123", "your order has shipped"));

        FactoryTest.runAll();
    }
}
