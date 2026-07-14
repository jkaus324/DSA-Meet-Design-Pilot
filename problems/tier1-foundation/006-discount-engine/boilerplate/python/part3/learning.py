# Data class (given — do not modify).
class CartItem:
    def __init__(self, name, price, quantity, category=""):
        self.name = name
        self.price = price
        self.quantity = quantity
        self.category = category

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def apply_percentage_discount(cart, percentage):
    # TODO: implement this
    return None

def apply_flat_discount(cart, amount):
    # TODO: implement this
    return None

def apply_bogo(cart, buyCount, freeCount):
    # TODO: implement this
    return None

def apply_percentage_with_eligibility(cart, percentage, minCartValue, requireFirstTimeUser, isFirstTimeUser, eligibleCategory):
    # TODO: implement this
    return None
