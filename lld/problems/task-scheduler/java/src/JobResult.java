public class JobResult {
    private final String jobId;
    private Status status;
    private int attempts;
    private Exception error;

    public JobResult(String jobId, Status status, int attempts, Exception error) {
        this.jobId = jobId;
        this.status = status;
        this.attempts = attempts;
        this.error = error;
    }

    public JobResult copy() {
        return new JobResult(jobId, status, attempts, error);
    }

    public String getJobId() {
        return jobId;
    }

    public Status getStatus() {
        return status;
    }

    public void setStatus(Status status) {
        this.status = status;
    }

    public int getAttempts() {
        return attempts;
    }

    public void setAttempts(int attempts) {
        this.attempts = attempts;
    }

    public Exception getError() {
        return error;
    }

    public void setError(Exception error) {
        this.error = error;
    }

    @Override
    public String toString() {
        return "JobResult{jobId='" + jobId + "', status=" + status + ", attempts=" + attempts
                + ", error=" + (error == null ? null : error.getMessage()) + "}";
    }
}
