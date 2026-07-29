public class WhipDecorator extends CoffeeDecorator {
    public WhipDecorator(Coffee wrapped) {
        super(wrapped);
    }

    @Override
    public double cost() {
        return wrapped.cost() + 0.75;
    }

    @Override
    public String description() {
        return wrapped.description() + " + Whip";
    }
}
