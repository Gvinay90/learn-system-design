public abstract class Approver {
    public static class NoApproverException extends RuntimeException {
        public NoApproverException(String message) { super(message); }
    }

    private final String name;
    private final double limit;
    private Approver next;

    protected Approver(String name, double limit) {
        this.name = name;
        this.limit = limit;
    }

    public void setNext(Approver next) {
        this.next = next;
    }

    public String approve(ExpenseRequest req) {
        if (req.getAmount() <= limit) {
            return name;
        }
        if (next != null) {
            return next.approve(req);
        }
        throw new NoApproverException("no approver in the chain can approve this amount");
    }
}
