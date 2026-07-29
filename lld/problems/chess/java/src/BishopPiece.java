public class BishopPiece extends AbstractPiece {
    public BishopPiece(Color color) {
        super(color);
    }

    @Override
    public PieceType getType() {
        return PieceType.BISHOP;
    }

    @Override
    public boolean canMove(Board board, Position from, Position to) {
        int dr = to.row - from.row;
        int dc = to.col - from.col;
        if (abs(dr) != abs(dc) || dr == 0) {
            return false;
        }
        return board.clearPath(from, to) && targetOK(board, to);
    }
}
