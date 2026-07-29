/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out ChessTest` directly.
 */
public class ChessTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testPawnOpeningMoveHappyPath();
        testIllegalMoveRejected();
        testCheckDetection();
        testFoolsMateCheckmate();
        System.out.println("All ChessTest cases passed.");
    }

    private static void testPawnOpeningMoveHappyPath() {
        Game g = new Game();
        Position from = Position.parse("e2");
        Position to = Position.parse("e4");

        g.move(from, to);
        assertTrue(g.getBoard().at(to) != null && g.getBoard().at(to).getType() == PieceType.PAWN,
                "expected pawn at e4");
        assertTrue(g.getBoard().at(from) == null, "expected e2 to be empty after move");
        assertEquals(Color.BLACK, g.getTurn(), "expected turn to switch to Black");
    }

    private static void testIllegalMoveRejected() {
        Game g = new Game();

        try {
            g.move(Position.parse("e2"), Position.parse("e5"));
            throw new AssertionError("expected IllegalMoveException for 3-square pawn move");
        } catch (Game.IllegalMoveException expected) {
            // expected
        }

        try {
            g.move(Position.parse("a1"), Position.parse("a3"));
            throw new AssertionError("expected IllegalMoveException for blocked rook move");
        } catch (Game.IllegalMoveException expected) {
            // expected
        }
    }

    private static void testCheckDetection() {
        Board board = new Board();
        board.set(Position.parse("e1"), PieceFactory.create(PieceType.KING, Color.WHITE));
        board.set(Position.parse("e8"), PieceFactory.create(PieceType.ROOK, Color.BLACK));
        board.set(Position.parse("a1"), PieceFactory.create(PieceType.KING, Color.BLACK));
        Game g = new Game(board, Color.WHITE);

        assertTrue(g.isInCheck(Color.WHITE), "expected White king on e1 to be in check from rook on e8");
        assertTrue(!g.isCheckmate(Color.WHITE), "expected check but not checkmate: king can step aside");
    }

    // Fastest possible checkmate: 1. f3 e5 2. g4 Qh4#
    private static void testFoolsMateCheckmate() {
        Game g = new Game();
        String[][] moves = {
                {"f2", "f3"},
                {"e7", "e5"},
                {"g2", "g4"},
                {"d8", "h4"},
        };
        for (String[] m : moves) {
            g.move(Position.parse(m[0]), Position.parse(m[1]));
        }

        assertTrue(g.isInCheck(Color.WHITE), "expected White king to be in check after Qh4");
        assertTrue(g.isCheckmate(Color.WHITE), "expected fool's mate to be checkmate");
    }

    private static void assertTrue(boolean condition, String label) {
        if (!condition) {
            throw new AssertionError(label);
        }
    }

    private static void assertEquals(Object expected, Object actual, String label) {
        if (!expected.equals(actual)) {
            throw new AssertionError(label + ": expected " + expected + " but got " + actual);
        }
    }
}
