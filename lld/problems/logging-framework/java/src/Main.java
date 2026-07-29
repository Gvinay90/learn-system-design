public class Main {
    public static void main(String[] args) {
        Logger logger = new Logger(Level.INFO, new ConsoleAppender());
        logger.debug("this is filtered out below INFO threshold");
        logger.info("service started");
        logger.warn("cache miss rate elevated");
        logger.error("failed to connect to downstream");

        LoggingFrameworkTest.runAll();
    }
}
