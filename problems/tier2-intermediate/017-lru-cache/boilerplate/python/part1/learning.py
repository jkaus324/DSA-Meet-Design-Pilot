# Data class (given — do not modify).
class LruOp:
    def __init__(self, kind, i1=0, i2=0, i3=0, i4=0):
        self.kind = kind
        self.i1 = i1
        self.i2 = i2
        self.i3 = i3
        self.i4 = i4

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def lru_simulate(ops):
    # TODO: implement this
    return None
