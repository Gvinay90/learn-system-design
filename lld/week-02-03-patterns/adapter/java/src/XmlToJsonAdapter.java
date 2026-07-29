import java.util.LinkedHashMap;
import java.util.Map;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Adapts LegacyXmlDataProvider (an incompatible interface) to DataProvider,
 * the interface the rest of the system depends on.
 */
public class XmlToJsonAdapter implements DataProvider {
    private static final Pattern RECORD_PATTERN =
            Pattern.compile("<record id=\"(.*?)\">(.*?)</record>");

    private final LegacyXmlDataProvider legacy;

    public XmlToJsonAdapter(LegacyXmlDataProvider legacy) {
        this.legacy = legacy;
    }

    @Override
    public Map<String, String> fetchJson(String id) throws RecordNotFoundException {
        String raw = legacy.fetchXml(id);
        Matcher m = RECORD_PATTERN.matcher(raw);
        if (!m.matches()) {
            throw new RecordNotFoundException(id);
        }
        Map<String, String> result = new LinkedHashMap<>();
        result.put("id", m.group(1));
        result.put("value", m.group(2));
        return result;
    }
}
