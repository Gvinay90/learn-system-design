import java.util.List;
import java.util.Map;

public class Main {
    public static void main(String[] args) {
        Document original = new Document(
                "Q3 Report",
                new Metadata("Vinay", List.of("finance", "draft")),
                List.of("Intro", "Numbers"),
                Map.of("status", "draft"));

        Document clone = original.deepClone();
        clone.getMeta().setAuthor("Someone Else");
        clone.getMeta().getTags().add("reviewed");
        clone.getSections().add("Appendix");
        clone.getProps().put("status", "final");

        System.out.println("Original author: " + original.getMeta().getAuthor());
        System.out.println("Clone author: " + clone.getMeta().getAuthor());
        System.out.println("Original tags: " + original.getMeta().getTags());
        System.out.println("Clone tags: " + clone.getMeta().getTags());
        System.out.println("Original sections: " + original.getSections());
        System.out.println("Clone sections: " + clone.getSections());
        System.out.println("Original props: " + original.getProps());
        System.out.println("Clone props: " + clone.getProps());

        PrototypeTest.runAll();
    }
}
