/**
 * Implemented by each piece type (Template/Strategy pattern): canMove encodes
 * that piece's movement pattern plus path-clearance, but not whether the move
 * leaves the mover's own king in check — Game is responsible for that.
 */
public interface Piece {
    Color getColor();
    PieceType getType();
    boolean canMove(Board board, Position from, Position to);
}
