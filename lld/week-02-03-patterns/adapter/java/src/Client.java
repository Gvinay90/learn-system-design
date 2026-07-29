import java.util.List;
import java.util.Map;

/** Client depends only on DataProvider, oblivious to legacy vs. native sourcing. */
public class Client {
    public static int fetchAndSum(DataProvider provider, List<String> ids) throws RecordNotFoundException {
        int total = 0;
        for (String id : ids) {
            Map<String, String> record = provider.fetchJson(id);
            total += Integer.parseInt(record.get("value"));
        }
        return total;
    }
}
