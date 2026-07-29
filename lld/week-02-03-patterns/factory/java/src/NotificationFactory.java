public class NotificationFactory {

    public static class UnknownNotificationTypeException extends RuntimeException {
        public UnknownNotificationTypeException(NotificationType type) {
            super("unknown notification type: " + type);
        }
    }

    public static Notification create(NotificationType type) {
        if (type == null) {
            throw new UnknownNotificationTypeException(null);
        }
        switch (type) {
            case EMAIL:
                return new EmailNotification();
            case SMS:
                return new SMSNotification();
            case PUSH:
                return new PushNotification();
            default:
                throw new UnknownNotificationTypeException(type);
        }
    }
}
