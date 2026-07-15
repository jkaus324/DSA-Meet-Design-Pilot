# Data class (given — do not modify).

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def reset_service():
    # TODO: implement this
    return None

def set_inventory(productId, qty):
    # TODO: implement this
    return None

def get_inventory(productId):
    # TODO: implement this
    return None

def create_order_simple(productId, quantity, totalAmount):
    # TODO: implement this
    return None

def get_order_state_str(orderId):
    # TODO: implement this
    return None

def confirm_order(orderId):
    # TODO: implement this
    return None

def ship_order(orderId):
    # TODO: implement this
    return None

def deliver_order(orderId):
    # TODO: implement this
    return None

def cancel_order(orderId):
    # TODO: implement this
    return None

def get_history_size(orderId):
    # TODO: implement this
    return None
