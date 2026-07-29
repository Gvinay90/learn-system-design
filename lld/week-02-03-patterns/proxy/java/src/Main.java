public class Main {
    public static void main(String[] args) {
        RealBankAccount real = new RealBankAccount("alice", 100);
        BankAccountProxy acc = new BankAccountProxy(real);

        acc.withdraw("alice", 40);
        System.out.println("Balance after owner withdraw: " + acc.getBalance());

        try {
            acc.withdraw("mallory", 10);
        } catch (Account.UnauthorizedException e) {
            System.out.println("Blocked: " + e.getMessage());
        }

        BankAccountProxyTest.runAll();
    }
}
