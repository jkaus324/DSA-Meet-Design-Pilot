# Data class (given — do not modify).
class PricingContext:
    def __init__(self, baseFare, availableDrivers, activeRideRequests, timeOfDay, weather):
        self.baseFare = baseFare
        self.availableDrivers = availableDrivers
        self.activeRideRequests = activeRideRequests
        self.timeOfDay = timeOfDay
        self.weather = weather

class RideRequest:
    def __init__(self, userId, pickup, dropoff, rideType):
        self.userId = userId
        self.pickup = pickup
        self.dropoff = dropoff
        self.rideType = rideType

from abc import ABC, abstractmethod

class Strategy(ABC):
    @abstractmethod
    def compare(self, a, b):
        """Return True iff `a` ranks strictly before `b`."""

def calculateSurge(ctx):
    # TODO: implement this
    return None

def calculateFare(req, ctx):
    # TODO: implement this
    return None
