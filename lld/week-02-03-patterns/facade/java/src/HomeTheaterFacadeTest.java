public class HomeTheaterFacadeTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testWatchMovieCoordinatesSubsystems();
        testEndMovieResetsSubsystems();
        System.out.println("All HomeTheaterFacadeTest cases passed.");
    }

    private static void testWatchMovieCoordinatesSubsystems() {
        HomeTheaterFacade theater = new HomeTheaterFacade();
        theater.watchMovie("Inception");

        assertTrue(theater.getLights().isDimmed(), "lights dimmed");
        assertTrue(theater.getScreen().isDown(), "screen down");
        assertTrue(theater.getProjector().isOn(), "projector on");
        assertEquals("dvd", theater.getProjector().getInput(), "projector input");
        assertTrue(theater.getAmp().isOn(), "amp on");
        assertEquals(7, theater.getAmp().getVolume(), "amp volume");
        assertEquals("Inception", theater.getDvd().getMovie(), "movie playing");
    }

    private static void testEndMovieResetsSubsystems() {
        HomeTheaterFacade theater = new HomeTheaterFacade();
        theater.watchMovie("Inception");
        theater.endMovie();

        assertEquals("", theater.getDvd().getMovie(), "dvd movie cleared");
        assertTrue(!theater.getDvd().isOn(), "dvd off");
        assertTrue(!theater.getAmp().isOn(), "amp off");
        assertTrue(!theater.getProjector().isOn(), "projector off");
        assertTrue(!theater.getScreen().isDown(), "screen raised");
        assertTrue(!theater.getLights().isDimmed(), "lights brightened");
    }

    private static void assertTrue(boolean condition, String label) {
        if (!condition) {
            throw new AssertionError(label + ": expected true");
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
