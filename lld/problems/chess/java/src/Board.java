/**
 * Holds the 8x8 grid of squares. Clone() deep-copies the grid for
 * check-safety simulation (trial move + IsInCheck) without mutating the
 * real board.
 */
public class Board {
    private final Piece[][] squares = new Piece[8][8];

    public Board() {
    }

    public static Board newGame() {
        Board b = new Board();
        PieceType[] backRank = {
                PieceType.ROOK, PieceType.KNIGHT, PieceType.BISHOP, PieceType.QUEEN,
                PieceType.KING, PieceType.BISHOP, PieceType.KNIGHT, PieceType.ROOK
        };
        for (int col = 0; col < 8; col++) {
            b.squares[0][col] = PieceFactory.create(backRank[col], Color.BLACK);
            b.squares[7][col] = PieceFactory.create(backRank[col], Color.WHITE);
            b.squares[1][col] = PieceFactory.create(PieceType.PAWN, Color.BLACK);
            b.squares[6][col] = PieceFactory.create(PieceType.PAWN, Color.WHITE);
        }
        return b;
    }

    public Piece at(Position pos) {
        return squares[pos.row][pos.col];
    }

    public void set(Position pos, Piece p) {
        squares[pos.row][pos.col] = p;
    }

    /** Relocates a piece with no legality checks; callers must validate first. */
    Piece move(Position from, Position to) {
        Piece captured = at(to);
        set(to, at(from));
        set(from, null);
        return captured;
    }

    public Board clone() {
        Board nb = new Board();
        for (int r = 0; r < 8; r++) {
            System.arraycopy(this.squares[r], 0, nb.squares[r], 0, 8);
        }
        return nb;
    }

    private static int sign(int n) {
        if (n > 0) return 1;
        if (n < 0) return -1;
        return 0;
    }

    /**
     * Reports whether all squares strictly between from and to (on a straight
     * or diagonal line) are empty. Callers already verified the line is
     * straight/diagonal; the destination square itself is not checked here.
     */
    boolean clearPath(Position from, Position to) {
        int stepR = sign(to.row - from.row);
        int stepC = sign(to.col - from.col);
        int r = from.row + stepR;
        int c = from.col + stepC;
        while (r != to.row || c != to.col) {
            if (squares[r][c] != null) {
                return false;
            }
            r += stepR;
            c += stepC;
        }
        return true;
    }

    public Position findKing(Color color) {
        for (int r = 0; r < 8; r++) {
            for (int c = 0; c < 8; c++) {
                Piece p = squares[r][c];
                if (p != null && p.getType() == PieceType.KING && p.getColor() == color) {
                    return new Position(r, c);
                }
            }
        }
        return null;
    }

    /**
     * Reports whether any piece of attacker's color can move onto pos. Only
     * meaningful when pos is occupied (used for check detection, where pos is
     * always the defending king's square).
     */
    public boolean isSquareAttacked(Position pos, Color attacker) {
        for (int r = 0; r < 8; r++) {
            for (int c = 0; c < 8; c++) {
                Piece p = squares[r][c];
                if (p != null && p.getColor() == attacker && p.canMove(this, new Position(r, c), pos)) {
                    return true;
                }
            }
        }
        return false;
    }

    public boolean isInCheck(Color color) {
        Position kingPos = findKing(color);
        if (kingPos == null) {
            return false;
        }
        return isSquareAttacked(kingPos, color.opposite());
    }

    /** Package-private accessor for Game's brute-force checkmate scan. */
    Piece rawAt(int r, int c) {
        return squares[r][c];
    }
}
