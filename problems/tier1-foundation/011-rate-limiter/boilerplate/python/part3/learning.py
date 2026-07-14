# Data class (given — do not modify).

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def reset_service():
    # TODO: implement this
    return None

def init_limiter(maxRequests, windowSize):
    # TODO: implement this
    return None

def allow_request_simple(clientId, timestamp, endpoint):
    # TODO: implement this
    return None

def get_request_count(clientId):
    # TODO: implement this
    return None

def allow_request_with_strategy_simple(algorithm, clientId, timestamp, endpoint):
    # TODO: implement this
    return None

def allow_request_for_tier_str(tier, clientId, timestamp, endpoint):
    # TODO: implement this
    return None
