import java.util.Random;

public class StandardDice implements Dice {
    private final Random random;

    public StandardDice(long seed) {
        this.random = new Random(seed);
    }

    @Override
    public int roll() {
        return random.nextInt(6) + 1;
    }
}
