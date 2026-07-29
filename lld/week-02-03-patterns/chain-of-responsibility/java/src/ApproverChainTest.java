public class ApproverChainTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testManagerApprovesSmallExpense();
        testDirectorApprovesMidExpense();
        testVpApprovesLargeExpense();
        testNoApproverForExcessiveExpense();
        testChainOrderMatters();
        System.out.println("All ApproverChainTest cases passed.");
    }

    private static Approver newTestChain() {
        return ApproverChain.build(new Manager(1000), new Director(5000), new VP(20000));
    }

    private static void testManagerApprovesSmallExpense() {
        Approver chain = newTestChain();
        String who = chain.approve(new ExpenseRequest(500, "team lunch"));
        assertEquals("Manager", who, "small expense approver");
    }

    private static void testDirectorApprovesMidExpense() {
        Approver chain = newTestChain();
        String who = chain.approve(new ExpenseRequest(3000, "conference"));
        assertEquals("Director", who, "mid expense approver");
    }

    private static void testVpApprovesLargeExpense() {
        Approver chain = newTestChain();
        String who = chain.approve(new ExpenseRequest(15000, "new hire equipment"));
        assertEquals("VP", who, "large expense approver");
    }

    private static void testNoApproverForExcessiveExpense() {
        Approver chain = newTestChain();
        try {
            chain.approve(new ExpenseRequest(100000, "private jet"));
            throw new AssertionError("expected NoApproverException");
        } catch (Approver.NoApproverException e) {
            // expected
        }
    }

    private static void testChainOrderMatters() {
        Approver chain = ApproverChain.build(new Manager(1000));
        try {
            chain.approve(new ExpenseRequest(2000, "no further approver"));
            throw new AssertionError("expected NoApproverException");
        } catch (Approver.NoApproverException e) {
            // expected
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
