public class Main {
    public static void main(String[] args) {
        LRUCache<String, Integer> cache = new LRUCache<>(2);
        cache.put("a", 1);
        cache.put("b", 2);
        System.out.println("get(a) = " + cache.get("a"));
        cache.put("c", 3);
        System.out.println("get(b) after eviction = " + cache.get("b"));
        System.out.println("get(c) = " + cache.get("c"));

        LRUCacheTest.runAll();
    }
}
