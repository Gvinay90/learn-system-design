import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.PriorityQueue;

/**
 * In-memory, priority-based job scheduler with delayed execution and
 * retries, running on a fixed worker pool. Jobs are ordered by earliest
 * due time first, tie-broken by higher priority first — mirrors the Go
 * implementation's min-heap semantics.
 */
public class Scheduler {
    private static final RetryPolicy DEFAULT_RETRY_POLICY = new ExponentialBackoff(10, 200);

    private final Object lock = new Object();
    private final PriorityQueue<Job> queue = new PriorityQueue<>(Scheduler::compareJobs);
    private final Map<String, JobResult> results = new HashMap<>();
    private final int workers;
    private final List<Thread> threads = new ArrayList<>();
    private volatile boolean stopped = true;

    public Scheduler(int workers) {
        this.workers = workers;
    }

    private static int compareJobs(Job a, Job b) {
        if (a.getRunAtMillis() == b.getRunAtMillis()) {
            return Integer.compare(b.getPriority(), a.getPriority());
        }
        return Long.compare(a.getRunAtMillis(), b.getRunAtMillis());
    }

    /** Enqueues a job. Safe to call before or after start(), and concurrently. */
    public void submit(Job job) {
        if (job.getRetryPolicy() == null) {
            job.setRetryPolicy(DEFAULT_RETRY_POLICY);
        }
        synchronized (lock) {
            queue.add(job);
            results.put(job.getId(), new JobResult(job.getId(), Status.PENDING, 0, null));
        }
    }

    public void start() {
        stopped = false;
        threads.clear();
        for (int i = 0; i < workers; i++) {
            Thread t = new Thread(this::workerLoop, "scheduler-worker-" + i);
            t.setDaemon(true);
            threads.add(t);
            t.start();
        }
    }

    public void stop() {
        stopped = true;
        for (Thread t : threads) {
            try {
                t.join();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }
    }

    private void workerLoop() {
        while (!stopped) {
            Job job = popDue();
            if (job != null) {
                execute(job);
            } else {
                sleepMillis(2);
            }
        }
    }

    /**
     * Removes and returns the highest-priority job whose runAt has already
     * elapsed, or null if none is ready.
     */
    private Job popDue() {
        synchronized (lock) {
            Job head = queue.peek();
            if (head == null || head.getRunAtMillis() > System.currentTimeMillis()) {
                return null;
            }
            return queue.poll();
        }
    }

    private void execute(Job job) {
        synchronized (lock) {
            results.get(job.getId()).setStatus(Status.RUNNING);
        }

        Exception err = null;
        try {
            job.getTask().run();
        } catch (Exception e) {
            err = e;
        }
        job.incrementAttempts();

        synchronized (lock) {
            JobResult res = results.get(job.getId());
            res.setAttempts(job.getAttempts());
            if (err == null) {
                res.setStatus(Status.SUCCEEDED);
                res.setError(null);
                return;
            }
            res.setError(err);
            if (job.getAttempts() <= job.getMaxRetries()) {
                res.setStatus(Status.RETRYING);
                long delay = job.getRetryPolicy().nextDelayMillis(job.getAttempts() - 1);
                job.setRunAtMillis(System.currentTimeMillis() + delay);
                queue.add(job);
                return;
            }
            res.setStatus(Status.FAILED);
        }
    }

    public Optional<JobResult> getResult(String id) {
        synchronized (lock) {
            JobResult r = results.get(id);
            return r == null ? Optional.empty() : Optional.of(r.copy());
        }
    }

    /**
     * Polls until job id reaches a terminal state (SUCCEEDED or FAILED) or
     * the timeout elapses.
     */
    public Optional<JobResult> waitForResult(String id, long timeoutMillis) {
        long deadline = System.currentTimeMillis() + timeoutMillis;
        while (true) {
            Optional<JobResult> r = getResult(id);
            if (r.isPresent()) {
                Status s = r.get().getStatus();
                if (s == Status.SUCCEEDED || s == Status.FAILED) {
                    return r;
                }
            }
            if (System.currentTimeMillis() > deadline) {
                return Optional.empty();
            }
            sleepMillis(2);
        }
    }

    private static void sleepMillis(long ms) {
        try {
            Thread.sleep(ms);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
