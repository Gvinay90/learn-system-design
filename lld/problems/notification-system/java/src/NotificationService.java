import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Dispatches rendered notifications to a user's preferred channels,
 * retrying each channel independently on failure.
 */
public class NotificationService {
    private final Map<Channel, Notifier> channels = new HashMap<>();
    private final Map<String, List<Channel>> preferences = new HashMap<>();
    private final RetryPolicy retry;

    public NotificationService(RetryPolicy retry) {
        this.retry = retry;
    }

    /** Makes a Notifier available for dispatch under its own channel() identity. */
    public synchronized void registerChannel(Notifier notifier) {
        channels.put(notifier.channel(), notifier);
    }

    /** Records which channels a given user wants to receive notifications on, in order. */
    public synchronized void setPreferences(String userId, Channel... channels) {
        preferences.put(userId, new ArrayList<>(Arrays.asList(channels)));
    }

    /**
     * Renders template with data and sends the result to every channel
     * preferred by userId, retrying each channel per the service's
     * RetryPolicy. Returns one SendResult per attempted channel and does not
     * stop early if one channel ultimately fails.
     *
     * @throws IllegalArgumentException if userId has no registered preferences.
     */
    public List<SendResult> notify(String userId, String recipient, String template, Map<String, String> data) {
        List<Channel> prefs;
        synchronized (this) {
            prefs = preferences.get(userId);
            if (prefs == null) {
                throw new IllegalArgumentException("no channel preferences registered for user " + userId);
            }
            prefs = new ArrayList<>(prefs);
        }

        String message = TemplateRenderer.render(template, data);

        List<SendResult> results = new ArrayList<>();
        for (Channel ch : prefs) {
            Notifier notifier;
            synchronized (this) {
                notifier = channels.get(ch);
            }
            if (notifier == null) {
                results.add(new SendResult(ch, 0, new SendFailedException("no channel registered for " + ch)));
                continue;
            }
            results.add(sendWithRetry(notifier, recipient, message));
        }
        return results;
    }

    private SendResult sendWithRetry(Notifier notifier, String recipient, String message) {
        int maxAttempts = Math.max(retry.maxAttempts, 1);

        Exception lastError = null;
        for (int attempt = 1; attempt <= maxAttempts; attempt++) {
            try {
                notifier.send(recipient, message);
                return new SendResult(notifier.channel(), attempt, null);
            } catch (SendFailedException e) {
                lastError = e;
                if (attempt < maxAttempts && retry.delayMillis > 0) {
                    try {
                        Thread.sleep(retry.delayMillis);
                    } catch (InterruptedException ie) {
                        Thread.currentThread().interrupt();
                    }
                }
            }
        }
        return new SendResult(notifier.channel(), maxAttempts, lastError);
    }
}
