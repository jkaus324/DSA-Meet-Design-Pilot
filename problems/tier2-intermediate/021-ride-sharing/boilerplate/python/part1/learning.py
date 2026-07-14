# Data class (given — do not modify).
class RideOp:
    def __init__(self, kind, s1="", s2="", s3="", s4="", i1=0, i2=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.s3 = s3
        self.s4 = s4
        self.i1 = i1
        self.i2 = i2

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def ride_simulate(ops):
    # TODO: implement this
    return None
