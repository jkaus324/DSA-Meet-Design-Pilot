# Data class (given — do not modify).
class TwitterOp:
    def __init__(self, kind, i1=0, i2=0):
        self.kind = kind
        self.i1 = i1
        self.i2 = i2

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def twitter_simulate(ops):
    # TODO: implement this
    return None
