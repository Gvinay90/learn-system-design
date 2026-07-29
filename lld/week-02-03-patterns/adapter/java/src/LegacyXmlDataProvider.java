import java.util.Map;

public class LegacyXmlDataProvider {
    private final Map<String, String> store;

    public LegacyXmlDataProvider(Map<String, String> store) {
        this.store = store;
    }

    public String fetchXml(String id) throws RecordNotFoundException {
        String value = store.get(id);
        if (value == null) {
            throw new RecordNotFoundException(id);
        }
        return "<record id=\"" + id + "\">" + value + "</record>";
    }
}
