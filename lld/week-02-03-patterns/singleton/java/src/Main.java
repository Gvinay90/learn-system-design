public class Main {
    public static void main(String[] args) {
        AppConfig config = AppConfig.getInstance();
        config.set("env", "production");
        System.out.println("Config instance id: " + config.getId());
        System.out.println("env=" + AppConfig.getInstance().get("env"));

        SingletonTest.runAll();
    }
}
