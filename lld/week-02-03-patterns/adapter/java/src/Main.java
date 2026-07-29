import java.util.List;
import java.util.Map;

public class Main {
    public static void main(String[] args) throws Exception {
        LegacyXmlDataProvider legacy = new LegacyXmlDataProvider(Map.of("u1", "42"));
        DataProvider adapted = new XmlToJsonAdapter(legacy);

        Map<String, String> record = adapted.fetchJson("u1");
        System.out.println("Adapted record: " + record);

        int sum = Client.fetchAndSum(adapted, List.of("u1"));
        System.out.println("Sum via adapter: " + sum);

        AdapterTest.runAll();
    }
}
