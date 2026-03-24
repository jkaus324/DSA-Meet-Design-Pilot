# Problem 010 — Ride Surge Pricing Engine

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + Observer | **DSA:** Priority Queue
**Companies:** Uber, Ola | **Time:** 60 minutes

---

## Problem Statement

Build a surge pricing engine for a ride-sharing platform. The system must:

1. Track supply (available drivers) and demand (active ride requests) in real-time
2. Calculate a **surge multiplier** when demand exceeds supply
3. Account for multiple surge factors: demand/supply ratio, time of day, and weather conditions
4. Notify downstream systems when surge changes significantly

---

## Before You Code

**Two patterns working together:**
- **Strategy**: The surge *calculation algorithm* is swappable — demand-based, weather-based, time-based factors can each be a separate strategy
- **Observer**: When the surge multiplier changes significantly, notify multiple downstream systems (drivers, ops dashboard, pricing display)

**The DSA angle:** Surge factors combine independently — the final multiplier is the max (or sum) of all active factors. Think about how to model this without a giant if-else chain.

---

## Data Structures

```cpp
struct RideRequest {
    string userId;
    string pickup;
    string dropoff;
    string rideType;  // "economy", "premium", "pool"
};

struct Driver {
    string id;
    double rating;
    string rideType;
    bool   available;
};

struct PricingContext {
    double baseFare;
    int    availableDrivers;
    int    activeRideRequests;
    string timeOfDay;   // "morning", "evening", "night"
    string weather;     // "clear", "rain", "storm"
};
```

---

## Part 1

**Base requirement — Surge calculation**

Implement a surge pricing engine that computes a surge multiplier from multiple independent factors:

| Factor | Rule |
|--------|------|
| **Demand surge** | demand/supply ratio > 1.5 → 1.25x; > 2.0 → 1.5x; > 3.0 → 2.0x; no drivers → 2.5x |
| **Weather surge** | storm → 2.0x; rain → 1.3x; clear → 1.0x |
| **Time surge** | evening peak → 1.4x; morning peak → 1.2x; other → 1.0x |

**Rule:** The final surge multiplier is the **maximum** of all active factors. Surge is always at least 1.0x.

**Entry points (tests will call these):**
```cpp
double calculateSurge(const PricingContext& ctx);
double calculateFare(const RideRequest& req, const PricingContext& ctx);
// fare = ctx.baseFare * calculateSurge(ctx)
```

**What to implement:**
```cpp
class SurgeStrategy {
public:
    virtual double multiplier(const PricingContext& ctx) = 0;
    virtual ~SurgeStrategy() = default;
};

class DemandSurge  : public SurgeStrategy { ... };
class WeatherSurge : public SurgeStrategy { ... };
class TimeSurge    : public SurgeStrategy { ... };

class SurgePricingEngine {
    vector<SurgeStrategy*> strategies;
public:
    void addStrategy(SurgeStrategy* s);
    double calculateSurge(const PricingContext& ctx);
};
```

Adding a new surge factor (e.g., "special event") must require **zero changes** to existing strategies.

---

## Part 2

**Extension 1 — Surge change notifications**

The ops team needs **real-time alerts** when surge multiplier changes significantly (by more than 0.5x):

- **Drivers** should be notified when surge increases (encourages them to go online)
- **Ops dashboard** receives all surge change events
- Driver notifications are filtered by `rideType` — economy drivers only see economy surge changes

**New entry points:**
```cpp
void registerSurgeObserver(SurgeObserver* obs);
// calculateFare() now also triggers notifications when surge changes significantly
```

**Observer interface:**
```cpp
class SurgeObserver {
public:
    virtual void onSurgeChange(double oldMult, double newMult,
                               const string& rideType) = 0;
};
```

**Design challenge:** How do you combine Strategy (surge calculation) with Observer (surge notifications) cleanly? Does the engine become the subject? Or is there a separate notifier?

---

## Running Tests

```bash
./run-tests.sh 010-ride-surge-pricing cpp
```
