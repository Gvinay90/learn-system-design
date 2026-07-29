/**
 * Replays a fixed sequence of rolls, cycling once exhausted. This is what makes
 * the game engine testable without flakiness: no real randomness in the test suite.
 */
public class ScriptedDice implements Dice {
    private final int[] rolls;
    private int pos = 0;

    public ScriptedDice(int... rolls) {
        this.rolls = rolls;
    }

    @Override
    public int roll() {
        if (rolls.length == 0) {
            return 1;
        }
        int v = rolls[pos % rolls.length];
        pos++;
        return v;
    }
}
