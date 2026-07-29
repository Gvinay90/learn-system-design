public class Amplifier {
    private boolean on;
    private int volume;

    public void turnOn() { on = true; }
    public void turnOff() { on = false; volume = 0; }
    public void setVolume(int v) { volume = v; }

    public boolean isOn() { return on; }
    public int getVolume() { return volume; }
}
