public class Projector {
    private boolean on;
    private String input = "";

    public void turnOn() { on = true; }
    public void turnOff() { on = false; input = ""; }
    public void setInput(String i) { input = i; }

    public boolean isOn() { return on; }
    public String getInput() { return input; }
}
