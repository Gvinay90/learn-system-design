import java.util.List;

/** Drives turn order, dice rolls, and win detection for a single match. */
public class Game {
    public static class NotEnoughPlayersException extends RuntimeException {
        public NotEnoughPlayersException() { super("need at least two players"); }
    }
    public static class GameAlreadyOverException extends RuntimeException {
        public GameAlreadyOverException() { super("game already has a winner"); }
    }

    private final Board board;
    private final Dice dice;
    private final List<Player> players;
    private int turn = 0;
    private Player winner;

    public Game(Board board, Dice dice, List<String> playerNames) {
        if (playerNames.size() < 2) {
            throw new NotEnoughPlayersException();
        }
        this.board = board;
        this.dice = dice;
        this.players = playerNames.stream().map(Player::new).collect(java.util.stream.Collectors.toList());
    }

    /**
     * Rolls the dice for the current player, applies the exact-landing rule and any
     * snake/ladder at the destination, then rotates turn order. Returns the player
     * who just moved.
     */
    public Player playTurn() {
        if (winner != null) {
            throw new GameAlreadyOverException();
        }

        Player player = players.get(turn);
        int roll = dice.roll();
        int target = player.getPosition() + roll;

        // Overshooting the last cell is not a legal move: the player stays put.
        if (target <= board.getSize()) {
            player.setPosition(board.resolve(target));
        }

        if (player.getPosition() == board.getSize()) {
            winner = player;
        } else {
            turn = (turn + 1) % players.size();
        }
        return player;
    }

    public Player getCurrentPlayer() {
        return players.get(turn);
    }

    public Player getWinner() {
        return winner;
    }

    /** Runs turns until a winner emerges, guarding against runaway loops. */
    public Player play(int maxTurns) {
        for (int i = 0; i < maxTurns; i++) {
            playTurn();
            if (winner != null) {
                return winner;
            }
        }
        throw new IllegalStateException("no winner after " + maxTurns + " turns");
    }
}
