import java.util.ArrayList;
import java.util.List;

public class LoggingDisplay implements Observer {
    private final List<String> history = new ArrayList<>();

    @Override
    public void update(double tempC, double humidity) {
        history.add(String.format("temp=%.1f humidity=%.1f", tempC, humidity));
    }

    public List<String> getHistory() { return history; }
}
