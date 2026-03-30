#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

struct RideRequest { string userId, pickup, dropoff, rideType; };
struct Driver { string id; double rating; string rideType; bool available; };
struct PricingContext { double baseFare; int availableDrivers, activeRideRequests; string timeOfDay, weather; };

// ─── SurgeStrategy Interface ──────────────────────────────────────────────────

class SurgeStrategy {
public:
    virtual double multiplier(const PricingContext& ctx) = 0;
    virtual ~SurgeStrategy() = default;
};

// Copy your Part 1 strategies here (DemandSurge, WeatherSurge, TimeSurge, etc.)

// ─── NEW: SurgeObserver Interface ────────────────────────────────────────────
// HINT: Observers receive the old multiplier, new multiplier, and ride type.

class SurgeObserver {
public:
    virtual void onSurgeChange(double oldMult, double newMult,
                               const string& rideType) = 0;
    virtual ~SurgeObserver() = default;
};

// TODO: Implement DriverObserver (filters by driver's rideType)
// TODO: Implement OpsDashboardObserver (receives all surge changes)

// ─── SurgePricingEngine ───────────────────────────────────────────────────────

class SurgePricingEngine {
    vector<SurgeStrategy*> strategies;
    vector<SurgeObserver*> observers;
    double lastMultiplier = 1.0;
    const double CHANGE_THRESHOLD = 0.5;
public:
    void addStrategy(SurgeStrategy* s) { strategies.push_back(s); }
    void addObserver(SurgeObserver* o) { observers.push_back(o); }

    double calculateSurge(const PricingContext& ctx) {
        double mult = 1.0;
        for (auto* s : strategies) mult = max(mult, s->multiplier(ctx));
        // TODO: If |mult - lastMultiplier| > CHANGE_THRESHOLD, notify observers
        lastMultiplier = mult;
        return mult;
    }
};

double calculateSurge(const PricingContext& ctx) {
    // TODO: implement using SurgePricingEngine
    return 1.0;
}

double calculateFare(const RideRequest& req, const PricingContext& ctx) {
    return ctx.baseFare * calculateSurge(ctx);
}

void registerSurgeObserver(SurgeObserver* obs) {
    // TODO: register with the global engine
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Part 2: Surge notifications — implement TODOs above." << endl;
    return 0;
}
#endif
