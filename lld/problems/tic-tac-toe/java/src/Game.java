import java.util.List;

/**
 * Orchestrates turn order, move validation and terminal-state tracking on
 * top of a Board. No locking here: unlike the parking lot, a single
 * tic-tac-toe game is played by one turn owner at a time, so there is no
 * concurrent-writer scenario to guard against.
 */
public class Game {
    public static class GameOverException extends RuntimeException {
        public GameOverException() { super("game is already over"); }
    }

    private final Board board;
    private final List<Player> players;
    private int turn = 0;
    private GameStatus status = GameStatus.IN_PROGRESS;
    private Player winner;

    public Game(int size, List<Player> players) {
        this.board = new Board(size);
        this.players = players;
    }

    public Board getBoard() { return board; }
    public GameStatus getStatus() { return status; }
    public Player getWinner() { return winner; }
    public Player getCurrentPlayer() { return players.get(turn); }

    /** Plays a cell for the current player. Returns true if this move wins the game. */
    public boolean move(int row, int col) {
        if (status != GameStatus.IN_PROGRESS) {
            throw new GameOverException();
        }

        Player player = players.get(turn);
        boolean won = board.place(row, col, player.getSymbol());

        if (won) {
            status = GameStatus.WON;
            winner = player;
            return true;
        }
        if (board.isFull()) {
            status = GameStatus.DRAW;
            return false;
        }

        turn = (turn + 1) % players.size();
        return false;
    }
}
