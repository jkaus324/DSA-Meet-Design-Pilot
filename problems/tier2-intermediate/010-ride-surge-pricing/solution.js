// Ride surge pricing — Strategy + Observer (JavaScript).

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

class DemandSurge {
  multiplier(ctx) {
    if (ctx.availableDrivers === 0) return 2.5;
    const ratio = ctx.activeRideRequests / ctx.availableDrivers;
    if (ratio > 3.0) return 2.0;
    if (ratio > 2.0) return 1.5;
    if (ratio > 1.5) return 1.25;
    return 1.0;
  }
}

class WeatherSurge {
  multiplier(ctx) {
    if (ctx.weather === 'storm') return 1.8;
    if (ctx.weather === 'rain') return 1.3;
    return 1.0;
  }
}

class TimeSurge {
  multiplier(ctx) {
    if (ctx.timeOfDay === 'evening') return 1.5;
    if (ctx.timeOfDay === 'morning') return 1.2;
    return 1.0;
  }
}

class PricingEngine {
  constructor() {
    this.strategies = [];
    this.observers = [];
    this.lastSurge = 1.0;
  }

  addStrategy(s) { this.strategies.push(s); }

  addObserver(o) { this.observers.push(o); }

  clearObservers() { this.observers = []; }

  calculateSurge(ctx, rideType = 'all') {
    let mult = 1.0;
    for (const s of this.strategies) {
      mult = Math.max(mult, s.multiplier(ctx));
    }
    mult = Math.min(mult, 3.0);
    if (Math.abs(mult - this.lastSurge) > PricingEngine.CHANGE_THRESHOLD) {
      for (const o of this.observers) {
        o.onSurgeChange(this.lastSurge, mult, rideType);
      }
    }
    this.lastSurge = mult;
    return mult;
  }

  calculateFare(ctx, rideType = 'all') {
    return ctx.baseFare * this.calculateSurge(ctx, rideType);
  }
}
PricingEngine.CHANGE_THRESHOLD = 0.5;

let _engine = new PricingEngine();
_engine.addStrategy(new DemandSurge());
_engine.addStrategy(new WeatherSurge());
_engine.addStrategy(new TimeSurge());

function calculateSurge(ctx) {
  return _engine.calculateSurge(ctx);
}

function calculateFare(req, ctx) {
  return _engine.calculateFare(ctx, req.rideType);
}

function registerSurgeObserver(obs) {
  _engine.clearObservers();
  _engine.addObserver(obs);
}

module.exports = {
  PricingContext, RideRequest,
  DemandSurge, WeatherSurge, TimeSurge, PricingEngine,
  calculateSurge, calculateFare, registerSurgeObserver,
};
