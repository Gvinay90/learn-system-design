public class RookPiece extends AbstractPiece {
    public RookPiece(Color color) {
        super(color);
    }

    @Override
    public PieceType getType() {
        return PieceType.ROOK;
    }

    @Override
    public boolean canMove(Board board, Position from, Position to) {
        int dr = to.row - from.row;
        int dc = to.col - from.col;
        if ((dr == 0) == (dc == 0)) {
            return false;
        }
        return board.clearPath(from, to) && targetOK(board, to);
    }
}
