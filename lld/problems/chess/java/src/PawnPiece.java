public class PawnPiece extends AbstractPiece {
    public PawnPiece(Color color) {
        super(color);
    }

    @Override
    public PieceType getType() {
        return PieceType.PAWN;
    }

    @Override
    public boolean canMove(Board board, Position from, Position to) {
        int dir = color == Color.WHITE ? -1 : 1;
        int startRow = color == Color.WHITE ? 6 : 1;

        int dr = to.row - from.row;
        int dc = to.col - from.col;

        if (dc == 0) {
            if (dr == dir) {
                return board.at(to) == null;
            }
            if (dr == 2 * dir && from.row == startRow) {
                Position mid = new Position(from.row + dir, from.col);
                return board.at(mid) == null && board.at(to) == null;
            }
            return false;
        }
        if (abs(dc) == 1 && dr == dir) {
            Piece occupant = board.at(to);
            return occupant != null && occupant.getColor() != color;
        }
        return false;
    }
}
