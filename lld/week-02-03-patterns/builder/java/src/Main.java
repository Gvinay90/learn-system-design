public class Main {
    public static void main(String[] args) {
        HttpRequest req = new HttpRequestBuilder()
                .method("POST")
                .url("https://api.example.com/orders")
                .header("Content-Type", "application/json")
                .header("Authorization", "Bearer token123")
                .body("{\"item\":\"widget\"}")
                .build();

        System.out.println(req);

        BuilderTest.runAll();
    }
}
