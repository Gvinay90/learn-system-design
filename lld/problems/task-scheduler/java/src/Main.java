import java.util.Optional;
import java.util.concurrent.atomic.AtomicInteger;

public class Main {
    public static void main(String[] args) throws InterruptedException {
        Scheduler scheduler = new Scheduler(2);
        scheduler.start();

        AtomicInteger ran = new AtomicInteger();
        scheduler.submit(new Job("demo", 1, System.currentTimeMillis(), () -> {
            ran.incrementAndGet();
            System.out.println("demo job executed");
        }));

        Optional<JobResult> res = scheduler.waitForResult("demo", 1000);
        System.out.println("Result: " + res.orElse(null));

        scheduler.stop();

        TaskSchedulerTest.runAll();
    }
}
