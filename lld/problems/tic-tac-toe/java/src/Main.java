import java.util.List;

public class Main {
    public static void main(String[] args) {
        Player alice = new Player("Alice", 'X');
        Player bob = new Player("Bob", 'O');
        Game game = new Game(3, List.of(alice, bob));

        int[][] moves = {{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}};
        for (int[] m : moves) {
            boolean won = game.move(m[0], m[1]);
            if (won) {
                System.out.println(game.getWinner().getName() + " wins!");
            }
        }
        System.out.println(game.getBoard());

        TicTacToeTest.runAll();
    }
}
