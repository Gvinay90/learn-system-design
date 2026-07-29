import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.BlockingQueue;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.function.Consumer;

/**
 * Runs a bounded producer-consumer pipeline: producer threads push items
 * into an ArrayBlockingQueue (the bounded buffer) and an ExecutorService
 * worker pool drains it concurrently. A poison-pill sentinel per worker
 * signals shutdown once every producer has finished.
 */
public class BoundedPipeline {
    private static final Item POISON = new Item(-1, -1);

    private final int numWorkers;
    private final int bufferCapacity;

    public BoundedPipeline(int numWorkers, int bufferCapacity) {
        this.numWorkers = numWorkers;
        this.bufferCapacity = bufferCapacity;
    }

    public void run(int numProducers, int itemsPerProducer, Consumer<Item> handle) {
        BlockingQueue<Item> buffer = new ArrayBlockingQueue<>(bufferCapacity);
        ExecutorService workerPool = Executors.newFixedThreadPool(numWorkers);
        CountDownLatch producersDone = new CountDownLatch(numProducers);

        List<Thread> producers = new ArrayList<>();
        for (int p = 0; p < numProducers; p++) {
            final int producerId = p;
            Thread t = new Thread(() -> {
                try {
                    for (int seq = 0; seq < itemsPerProducer; seq++) {
                        buffer.put(new Item(producerId, seq));
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    producersDone.countDown();
                }
            });
            producers.add(t);
            t.start();
        }

        for (int w = 0; w < numWorkers; w++) {
            workerPool.submit(() -> {
                try {
                    while (true) {
                        Item item = buffer.take();
                        if (item == POISON) {
                            return;
                        }
                        handle.accept(item);
                    }
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            });
        }

        try {
            producersDone.await();
            for (int w = 0; w < numWorkers; w++) {
                buffer.put(POISON);
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }

        workerPool.shutdown();
        try {
            workerPool.awaitTermination(30, TimeUnit.SECONDS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
