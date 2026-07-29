public class Main {
    public static void main(String[] args) {
        Coffee coffee = new Espresso();
        coffee = new MilkDecorator(coffee);
        coffee = new SugarDecorator(coffee);
        coffee = new WhipDecorator(coffee);

        System.out.println(coffee.description() + " costs " + coffee.cost());

        DecoratorTest.runAll();
    }
}
