public class Main {
    public static void main(String[] args) {
        Game game = new Game();

        game.move(Position.parse("e2"), Position.parse("e4"));
        System.out.println("White plays e2-e4; turn is now " + game.getTurn());

        game.move(Position.parse("e7"), Position.parse("e5"));
        System.out.println("Black plays e7-e5; turn is now " + game.getTurn());

        // Fool's mate: fastest possible checkmate.
        Game mate = new Game();
        mate.move(Position.parse("f2"), Position.parse("f3"));
        mate.move(Position.parse("e7"), Position.parse("e5"));
        mate.move(Position.parse("g2"), Position.parse("g4"));
        mate.move(Position.parse("d8"), Position.parse("h4"));
        System.out.println("After Fool's Mate: White in check = " + mate.isInCheck(Color.WHITE)
                + ", checkmate = " + mate.isCheckmate(Color.WHITE));

        ChessTest.runAll();
    }
}
