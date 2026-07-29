public abstract class CoffeeDecorator implements Coffee {
    protected final Coffee wrapped;

    protected CoffeeDecorator(Coffee wrapped) {
        this.wrapped = wrapped;
    }
}
