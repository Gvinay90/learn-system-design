import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

public class Main {
    public static void main(String[] args) {
        Set<String> consumed = ConcurrentHashMap.newKeySet();
        BoundedPipeline pipeline = new BoundedPipeline(4, 8);
        pipeline.run(3, 20, item -> consumed.add(item.id()));

        System.out.println("Consumed " + consumed.size() + " unique items via producer-consumer pipeline");

        ConcurrencyTest.runAll();
    }
}
