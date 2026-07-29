import java.util.LinkedHashMap;
import java.util.Map;

public class ModernDataProvider implements DataProvider {
    private final Map<String, String> store;

    public ModernDataProvider(Map<String, String> store) {
        this.store = store;
    }

    @Override
    public Map<String, String> fetchJson(String id) throws RecordNotFoundException {
        String value = store.get(id);
        if (value == null) {
            throw new RecordNotFoundException(id);
        }
        Map<String, String> result = new LinkedHashMap<>();
        result.put("id", id);
        result.put("value", value);
        return result;
    }
}
