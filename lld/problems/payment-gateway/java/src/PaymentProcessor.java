public interface PaymentProcessor {
    void process(PaymentRequest request) throws PaymentProcessingException;

    class PaymentProcessingException extends Exception {
        public PaymentProcessingException(String message) { super(message); }
    }
}
