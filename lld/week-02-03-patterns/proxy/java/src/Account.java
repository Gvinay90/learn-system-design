/**
 * Common interface shared by the real subject and its proxy, so callers can
 * use either interchangeably.
 */
public interface Account {
    void deposit(String requester, double amount);
    void withdraw(String requester, double amount);
    double getBalance();

    class UnauthorizedException extends RuntimeException {
        public UnauthorizedException(String message) { super(message); }
    }

    class InsufficientFundsException extends RuntimeException {
        public InsufficientFundsException(String message) { super(message); }
    }
}
