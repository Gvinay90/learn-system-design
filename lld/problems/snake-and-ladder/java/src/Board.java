import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class Board {
    private final int size;
    private final Map<Integer, Integer> entities = new HashMap<>();

    public Board(int size, List<Entity> snakes, List<Entity> ladders) {
        this.size = size;
        for (Entity s : snakes) {
            if (s.getStart() <= s.getEnd()) {
                throw new IllegalArgumentException(
                        "snake start " + s.getStart() + " must be greater than end " + s.getEnd());
            }
            addEntity(s.getStart(), s.getEnd());
        }
        for (Entity l : ladders) {
            if (l.getStart() >= l.getEnd()) {
                throw new IllegalArgumentException(
                        "ladder start " + l.getStart() + " must be less than end " + l.getEnd());
            }
            addEntity(l.getStart(), l.getEnd());
        }
    }

    private void addEntity(int from, int to) {
        if (from <= 1 || from >= size) {
            throw new IllegalArgumentException("entity start " + from + " must be within (1, " + size + ")");
        }
        if (entities.containsKey(from)) {
            throw new IllegalArgumentException("cell " + from + " already has a snake or ladder");
        }
        entities.put(from, to);
    }

    public int getSize() { return size; }

    /** Applies any snake/ladder at the landed-on cell, returning the final resting cell. */
    public int resolve(int cell) {
        return entities.getOrDefault(cell, cell);
    }
}
