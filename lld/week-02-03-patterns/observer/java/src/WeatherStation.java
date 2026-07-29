import java.util.ArrayList;
import java.util.List;

/**
 * The subject: holds the current readings and notifies all subscribed
 * observers whenever they change.
 */
public class WeatherStation {
    private final List<Observer> observers = new ArrayList<>();
    private double tempC;
    private double humidity;

    public void subscribe(Observer o) {
        observers.add(o);
    }

    public void unsubscribe(Observer o) {
        observers.remove(o);
    }

    public void setMeasurements(double tempC, double humidity) {
        this.tempC = tempC;
        this.humidity = humidity;
        for (Observer o : observers) {
            o.update(tempC, humidity);
        }
    }
}
