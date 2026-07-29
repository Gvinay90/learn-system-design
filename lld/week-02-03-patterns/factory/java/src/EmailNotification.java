public class EmailNotification implements Notification {
    @Override
    public String send(String recipient, String message) {
        return "Email to " + recipient + ": " + message;
    }
}
