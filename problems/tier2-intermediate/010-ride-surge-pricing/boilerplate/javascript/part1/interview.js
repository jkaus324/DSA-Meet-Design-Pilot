// Data class (given).
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

// TODO: design and implement your solution.
// Required functions:
//   function calculateSurge(ctx)
//   function calculateFare(req, ctx)

function calculateSurge(ctx) {
  // TODO: write your solution
  return null;
}

function calculateFare(req, ctx) {
  // TODO: write your solution
  return null;
}

// ── Export everything the test runner needs (do not remove) ──
// If you add classes (e.g. strategy subclasses), add them here too.
module.exports = { PricingContext, RideRequest, calculateSurge, calculateFare };
