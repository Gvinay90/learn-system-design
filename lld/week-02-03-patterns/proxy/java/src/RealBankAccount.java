/**
 * The expensive/sensitive real subject. It performs no authorization checks
 * of its own — that's the proxy's job.
 */
public class RealBankAccount implements Account {
    private final String owner;
    private double balance;

    public RealBankAccount(String owner, double initialBalance) {
        this.owner = owner;
        this.balance = initialBalance;
    }

    public String getOwner() { return owner; }

    @Override
    public void deposit(String requester, double amount) {
        balance += amount;
    }

    @Override
    public void withdraw(String requester, double amount) {
        if (amount > balance) {
            throw new InsufficientFundsException("insufficient funds");
        }
        balance -= amount;
    }

    @Override
    public double getBalance() { return balance; }
}
