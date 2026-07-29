public class SugarDecorator extends CoffeeDecorator {
    public SugarDecorator(Coffee wrapped) {
        super(wrapped);
    }

    @Override
    public double cost() {
        return wrapped.cost() + 0.25;
    }

    @Override
    public String description() {
        return wrapped.description() + " + Sugar";
    }
}
