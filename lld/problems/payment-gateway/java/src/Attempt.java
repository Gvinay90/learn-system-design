import java.time.Instant;

public class Attempt {
    private final int number;
    private final boolean success;
    private final String error;
    private final Instant at;

    public Attempt(int number, boolean success, String error, Instant at) {
        this.number = number;
        this.success = success;
        this.error = error;
        this.at = at;
    }

    public int getNumber() { return number; }
    public boolean isSuccess() { return success; }
    public String getError() { return error; }
    public Instant getAt() { return at; }
}
