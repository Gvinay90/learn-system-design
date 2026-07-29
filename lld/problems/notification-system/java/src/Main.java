import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Main {
    public static void main(String[] args) {
        EmailChannel email = new EmailChannel();
        SMSChannel sms = new SMSChannel();
        PushChannel push = new PushChannel();

        NotificationService service = new NotificationService(RetryPolicy.defaultPolicy());
        service.registerChannel(email);
        service.registerChannel(sms);
        service.registerChannel(push);

        service.setPreferences("u1", Channel.EMAIL, Channel.SMS);

        Map<String, String> data = new HashMap<>();
        data.put("name", "Ann");
        data.put("orderId", "42");

        List<SendResult> results = service.notify("u1", "ann@example.com", "Hello {name}, your order {orderId} shipped", data);
        for (SendResult r : results) {
            System.out.println(r.channel + " -> success=" + r.success() + " attempts=" + r.attempts);
        }
        System.out.println("email sent: " + email.sent().size());
        System.out.println("sms sent: " + sms.sent().size());

        NotificationSystemTest.runAll();
    }
}
