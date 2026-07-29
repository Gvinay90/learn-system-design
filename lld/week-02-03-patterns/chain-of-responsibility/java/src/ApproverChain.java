public class ApproverChain {
    private ApproverChain() {}

    public static Approver build(Approver... approvers) {
        for (int i = 0; i < approvers.length - 1; i++) {
            approvers[i].setNext(approvers[i + 1]);
        }
        return approvers[0];
    }
}
