import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class NotificationSystemTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testRenderTemplate();
        testNotifyDispatchesPerPreference();
        testNotifyUnknownUserThrows();
        testRetrySucceedsAfterNFailures();
        testRetryGivesUpAfterMaxAttempts();
        testNotifyMultipleChannelsIndependentRetry();
        System.out.println("All NotificationSystemTest cases passed.");
    }

    private static void testRenderTemplate() {
        Map<String, String> data = new HashMap<>();
        data.put("name", "Alice");
        assertEquals("Hello Alice", TemplateRenderer.render("Hello {name}", data), "single placeholder");

        data.put("orderId", "42");
        assertEquals("Hello Alice, your order 42 shipped",
                TemplateRenderer.render("Hello {name}, your order {orderId} shipped", data),
                "multiple placeholders");

        Map<String, String> partial = new HashMap<>();
        partial.put("name", "Carl");
        assertEquals("Hi Carl, code {code}", TemplateRenderer.render("Hi {name}, code {code}", partial), "missing key left untouched");

        assertEquals("plain message", TemplateRenderer.render("plain message", new HashMap<>()), "no placeholders");
        assertEquals("Hello {name", TemplateRenderer.render("Hello {name", data), "unterminated brace");
    }

    private static void testNotifyDispatchesPerPreference() {
        EmailChannel email = new EmailChannel();
        SMSChannel sms = new SMSChannel();
        PushChannel push = new PushChannel();

        NotificationService service = new NotificationService(new RetryPolicy(1, 0));
        service.registerChannel(email);
        service.registerChannel(sms);
        service.registerChannel(push);

        service.setPreferences("u1", Channel.EMAIL, Channel.SMS);
        service.setPreferences("u2", Channel.PUSH);

        Map<String, String> data = new HashMap<>();
        data.put("name", "Ann");
        List<SendResult> results = service.notify("u1", "u1@example.com", "Hello {name}", data);

        if (results.size() != 2) {
            throw new AssertionError("expected 2 results for u1, got " + results.size());
        }
        for (SendResult r : results) {
            if (!r.success()) {
                throw new AssertionError("unexpected send error on channel " + r.channel);
            }
        }
        assertEquals(1, email.sent().size(), "email sent count");
        assertEquals("Hello Ann", email.sent().get(0).message, "email rendered message");
        assertEquals(1, sms.sent().size(), "sms sent count");
        assertEquals(0, push.sent().size(), "push should not receive u1's message");

        Map<String, String> data2 = new HashMap<>();
        data2.put("name", "Bo");
        service.notify("u2", "device-token", "Ping {name}", data2);
        assertEquals(1, push.sent().size(), "push sent count for u2");
        assertEquals("Ping Bo", push.sent().get(0).message, "push rendered message");
    }

    private static void testNotifyUnknownUserThrows() {
        NotificationService service = new NotificationService(RetryPolicy.defaultPolicy());
        service.registerChannel(new EmailChannel());

        try {
            service.notify("ghost", "x", "hi", new HashMap<>());
            throw new AssertionError("expected IllegalArgumentException for unknown user");
        } catch (IllegalArgumentException expected) {
            // expected
        }
    }

    private static void testRetrySucceedsAfterNFailures() {
        EmailChannel underlying = new EmailChannel();
        FlakyChannel flaky = new FlakyChannel(underlying, 2);

        NotificationService service = new NotificationService(new RetryPolicy(3, 1));
        service.registerChannel(flaky);
        service.setPreferences("u1", Channel.EMAIL);

        Map<String, String> data = new HashMap<>();
        data.put("orderId", "7");
        List<SendResult> results = service.notify("u1", "u1@example.com", "Order {orderId} shipped", data);

        assertEquals(1, results.size(), "result count");
        SendResult r = results.get(0);
        if (!r.success()) {
            throw new AssertionError("expected eventual success, got err: " + r.error);
        }
        assertEquals(3, r.attempts, "attempts (2 failures + 1 success)");
        assertEquals(3, flaky.attempts(), "channel invocation count");
        assertEquals(1, underlying.sent().size(), "underlying successful sends");
        assertEquals("Order 7 shipped", underlying.sent().get(0).message, "underlying message content");
    }

    private static void testRetryGivesUpAfterMaxAttempts() {
        AlwaysFailChannel failing = new AlwaysFailChannel(Channel.SMS);

        NotificationService service = new NotificationService(new RetryPolicy(4, 1));
        service.registerChannel(failing);
        service.setPreferences("u1", Channel.SMS);

        Map<String, String> data = new HashMap<>();
        data.put("name", "X");
        List<SendResult> results = service.notify("u1", "5551234", "hi {name}", data);

        assertEquals(1, results.size(), "result count");
        SendResult r = results.get(0);
        if (r.success()) {
            throw new AssertionError("expected send to ultimately fail");
        }
        assertEquals(4, r.attempts, "attempts equal to MaxAttempts");
        assertEquals(4, failing.attempts(), "channel invocation count");
    }

    private static void testNotifyMultipleChannelsIndependentRetry() {
        EmailChannel email = new EmailChannel();
        FlakyChannel flakySms = new FlakyChannel(new SMSChannel(), 1);

        NotificationService service = new NotificationService(new RetryPolicy(2, 1));
        service.registerChannel(email);
        service.registerChannel(flakySms);
        service.setPreferences("u1", Channel.EMAIL, Channel.SMS);

        List<SendResult> results = service.notify("u1", "recipient", "msg", new HashMap<>());
        assertEquals(2, results.size(), "result count");
        for (SendResult r : results) {
            if (!r.success()) {
                throw new AssertionError("expected both channels to eventually succeed, channel " + r.channel + " failed: " + r.error);
            }
        }
        assertEquals(1, results.get(0).attempts, "email succeeds on first attempt");
        assertEquals(2, results.get(1).attempts, "sms succeeds on second attempt");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
