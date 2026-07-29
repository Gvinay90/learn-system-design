from decorator import Espresso, MilkDecorator, SugarDecorator, WhipDecorator


def test_base_coffee():
    c = Espresso()
    assert c.cost() == 2.0
    assert c.description() == "Espresso"


def test_single_decorator():
    c = MilkDecorator(Espresso())
    assert c.cost() == 2.5
    assert c.description() == "Espresso + Milk"


def test_stacked_decorators_cumulative_cost_and_description():
    c = Espresso()
    c = MilkDecorator(c)
    c = SugarDecorator(c)
    c = WhipDecorator(c)

    assert c.cost() == 2.0 + 0.5 + 0.25 + 0.75
    assert c.description() == "Espresso + Milk + Sugar + Whip"


def test_decorator_order_independence_of_cost():
    a = WhipDecorator(MilkDecorator(Espresso()))
    b = MilkDecorator(WhipDecorator(Espresso()))
    assert a.cost() == b.cost()
    assert a.description() != b.description()
