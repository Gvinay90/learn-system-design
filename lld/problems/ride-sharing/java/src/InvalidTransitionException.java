public class InvalidTransitionException extends Exception {
    public InvalidTransitionException() {
        super("invalid trip status transition");
    }
}
