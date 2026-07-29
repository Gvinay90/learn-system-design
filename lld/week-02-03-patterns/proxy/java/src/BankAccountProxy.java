import java.util.ArrayList;
import java.util.List;

/**
 * Protection proxy: enforces that only the account owner may deposit or
 * withdraw, and logs every access attempt, without RealBankAccount needing
 * to know about authorization at all.
 */
public class BankAccountProxy implements Account {
    private final RealBankAccount real;
    private final List<String> accessLog = new ArrayList<>();

    public BankAccountProxy(RealBankAccount real) {
        this.real = real;
    }

    @Override
    public void deposit(String requester, double amount) {
        accessLog.add("deposit attempt by " + requester);
        if (!requester.equals(real.getOwner())) {
            throw new UnauthorizedException("requester is not the account owner");
        }
        real.deposit(requester, amount);
    }

    @Override
    public void withdraw(String requester, double amount) {
        accessLog.add("withdraw attempt by " + requester);
        if (!requester.equals(real.getOwner())) {
            throw new UnauthorizedException("requester is not the account owner");
        }
        real.withdraw(requester, amount);
    }

    @Override
    public double getBalance() { return real.getBalance(); }

    public List<String> getAccessLog() { return accessLog; }
}
