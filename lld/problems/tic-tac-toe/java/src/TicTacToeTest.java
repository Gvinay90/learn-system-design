import java.util.List;

/**
 * Plain assert-based test runner (no JUnit dependency needed).
 * Run via `java -cp out Main` (invoked from Main) or `java -cp out TicTacToeTest` directly.
 */
public class TicTacToeTest {

    public static void main(String[] args) {
        runAll();
    }

    public static void runAll() {
        testRowWin();
        testColumnWin();
        testDiagonalWin();
        testDrawDetection();
        testOccupiedAndOutOfBoundsRejected();
        testMoveRejectedAfterGameOver();
        testNxNGeneralization();
        System.out.println("All TicTacToeTest cases passed.");
    }

    private static Game newTestGame() {
        return new Game(3, List.of(new Player("Alice", 'X'), new Player("Bob", 'O')));
    }

    private static void testRowWin() {
        Game g = newTestGame();
        int[][] moves = {{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}};
        boolean won = false;
        for (int[] m : moves) {
            won = g.move(m[0], m[1]);
        }
        assertTrue(won, "expected last move to win");
        assertEquals(GameStatus.WON, g.getStatus(), "status after row win");
        assertEquals("Alice", g.getWinner().getName(), "winner of row win");
    }

    private static void testColumnWin() {
        Game g = newTestGame();
        int[][] moves = {{0, 0}, {0, 1}, {1, 0}, {1, 1}, {2, 0}};
        boolean won = false;
        for (int[] m : moves) {
            won = g.move(m[0], m[1]);
        }
        assertTrue(won, "expected column win");
        assertEquals("Alice", g.getWinner().getName(), "winner of column win");
    }

    private static void testDiagonalWin() {
        Game g = newTestGame();
        int[][] moves = {{0, 0}, {0, 1}, {1, 1}, {0, 2}, {2, 2}};
        boolean won = false;
        for (int[] m : moves) {
            won = g.move(m[0], m[1]);
        }
        assertTrue(won, "expected diagonal win");
        assertEquals("Alice", g.getWinner().getName(), "winner of diagonal win");
    }

    private static void testDrawDetection() {
        Game g = newTestGame();
        int[][] moves = {
                {0, 0}, {0, 1},
                {0, 2}, {1, 1},
                {1, 0}, {1, 2},
                {2, 1}, {2, 0},
                {2, 2}
        };
        boolean lastWon = false;
        for (int[] m : moves) {
            lastWon = g.move(m[0], m[1]);
        }
        assertTrue(!lastWon, "expected no winner in a drawn game");
        assertEquals(GameStatus.DRAW, g.getStatus(), "status after draw");
    }

    private static void testOccupiedAndOutOfBoundsRejected() {
        Game g = newTestGame();
        g.move(0, 0);
        try {
            g.move(0, 0);
            throw new AssertionError("expected IllegalStateException for occupied cell");
        } catch (IllegalStateException e) {
            // expected
        }
        try {
            g.move(5, 5);
            throw new AssertionError("expected IllegalArgumentException for out-of-bounds move");
        } catch (IllegalArgumentException e) {
            // expected
        }
    }

    private static void testMoveRejectedAfterGameOver() {
        Game g = newTestGame();
        int[][] moves = {{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}};
        for (int[] m : moves) {
            g.move(m[0], m[1]);
        }
        assertEquals(GameStatus.WON, g.getStatus(), "expected game won before extra move");
        try {
            g.move(2, 2);
            throw new AssertionError("expected GameOverException");
        } catch (Game.GameOverException e) {
            // expected
        }
    }

    private static void testNxNGeneralization() {
        Game g = new Game(5, List.of(new Player("Alice", 'X'), new Player("Bob", 'O')));
        int[][] moves = {
                {0, 4}, {0, 0},
                {1, 3}, {0, 1},
                {2, 2}, {0, 2},
                {3, 1}, {0, 3},
                {4, 0}
        };
        boolean won = false;
        Player winner = null;
        for (int[] m : moves) {
            if (g.move(m[0], m[1])) {
                won = true;
                winner = g.getWinner();
            }
        }
        assertTrue(won, "expected a win on the 5x5 anti-diagonal");
        assertEquals("Alice", winner.getName(), "winner of 5x5 anti-diagonal");
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
