public class BankAccountProxyTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testOwnerCanWithdraw();
        testNonOwnerCannotWithdraw();
        testNonOwnerCannotDeposit();
        testInsufficientFunds();
        testAccessLogRecordsAllAttempts();
        System.out.println("All BankAccountProxyTest cases passed.");
    }

    private static void testOwnerCanWithdraw() {
        BankAccountProxy acc = new BankAccountProxy(new RealBankAccount("alice", 100));
        acc.withdraw("alice", 40);
        assertEquals(60.0, acc.getBalance(), "balance after owner withdraw");
    }

    private static void testNonOwnerCannotWithdraw() {
        BankAccountProxy acc = new BankAccountProxy(new RealBankAccount("alice", 100));
        try {
            acc.withdraw("mallory", 10);
            throw new AssertionError("expected UnauthorizedException");
        } catch (Account.UnauthorizedException e) {
            // expected
        }
        assertEquals(100.0, acc.getBalance(), "balance unchanged");
    }

    private static void testNonOwnerCannotDeposit() {
        BankAccountProxy acc = new BankAccountProxy(new RealBankAccount("alice", 100));
        try {
            acc.deposit("mallory", 500);
            throw new AssertionError("expected UnauthorizedException");
        } catch (Account.UnauthorizedException e) {
            // expected
        }
        assertEquals(100.0, acc.getBalance(), "balance unchanged");
    }

    private static void testInsufficientFunds() {
        BankAccountProxy acc = new BankAccountProxy(new RealBankAccount("alice", 20));
        try {
            acc.withdraw("alice", 50);
            throw new AssertionError("expected InsufficientFundsException");
        } catch (Account.InsufficientFundsException e) {
            // expected
        }
    }

    private static void testAccessLogRecordsAllAttempts() {
        BankAccountProxy acc = new BankAccountProxy(new RealBankAccount("alice", 100));
        acc.withdraw("alice", 10);
        try {
            acc.deposit("mallory", 5);
        } catch (Account.UnauthorizedException ignored) {
        }
        assertEquals(2, acc.getAccessLog().size(), "access log entries");
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
