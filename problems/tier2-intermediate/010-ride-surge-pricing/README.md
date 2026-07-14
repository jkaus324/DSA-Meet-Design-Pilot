# Problem 010 — Ride Surge Pricing Engine

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + Observer | **DSA:** Priority Queue
**Companies:** Uber, Ola | **Time:** 60 minutes

---

## Problem Statement

You are building a surge pricing engine for a ride-sharing platform. The engine computes a multiplier applied to base fares based on real-time supply, demand, time of day, and weather conditions. Multiple independent surge factors combine to produce a final multiplier. Downstream systems (driver apps, ops dashboards) must be notified when the multiplier changes significantly.

**Constraints:**
- Surge multiplier is always >= 1.0
- Final multiplier = max of all active factor multipliers
- Surge change threshold for notifications: delta > 0.5x
- Supported ride types: "economy", "premium", "pool"

---

## Base Requirement — Multi-Factor Surge Calculation

Implement a `SurgePricingEngine` that computes a surge multiplier from multiple independent factors. Each factor is a separate strategy. The engine applies all registered strategies and returns the maximum multiplier.

| Factor | Rule |
|---|---|
| DemandSurge | ratio > 3.0 → 2.0x; ratio > 2.0 → 1.5x; ratio > 1.5 → 1.25x; no drivers → 2.5x |
| WeatherSurge | storm → 2.0x; rain → 1.3x; clear → 1.0x |
| TimeSurge | evening peak → 1.4x; morning peak → 1.2x; other → 1.0x |

**Example:**
```
ctx = {drivers=5, requests=12, weather="rain", timeOfDay="evening", baseFare=100.0}
// DemandSurge: ratio=2.4 → 1.5x
// WeatherSurge: rain → 1.3x
// TimeSurge: evening → 1.4x
// Final: max(1.5, 1.3, 1.4) = 1.5x
calculateSurge(ctx)  →  1.5
calculateFare(req, ctx)  →  150.0
```

**Public methods:**
- `void addStrategy(SurgeStrategy* strategy)`
- `double calculateSurge(const PricingContext& ctx)`
- `double calculateFare(const RideRequest& req, const PricingContext& ctx)`

Adding a new surge factor (e.g., "special event") must require zero changes to existing strategy classes.

---

## Extension 1 — Surge Change Notifications

Register observers that receive alerts when the surge multiplier changes by more than 0.5x compared to the previous calculation. Driver observers are filtered by ride type — an economy driver only receives notifications for economy surge changes.

**Example:**
```
registerSurgeObserver(&opsObserver)
registerSurgeObserver(&economyDriverObserver)  // filtered to rideType="economy"

// Previous surge: 1.0x
calculateFare(economyReq, highDemandCtx)  →  200.0, surge=2.0x
// delta = 1.0 > 0.5 threshold → both observers notified
// opsObserver.onSurgeChange(1.0, 2.0, "economy")  called
// economyDriverObserver.onSurgeChange(1.0, 2.0, "economy")  called
```

**Public method:**
- `void registerSurgeObserver(SurgeObserver* obs)`

---

## Running Tests

```bash
./run-tests.sh 010-ride-surge-pricing cpp
```
