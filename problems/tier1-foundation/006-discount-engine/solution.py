"""Discount Engine — Strategy + Decorator reference solution (Python)."""

from abc import ABC, abstractmethod


class CartItem:
    def __init__(self, name, price, quantity, category=""):
        self.name = name
        self.price = price
        self.quantity = quantity
        self.category = category


class Discount(ABC):
    @abstractmethod
    def apply(self, cart):
        ...


class PercentageDiscount(Discount):
    def __init__(self, pct):
        self.pct = pct

    def apply(self, cart):
        total = sum(i.price * i.quantity for i in cart)
        return total * (1.0 - self.pct / 100.0)


class FlatDiscount(Discount):
    def __init__(self, amount):
        self.amount = amount

    def apply(self, cart):
        total = sum(i.price * i.quantity for i in cart)
        return max(0.0, total - self.amount)


class BuyXGetYDiscount(Discount):
    def __init__(self, buy, free):
        self.buy = buy
        self.free = free

    def apply(self, cart):
        group = self.buy + self.free
        total = 0.0
        for it in cart:
            groups = it.quantity // group
            remainder = it.quantity % group
            paid = groups * self.buy + min(remainder, self.buy)
            total += paid * it.price
        return total


class StackedDiscount(Discount):
    def __init__(self, discounts):
        self.discounts = discounts

    def apply(self, cart):
        current = sum(i.price * i.quantity for i in cart)
        for d in self.discounts:
            temp = [CartItem("subtotal", current, 1, "")]
            current = d.apply(temp)
        return current


def apply_percentage_discount(cart, percentage):
    return PercentageDiscount(percentage).apply(cart)


def apply_flat_discount(cart, amount):
    return FlatDiscount(amount).apply(cart)


def apply_bogo(cart, buy, free):
    return BuyXGetYDiscount(buy, free).apply(cart)


def apply_percentage_with_eligibility(cart, percentage, minCartValue,
                                      requireFirstTimeUser, isFirstTimeUser,
                                      eligibleCategory):
    raw = sum(i.price * i.quantity for i in cart)
    if raw < minCartValue:
        return raw
    if requireFirstTimeUser and not isFirstTimeUser:
        return raw
    if eligibleCategory:
        eligible = [i for i in cart if i.category == eligibleCategory]
        non_eligible = sum(i.price * i.quantity for i in cart if i.category != eligibleCategory)
        return PercentageDiscount(percentage).apply(eligible) + non_eligible
    return PercentageDiscount(percentage).apply(cart)


def reset_service():
    pass
