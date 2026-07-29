import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

public class LRUCache<K, V> {
    private final int capacity;
    private final Map<K, Node<K, V>> items = new HashMap<>();
    private final Node<K, V> head = new Node<>(null, null);
    private final Node<K, V> tail = new Node<>(null, null);

    public LRUCache(int capacity) {
        if (capacity <= 0) {
            throw new IllegalArgumentException("capacity must be positive");
        }
        this.capacity = capacity;
        head.next = tail;
        tail.prev = head;
    }

    private void unlink(Node<K, V> n) {
        n.prev.next = n.next;
        n.next.prev = n.prev;
    }

    private void pushFront(Node<K, V> n) {
        n.prev = head;
        n.next = head.next;
        head.next.prev = n;
        head.next = n;
    }

    public synchronized Optional<V> get(K key) {
        Node<K, V> n = items.get(key);
        if (n == null) {
            return Optional.empty();
        }
        unlink(n);
        pushFront(n);
        return Optional.of(n.value);
    }

    public synchronized void put(K key, V value) {
        Node<K, V> n = items.get(key);
        if (n != null) {
            n.value = value;
            unlink(n);
            pushFront(n);
            return;
        }

        n = new Node<>(key, value);
        items.put(key, n);
        pushFront(n);

        if (items.size() > capacity) {
            Node<K, V> lru = tail.prev;
            unlink(lru);
            items.remove(lru.key);
        }
    }

    public synchronized int size() {
        return items.size();
    }

    synchronized int listLength() {
        int count = 0;
        for (Node<K, V> n = head.next; n != tail; n = n.next) {
            count++;
        }
        return count;
    }
}
