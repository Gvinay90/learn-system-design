public class Job {
    private final String id;
    private final int priority;
    private volatile long runAtMillis;
    private final Task task;
    private final int maxRetries;
    private volatile RetryPolicy retryPolicy;
    private volatile int attempts;

    public Job(String id, int priority, long runAtMillis, Task task, int maxRetries, RetryPolicy retryPolicy) {
        this.id = id;
        this.priority = priority;
        this.runAtMillis = runAtMillis;
        this.task = task;
        this.maxRetries = maxRetries;
        this.retryPolicy = retryPolicy;
    }

    public Job(String id, int priority, long runAtMillis, Task task) {
        this(id, priority, runAtMillis, task, 0, null);
    }

    public String getId() {
        return id;
    }

    public int getPriority() {
        return priority;
    }

    public long getRunAtMillis() {
        return runAtMillis;
    }

    public void setRunAtMillis(long runAtMillis) {
        this.runAtMillis = runAtMillis;
    }

    public Task getTask() {
        return task;
    }

    public int getMaxRetries() {
        return maxRetries;
    }

    public RetryPolicy getRetryPolicy() {
        return retryPolicy;
    }

    public void setRetryPolicy(RetryPolicy retryPolicy) {
        this.retryPolicy = retryPolicy;
    }

    public int getAttempts() {
        return attempts;
    }

    public void incrementAttempts() {
        attempts++;
    }
}
