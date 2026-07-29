public class Main {
    public static void main(String[] args) throws InterruptedException {
        TokenBucketLimiter tokenBucket = new TokenBucketLimiter(2, 1);
        System.out.println("token bucket allow(a) = " + tokenBucket.allow("a"));
        System.out.println("token bucket allow(a) = " + tokenBucket.allow("a"));
        System.out.println("token bucket allow(a) [should block] = " + tokenBucket.allow("a"));

        SlidingWindowLimiter slidingWindow = new SlidingWindowLimiter(2, 1000);
        System.out.println("sliding window allow(b) = " + slidingWindow.allow("b"));
        System.out.println("sliding window allow(b) = " + slidingWindow.allow("b"));
        System.out.println("sliding window allow(b) [should block] = " + slidingWindow.allow("b"));

        RateLimiterTest.runAll();
    }
}
