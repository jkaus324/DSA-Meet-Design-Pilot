# Data class (given).
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

# HINT: introduce an abstraction so new ranking rules don't change existing code.

# HINT: pick the field that defines 'better' for this ranking and compare the two.
def calculateSurge(ctx):
    # TODO: write your solution
    return None

# HINT: pick the field that defines 'better' for this ranking and compare the two.
def calculateFare(req, ctx):
    # TODO: write your solution
    return None
