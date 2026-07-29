public class MilkDecorator extends CoffeeDecorator {
    public MilkDecorator(Coffee wrapped) {
        super(wrapped);
    }

    @Override
    public double cost() {
        return wrapped.cost() + 0.5;
    }

    @Override
    public String description() {
        return wrapped.description() + " + Milk";
    }
}
