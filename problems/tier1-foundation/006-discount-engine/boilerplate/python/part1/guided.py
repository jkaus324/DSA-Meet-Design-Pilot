# Data class (given).
class CartItem:
    def __init__(self, name, price, quantity, category=""):
        self.name = name
        self.price = price
        self.quantity = quantity
        self.category = category

# HINT: introduce an abstraction so new ranking rules don't change existing code.

# HINT: pick the field that defines 'better' for this ranking and compare the two.
def apply_percentage_discount(cart, percentage):
    # TODO: write your solution
    return None

# HINT: pick the field that defines 'better' for this ranking and compare the two.
def apply_flat_discount(cart, amount):
    # TODO: write your solution
    return None

# HINT: pick the field that defines 'better' for this ranking and compare the two.
def apply_bogo(cart, buyCount, freeCount):
    # TODO: write your solution
    return None
