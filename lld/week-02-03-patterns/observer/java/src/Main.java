public class Main {
    public static void main(String[] args) {
        WeatherStation station = new WeatherStation();
        PhoneDisplay phone = new PhoneDisplay();
        LoggingDisplay logger = new LoggingDisplay();
        station.subscribe(phone);
        station.subscribe(logger);

        station.setMeasurements(21.5, 60);
        System.out.println("Phone display: " + phone.getLastTempC() + "C, " + phone.getLastHumidity() + "%");

        WeatherStationTest.runAll();
    }
}
