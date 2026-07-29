public class WeatherStationTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testAllObserversNotifiedOnChange();
        testMultipleUpdatesAccumulate();
        testUnsubscribeStopsNotifications();
        testNoObserversDoesNotThrow();
        System.out.println("All WeatherStationTest cases passed.");
    }

    private static void testAllObserversNotifiedOnChange() {
        WeatherStation station = new WeatherStation();
        PhoneDisplay phone = new PhoneDisplay();
        LoggingDisplay logger = new LoggingDisplay();
        station.subscribe(phone);
        station.subscribe(logger);

        station.setMeasurements(21.5, 60);

        assertEquals(21.5, phone.getLastTempC(), "phone temp");
        assertEquals(60.0, phone.getLastHumidity(), "phone humidity");
        assertEquals(1, logger.getHistory().size(), "logger history size");
    }

    private static void testMultipleUpdatesAccumulate() {
        WeatherStation station = new WeatherStation();
        PhoneDisplay phone = new PhoneDisplay();
        station.subscribe(phone);

        station.setMeasurements(20, 50);
        station.setMeasurements(22, 55);

        assertEquals(2, phone.getUpdateCount(), "update count");
        assertEquals(22.0, phone.getLastTempC(), "latest temp");
    }

    private static void testUnsubscribeStopsNotifications() {
        WeatherStation station = new WeatherStation();
        PhoneDisplay phone = new PhoneDisplay();
        station.subscribe(phone);
        station.setMeasurements(20, 50);

        station.unsubscribe(phone);
        station.setMeasurements(99, 99);

        assertEquals(20.0, phone.getLastTempC(), "phone unaffected after unsubscribe");
        assertEquals(1, phone.getUpdateCount(), "only 1 update recorded");
    }

    private static void testNoObserversDoesNotThrow() {
        WeatherStation station = new WeatherStation();
        station.setMeasurements(10, 10);
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
