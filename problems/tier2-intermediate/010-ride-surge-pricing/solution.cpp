#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <cmath>
using namespace std;

// ─── Data Models ─────────────────────────────────────────────────────────────

struct RideRequest {
    string userId;
    string pickup;
    string dropoff;
    string rideType;
};

struct PricingContext {
    double baseFare;
    int    availableDrivers;
    int    activeRideRequests;
    string timeOfDay;   // "morning", "evening", "night"
    string weather;     // "clear", "rain", "storm"
};

// ─── Strategy Interface ──────────────────────────────────────────────────────

class SurgeStrategy {
public:
    virtual double multiplier(const PricingContext& ctx) = 0;
    virtual ~SurgeStrategy() = default;
};

// ─── Concrete Strategies ─────────────────────────────────────────────────────

class DemandSurge : public SurgeStrategy {
public:
    double multiplier(const PricingContext& ctx) override {
        if (ctx.availableDrivers == 0) return 2.5;
        double ratio = (double)ctx.activeRideRequests / ctx.availableDrivers;
        if (ratio > 3.0) return 2.0;
        if (ratio > 2.0) return 1.5;
        if (ratio > 1.5) return 1.25;
        return 1.0;
    }
};

class WeatherSurge : public SurgeStrategy {
public:
    double multiplier(const PricingContext& ctx) override {
        if (ctx.weather == "storm") return 1.8;
        if (ctx.weather == "rain")  return 1.3;
        return 1.0;
    }
};

class TimeSurge : public SurgeStrategy {
public:
    double multiplier(const PricingContext& ctx) override {
        if (ctx.timeOfDay == "evening") return 1.5;
        if (ctx.timeOfDay == "morning") return 1.2;
        return 1.0;
    }
};

// ─── Observer Interface (3-arg for Part 2) ───────────────────────────────────

class SurgeObserver {
public:
    virtual void onSurgeChange(double oldMultiplier, double newMultiplier,
                               const string& rideType) = 0;
    virtual ~SurgeObserver() = default;
};

// ─── Pricing Engine ──────────────────────────────────────────────────────────

class PricingEngine {
    vector<SurgeStrategy*> strategies;
    vector<SurgeObserver*> observers;
    double lastSurge = 1.0;
    static constexpr double CHANGE_THRESHOLD = 0.5;

public:
    void addStrategy(SurgeStrategy* s) { strategies.push_back(s); }
    void addObserver(SurgeObserver* o)  { observers.push_back(o); }
    void clearObservers()               { observers.clear(); }

    double calculateSurge(const PricingContext& ctx,
                          const string& rideType = "all") {
        // Combine strategies by taking the max multiplier
        double mult = 1.0;
        for (auto* s : strategies) {
            mult = max(mult, s->multiplier(ctx));
        }
        // Cap at 3.0x
        mult = min(mult, 3.0);

        // Notify observers if change exceeds threshold
        if (fabs(mult - lastSurge) > CHANGE_THRESHOLD) {
            for (auto* o : observers) {
                o->onSurgeChange(lastSurge, mult, rideType);
            }
        }
        lastSurge = mult;
        return mult;
    }

    double calculateFare(const PricingContext& ctx,
                         const string& rideType = "all") {
        return ctx.baseFare * calculateSurge(ctx, rideType);
    }
};

// ─── Global Engine & Free Functions ──────────────────────────────────────────

static DemandSurge  g_demandSurge;
static WeatherSurge g_weatherSurge;
static TimeSurge    g_timeSurge;

static PricingEngine& getGlobalEngine() {
    static PricingEngine engine;
    static bool initialized = false;
    if (!initialized) {
        engine.addStrategy(&g_demandSurge);
        engine.addStrategy(&g_weatherSurge);
        engine.addStrategy(&g_timeSurge);
        initialized = true;
    }
    return engine;
}

double calculateSurge(const PricingContext& ctx) {
    return getGlobalEngine().calculateSurge(ctx);
}

double calculateFare(const RideRequest& req, const PricingContext& ctx) {
    return getGlobalEngine().calculateFare(ctx, req.rideType);
}

void registerSurgeObserver(SurgeObserver* obs) {
    getGlobalEngine().clearObservers();
    getGlobalEngine().addObserver(obs);
}

// ─── Main ────────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    PricingContext ctx = {100.0, 5, 15, "evening", "storm"};
    RideRequest req = {"user1", "Downtown", "Airport", "economy"};

    double surge = calculateSurge(ctx);
    double fare  = calculateFare(req, ctx);

    cout << "Surge multiplier: " << surge << "x" << endl;
    cout << "Fare: $" << fare << endl;

    return 0;
}
#endif
