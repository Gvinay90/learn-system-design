public class Lights {
    private boolean dimmed;

    public void dim() { dimmed = true; }
    public void brighten() { dimmed = false; }

    public boolean isDimmed() { return dimmed; }
}
