public class QueenPiece extends AbstractPiece {
    public QueenPiece(Color color) {
        super(color);
    }

    @Override
    public PieceType getType() {
        return PieceType.QUEEN;
    }

    @Override
    public boolean canMove(Board board, Position from, Position to) {
        int dr = to.row - from.row;
        int dc = to.col - from.col;
        boolean straight = (dr == 0) != (dc == 0);
        boolean diagonal = dr != 0 && abs(dr) == abs(dc);
        if (!straight && !diagonal) {
            return false;
        }
        return board.clearPath(from, to) && targetOK(board, to);
    }
}
