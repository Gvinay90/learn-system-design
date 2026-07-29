public class KnightPiece extends AbstractPiece {
    public KnightPiece(Color color) {
        super(color);
    }

    @Override
    public PieceType getType() {
        return PieceType.KNIGHT;
    }

    @Override
    public boolean canMove(Board board, Position from, Position to) {
        int dr = abs(to.row - from.row);
        int dc = abs(to.col - from.col);
        if (!((dr == 1 && dc == 2) || (dr == 2 && dc == 1))) {
            return false;
        }
        return targetOK(board, to);
    }
}
