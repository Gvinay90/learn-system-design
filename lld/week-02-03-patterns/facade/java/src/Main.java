public class Main {
    public static void main(String[] args) {
        HomeTheaterFacade theater = new HomeTheaterFacade();
        theater.watchMovie("Inception");
        System.out.println("Now playing: " + theater.getDvd().getMovie());
        theater.endMovie();
        System.out.println("Movie ended, lights dimmed: " + theater.getLights().isDimmed());

        HomeTheaterFacadeTest.runAll();
    }
}
