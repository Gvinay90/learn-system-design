public class KingPiece extends AbstractPiece {
    public KingPiece(Color color) {
        super(color);
    }

    @Override
    public PieceType getType() {
        return PieceType.KING;
    }

    @Override
    public boolean canMove(Board board, Position from, Position to) {
        int dr = abs(to.row - from.row);
        int dc = abs(to.col - from.col);
        if (dr > 1 || dc > 1 || (dr == 0 && dc == 0)) {
            return false;
        }
        return targetOK(board, to);
    }
}
