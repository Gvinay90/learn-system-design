import java.util.Objects;

public class Position {
    public final int row;
    public final int col;

    public Position(int row, int col) {
        this.row = row;
        this.col = col;
    }

    public boolean inBounds() {
        return row >= 0 && row < 8 && col >= 0 && col < 8;
    }

    /** Parses algebraic notation ("e2") into a Position. */
    public static Position parse(String s) {
        if (s.length() != 2) {
            throw new IllegalArgumentException("invalid square " + s);
        }
        int col = s.charAt(0) - 'a';
        int row = 8 - (s.charAt(1) - '0');
        Position pos = new Position(row, col);
        if (!pos.inBounds()) {
            throw new IllegalArgumentException("invalid square " + s);
        }
        return pos;
    }

    @Override
    public boolean equals(Object o) {
        if (!(o instanceof Position)) return false;
        Position p = (Position) o;
        return row == p.row && col == p.col;
    }

    @Override
    public int hashCode() {
        return Objects.hash(row, col);
    }
}
