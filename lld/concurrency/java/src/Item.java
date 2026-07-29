/** A unit of work flowing from a producer to a worker pool. */
public class Item {
    private final int producerId;
    private final int seq;

    public Item(int producerId, int seq) {
        this.producerId = producerId;
        this.seq = seq;
    }

    public String id() {
        return "p" + producerId + "-" + seq;
    }
}
