/** Shared base for all concrete pieces: stores color and a target-square helper. */
public abstract class AbstractPiece implements Piece {
    protected final Color color;

    protected AbstractPiece(Color color) {
        this.color = color;
    }

    @Override
    public Color getColor() {
        return color;
    }

    /** True if `to` is empty or occupied by the opposing color (capturable). */
    protected boolean targetOK(Board board, Position to) {
        Piece occupant = board.at(to);
        return occupant == null || occupant.getColor() != color;
    }

    protected static int abs(int n) {
        return n < 0 ? -n : n;
    }
}
