public class PushNotification implements Notification {
    @Override
    public String send(String recipient, String message) {
        return "Push to " + recipient + ": " + message;
    }
}
