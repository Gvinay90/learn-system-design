import java.util.Map;

public interface DataProvider {
    Map<String, String> fetchJson(String id) throws RecordNotFoundException;
}
