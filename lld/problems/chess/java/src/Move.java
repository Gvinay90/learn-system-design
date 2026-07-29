/** A recorded move in the game history. */
public class Move {
    public final Position from;
    public final Position to;
    public final Piece piece;
    public final Piece captured;

    public Move(Position from, Position to, Piece piece, Piece captured) {
        this.from = from;
        this.to = to;
        this.piece = piece;
        this.captured = captured;
    }
}
