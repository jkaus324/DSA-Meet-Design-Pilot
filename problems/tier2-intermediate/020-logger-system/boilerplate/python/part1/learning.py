# Data class (given — do not modify).
class LogOp:
    def __init__(self, kind, s1="", s2="", i1=0):
        self.kind = kind
        self.s1 = s1
        self.s2 = s2
        self.i1 = i1

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def logger_simulate(ops):
    # TODO: implement this
    return None
