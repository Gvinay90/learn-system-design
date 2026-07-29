public class Main {
    public static void main(String[] args) {
        Approver chain = ApproverChain.build(new Manager(1000), new Director(5000), new VP(20000));

        String who = chain.approve(new ExpenseRequest(3000, "conference"));
        System.out.println("Approved by: " + who);

        ApproverChainTest.runAll();
    }
}
