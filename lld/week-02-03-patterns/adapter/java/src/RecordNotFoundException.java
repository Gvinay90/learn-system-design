public class RecordNotFoundException extends Exception {
    public RecordNotFoundException(String id) {
        super("record not found: " + id);
    }
}
