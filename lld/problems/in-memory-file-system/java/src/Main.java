import java.util.List;

public class Main {
    public static void main(String[] args) {
        FileSystem fs = new FileSystem();
        fs.mkdir("/home");
        fs.mkdir("/home/docs");
        fs.createFile("/home/docs/readme.md", "hello world");

        System.out.println("pwd: " + fs.pwd());
        fs.cd("/home/docs");
        System.out.println("pwd after cd: " + fs.pwd());
        System.out.println("readme.md: " + fs.readFile("readme.md"));

        fs.writeFile("readme.md", "\nmore content", true);
        System.out.println("readme.md after append: " + fs.readFile("readme.md"));

        List<Entry> entries = fs.ls("/home");
        System.out.print("ls /home: ");
        for (Entry e : entries) {
            System.out.print(e.getName() + (e.isDir() ? "/ " : " "));
        }
        System.out.println();

        FileSystemTest.runAll();
    }
}
