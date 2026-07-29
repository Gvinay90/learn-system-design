/** Factory pattern: constructs the concrete Piece for a given type/color. */
public final class PieceFactory {
    private PieceFactory() {
    }

    public static Piece create(PieceType type, Color color) {
        switch (type) {
            case PAWN:
                return new PawnPiece(color);
            case KNIGHT:
                return new KnightPiece(color);
            case BISHOP:
                return new BishopPiece(color);
            case ROOK:
                return new RookPiece(color);
            case QUEEN:
                return new QueenPiece(color);
            case KING:
                return new KingPiece(color);
            default:
                throw new IllegalArgumentException("unknown piece type " + type);
        }
    }
}
