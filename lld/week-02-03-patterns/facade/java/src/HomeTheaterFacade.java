/**
 * Coordinates five independent subsystems behind two calls, so a client
 * never needs to know the correct on/off sequence for a movie night.
 */
public class HomeTheaterFacade {
    private final Amplifier amp = new Amplifier();
    private final DvdPlayer dvd = new DvdPlayer();
    private final Projector projector = new Projector();
    private final Screen screen = new Screen();
    private final Lights lights = new Lights();

    public void watchMovie(String movie) {
        lights.dim();
        screen.lower();
        projector.turnOn();
        projector.setInput("dvd");
        amp.turnOn();
        amp.setVolume(7);
        dvd.turnOn();
        dvd.play(movie);
    }

    public void endMovie() {
        dvd.stop();
        dvd.turnOff();
        amp.turnOff();
        projector.turnOff();
        screen.raise();
        lights.brighten();
    }

    public Amplifier getAmp() { return amp; }
    public DvdPlayer getDvd() { return dvd; }
    public Projector getProjector() { return projector; }
    public Screen getScreen() { return screen; }
    public Lights getLights() { return lights; }
}
