"""Ride surge pricing — Strategy + Observer."""


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


class DemandSurge:
    def multiplier(self, ctx):
        if ctx.availableDrivers == 0:
            return 2.5
        ratio = ctx.activeRideRequests / ctx.availableDrivers
        if ratio > 3.0:
            return 2.0
        if ratio > 2.0:
            return 1.5
        if ratio > 1.5:
            return 1.25
        return 1.0


class WeatherSurge:
    def multiplier(self, ctx):
        if ctx.weather == "storm":
            return 1.8
        if ctx.weather == "rain":
            return 1.3
        return 1.0


class TimeSurge:
    def multiplier(self, ctx):
        if ctx.timeOfDay == "evening":
            return 1.5
        if ctx.timeOfDay == "morning":
            return 1.2
        return 1.0


class PricingEngine:
    CHANGE_THRESHOLD = 0.5

    def __init__(self):
        self.strategies = []
        self.observers = []
        self.lastSurge = 1.0

    def addStrategy(self, s):
        self.strategies.append(s)

    def addObserver(self, o):
        self.observers.append(o)

    def clearObservers(self):
        self.observers = []

    def calculateSurge(self, ctx, rideType="all"):
        mult = 1.0
        for s in self.strategies:
            mult = max(mult, s.multiplier(ctx))
        mult = min(mult, 3.0)
        if abs(mult - self.lastSurge) > self.CHANGE_THRESHOLD:
            for o in self.observers:
                o.onSurgeChange(self.lastSurge, mult, rideType)
        self.lastSurge = mult
        return mult

    def calculateFare(self, ctx, rideType="all"):
        return ctx.baseFare * self.calculateSurge(ctx, rideType)


_engine = PricingEngine()
_engine.addStrategy(DemandSurge())
_engine.addStrategy(WeatherSurge())
_engine.addStrategy(TimeSurge())


def calculateSurge(ctx):
    return _engine.calculateSurge(ctx)


def calculateFare(req, ctx):
    return _engine.calculateFare(ctx, req.rideType)


def registerSurgeObserver(obs):
    _engine.clearObservers()
    _engine.addObserver(obs)
