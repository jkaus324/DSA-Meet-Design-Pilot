# Data class (given — do not modify).

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def reset_service():
    # TODO: implement this
    return None

def notify_event(event, userIds, subscribedChannels):
    # TODO: implement this
    return None

def notify_priority(event, priority, userIds, subscribedChannels, minPriority):
    # TODO: implement this
    return None

def notify_priority_level(priority):
    # TODO: implement this
    return None
