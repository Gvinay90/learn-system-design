public class SMSNotification implements Notification {
    @Override
    public String send(String recipient, String message) {
        return "SMS to " + recipient + ": " + message;
    }
}
