public class Screen {
    private boolean down;

    public void lower() { down = true; }
    public void raise() { down = false; }

    public boolean isDown() { return down; }
}
