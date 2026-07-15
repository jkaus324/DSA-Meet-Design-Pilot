# Data class (given — do not modify).

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def reset_service():
    # TODO: implement this
    return None

def reset():
    # TODO: implement this
    return None

def getState():
    # TODO: implement this
    return None

def selectItem(item):
    # TODO: implement this
    return None

def insertMoney(amount):
    # TODO: implement this
    return None

def dispense():
    # TODO: implement this
    return None

def cancel():
    # TODO: implement this
    return None

def enterMaintenance(pin):
    # TODO: implement this
    return None

def exitMaintenance(pin):
    # TODO: implement this
    return None

def restock(item, qty):
    # TODO: implement this
    return None

def vm_get_quantity(item):
    # TODO: implement this
    return None
