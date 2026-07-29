public class Main {
    public static void main(String[] args) {
        PaymentGateway gateway = new PaymentGateway(FakePaymentProcessor.succeedsImmediately(), new RetryPolicy(3, 1));

        PaymentRequest request = new PaymentRequest("order-123", "payer-1", "payee-1", 250.0, "INR");
        PaymentResult result = gateway.charge(request);
        System.out.println("Charge result: " + result.getStatus() + " (id=" + result.getId() + ")");

        PaymentResult cached = gateway.charge(request);
        System.out.println("Repeat with same key returns cached id=" + cached.getId());

        System.out.println("Ledger entries: " + gateway.getLedger().getEntries().size());

        PaymentGatewayTest.runAll();
    }
}
