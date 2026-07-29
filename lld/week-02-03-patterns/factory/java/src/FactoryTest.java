/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out FactoryTest` directly.
 */
public class FactoryTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testCreateEmailNotification();
        testCreateSMSNotification();
        testCreatePushNotification();
        testCreateUnknownNotificationType();
        System.out.println("All FactoryTest cases passed.");
    }

    private static void testCreateEmailNotification() {
        Notification n = NotificationFactory.create(NotificationType.EMAIL);
        assertEquals("Email to alice@example.com: hello", n.send("alice@example.com", "hello"), "email format");
    }

    private static void testCreateSMSNotification() {
        Notification n = NotificationFactory.create(NotificationType.SMS);
        assertEquals("SMS to +1-555-0100: hello", n.send("+1-555-0100", "hello"), "sms format");
    }

    private static void testCreatePushNotification() {
        Notification n = NotificationFactory.create(NotificationType.PUSH);
        assertEquals("Push to device-123: hello", n.send("device-123", "hello"), "push format");
    }

    private static void testCreateUnknownNotificationType() {
        try {
            NotificationFactory.create(null);
            throw new AssertionError("expected UnknownNotificationTypeException");
        } catch (NotificationFactory.UnknownNotificationTypeException e) {
            // expected
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
