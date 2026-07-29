import java.util.HashMap;
import java.util.Map;

/**
 * A generic NxN grid. Win detection is O(1) per move: instead of rescanning
 * the whole board, we keep running counts per row, column and the two
 * diagonals for whichever symbol is placed. A move only ever affects one
 * row, one column and (at most) two diagonals, so checking those counts
 * after each move is enough to detect a win.
 */
public class Board {
    public static final char EMPTY = '\0';

    private final int size;
    private final char[][] cells;
    private final Map<Character, int[]> rowCounts = new HashMap<>();
    private final Map<Character, int[]> colCounts = new HashMap<>();
    private final Map<Character, Integer> diagCount = new HashMap<>();
    private final Map<Character, Integer> antiDiagCount = new HashMap<>();
    private int filledCells = 0;

    public Board(int size) {
        this.size = size;
        this.cells = new char[size][size];
    }

    public int getSize() { return size; }

    public char at(int row, int col) { return cells[row][col]; }

    public boolean isFull() { return filledCells == size * size; }

    private boolean inBounds(int row, int col) {
        return row >= 0 && row < size && col >= 0 && col < size;
    }

    /** Places symbol at (row, col) and returns true if that move completes a line. */
    boolean place(int row, int col, char symbol) {
        if (!inBounds(row, col)) {
            throw new IllegalArgumentException("move is out of bounds");
        }
        if (cells[row][col] != EMPTY) {
            throw new IllegalStateException("cell is already occupied");
        }
        cells[row][col] = symbol;
        filledCells++;

        rowCounts.putIfAbsent(symbol, new int[size]);
        colCounts.putIfAbsent(symbol, new int[size]);
        int[] rows = rowCounts.get(symbol);
        int[] cols = colCounts.get(symbol);
        rows[row]++;
        cols[col]++;
        if (row == col) {
            diagCount.merge(symbol, 1, Integer::sum);
        }
        if (row + col == size - 1) {
            antiDiagCount.merge(symbol, 1, Integer::sum);
        }

        return rows[row] == size
                || cols[col] == size
                || diagCount.getOrDefault(symbol, 0) == size
                || antiDiagCount.getOrDefault(symbol, 0) == size;
    }

    @Override
    public String toString() {
        StringBuilder sb = new StringBuilder();
        for (int r = 0; r < size; r++) {
            for (int c = 0; c < size; c++) {
                char ch = cells[r][c];
                sb.append(ch == EMPTY ? '.' : ch).append(' ');
            }
            sb.append('\n');
        }
        return sb.toString();
    }
}
