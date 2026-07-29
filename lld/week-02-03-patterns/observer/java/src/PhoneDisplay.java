public class PhoneDisplay implements Observer {
    private double lastTempC;
    private double lastHumidity;
    private int updateCount;

    @Override
    public void update(double tempC, double humidity) {
        lastTempC = tempC;
        lastHumidity = humidity;
        updateCount++;
    }

    public double getLastTempC() { return lastTempC; }
    public double getLastHumidity() { return lastHumidity; }
    public int getUpdateCount() { return updateCount; }
}
