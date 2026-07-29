public class DvdPlayer {
    private boolean on;
    private String movie = "";

    public void turnOn() { on = true; }
    public void turnOff() { on = false; movie = ""; }
    public void play(String m) { movie = m; }
    public void stop() { movie = ""; }

    public boolean isOn() { return on; }
    public String getMovie() { return movie; }
}
