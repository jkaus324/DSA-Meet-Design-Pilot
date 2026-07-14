# Data class (given — do not modify).
class ElevOp:
    def __init__(self, kind, s1="", i1=0, i2=0):
        self.kind = kind
        self.s1 = s1
        self.i1 = i1
        self.i2 = i2

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def elevator_simulate(ops):
    # TODO: implement this
    return None
