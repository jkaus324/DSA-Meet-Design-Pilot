// Data class (given — do not modify).
class PricingContext {
  constructor(baseFare, availableDrivers, activeRideRequests, timeOfDay, weather) {
    this.baseFare = baseFare;
    this.availableDrivers = availableDrivers;
    this.activeRideRequests = activeRideRequests;
    this.timeOfDay = timeOfDay;
    this.weather = weather;
  }
}

class RideRequest {
  constructor(userId, pickup, dropoff, rideType) {
    this.userId = userId;
    this.pickup = pickup;
    this.dropoff = dropoff;
    this.rideType = rideType;
  }
}

// Strategy — base strategy. Subclasses implement compare().
class Strategy {
  // Return true iff `a` ranks strictly before `b`.
  compare(a, b) { throw new Error('not implemented'); }
}

function calculateSurge(ctx) {
  // TODO: implement this
  return null;
}

function calculateFare(req, ctx) {
  // TODO: implement this
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
module.exports = { PricingContext, RideRequest, calculateSurge, calculateFare };
