import java.util.ArrayList;
import java.util.List;

/** Orchestrates turn alternation and rule enforcement on top of Board. */
public class Game {
    public static class NotYourTurnException extends RuntimeException {
        public NotYourTurnException() { super("not your turn"); }
    }
    public static class NoPieceAtSourceException extends RuntimeException {
        public NoPieceAtSourceException() { super("no piece at source square"); }
    }
    public static class OwnPieceAtTargetException extends RuntimeException {
        public OwnPieceAtTargetException() { super("cannot capture your own piece"); }
    }
    public static class IllegalMoveException extends RuntimeException {
        public IllegalMoveException() { super("illegal move for this piece"); }
    }
    public static class MoveLeavesKingInCheckException extends RuntimeException {
        public MoveLeavesKingInCheckException() { super("move would leave your own king in check"); }
    }

    private final Board board;
    private Color turn = Color.WHITE;
    private final List<Move> history = new ArrayList<>();

    public Game() {
        this.board = Board.newGame();
    }

    public Game(Board board, Color turn) {
        this.board = board;
        this.turn = turn;
    }

    public Board getBoard() {
        return board;
    }

    public Color getTurn() {
        return turn;
    }

    public List<Move> getHistory() {
        return history;
    }

    /** Validates and applies a move, then switches turns. */
    public void move(Position from, Position to) {
        Piece piece = board.at(from);
        if (piece == null) {
            throw new NoPieceAtSourceException();
        }
        if (piece.getColor() != turn) {
            throw new NotYourTurnException();
        }
        Piece target = board.at(to);
        if (target != null && target.getColor() == piece.getColor()) {
            throw new OwnPieceAtTargetException();
        }
        if (!piece.canMove(board, from, to)) {
            throw new IllegalMoveException();
        }

        Board trial = board.clone();
        trial.move(from, to);
        if (trial.isInCheck(piece.getColor())) {
            throw new MoveLeavesKingInCheckException();
        }

        Piece captured = board.move(from, to);
        history.add(new Move(from, to, piece, captured));
        turn = turn.opposite();
    }

    public boolean isInCheck(Color color) {
        return board.isInCheck(color);
    }

    /**
     * Reports whether color is in check with no legal move (by any of its
     * pieces, to any square) that escapes check. This brute-forces every
     * (piece, destination) pair rather than generating a minimal legal-move
     * set, which is intentionally simple at interview scope.
     */
    public boolean isCheckmate(Color color) {
        if (!board.isInCheck(color)) {
            return false;
        }
        for (int r = 0; r < 8; r++) {
            for (int c = 0; c < 8; c++) {
                Piece piece = board.rawAt(r, c);
                if (piece == null || piece.getColor() != color) {
                    continue;
                }
                Position from = new Position(r, c);
                for (int tr = 0; tr < 8; tr++) {
                    for (int tc = 0; tc < 8; tc++) {
                        Position to = new Position(tr, tc);
                        if (!piece.canMove(board, from, to)) {
                            continue;
                        }
                        Board trial = board.clone();
                        trial.move(from, to);
                        if (!trial.isInCheck(color)) {
                            return false;
                        }
                    }
                }
            }
        }
        return true;
    }
}
