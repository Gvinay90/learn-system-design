from strategy import (
    ClearancePricing,
    Item,
    PercentageDiscountPricing,
    RegularPricing,
    ShoppingCart,
)


def test_regular_pricing_no_discount():
    cart = ShoppingCart(RegularPricing())
    cart.add_item(Item("book", 20, 2))
    assert cart.checkout() == 40


def test_percentage_discount_pricing():
    cart = ShoppingCart(PercentageDiscountPricing(percent_off=10))
    cart.add_item(Item("shoes", 100, 1))
    assert cart.checkout() == 90


def test_clearance_pricing():
    cart = ShoppingCart(ClearancePricing(percent_off=20, flat_off=5))
    cart.add_item(Item("jacket", 100, 1))
    assert cart.checkout() == 75


def test_clearance_pricing_floors_at_zero():
    cart = ShoppingCart(ClearancePricing(percent_off=50, flat_off=100))
    cart.add_item(Item("sticker", 5, 1))
    assert cart.checkout() == 0


def test_switching_strategy_at_runtime():
    cart = ShoppingCart(RegularPricing())
    cart.add_item(Item("widget", 50, 2))
    assert cart.checkout() == 100

    cart.strategy = PercentageDiscountPricing(percent_off=25)
    assert cart.checkout() == 75
