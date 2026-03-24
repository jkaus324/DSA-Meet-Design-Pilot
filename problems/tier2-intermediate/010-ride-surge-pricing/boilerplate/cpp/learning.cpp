#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct PricingContext {
    double baseFare;
    int    availableDrivers;
    int    activeRideRequests;
    string timeOfDay;   // "morning", "evening", "night"
    string weather;     // "clear", "rain", "storm"
};

// ─── Surge Strategy Interface ─────────────────────────────────────────────────

class SurgeStrategy {
public:
    virtual double multiplier(const PricingContext& ctx) = 0;
    virtual ~SurgeStrategy() = default;
};

// ─── Concrete Surge Strategies ────────────────────────────────────────────────

class DemandSurge : public SurgeStrategy {
public:
    double multiplier(const PricingContext& ctx) override {
        // TODO: Calculate demand ratio (requests / drivers)
        //       Return 1.0 if balanced, up to 2.0x if demand >> supply
        return 1.0;
    }
};

class WeatherSurge : public SurgeStrategy {
public:
    double multiplier(const PricingContext& ctx) override {
        // TODO: Return multiplier based on weather
        //       clear → 1.0, rain → 1.3, storm → 1.8
        return 1.0;
    }
};

class TimeSurge : public SurgeStrategy {
public:
    double multiplier(const PricingContext& ctx) override {
        // TODO: Return multiplier based on time of day
        //       morning → 1.2 (peak), evening → 1.5 (peak), night → 1.0
        return 1.0;
    }
};

// ─── Surge Observer Interface ─────────────────────────────────────────────────

class SurgeObserver {
public:
    virtual void onSurgeChange(double oldMultiplier, double newMultiplier) = 0;
    virtual ~SurgeObserver() = default;
};

class DriverNotifier : public SurgeObserver {
public:
    void onSurgeChange(double old_, double new_) override {
        // TODO: Print driver-facing surge alert if new surge > 1.5x
    }
};

class RiderNotifier : public SurgeObserver {
public:
    void onSurgeChange(double old_, double new_) override {
        // TODO: Print rider-facing surge warning if new surge > 1.5x
    }
};

// ─── Pricing Engine ──────────────────────────────────────────────────────────

class PricingEngine {
private:
    vector<SurgeStrategy*> strategies;
    vector<SurgeObserver*> observers;
    double lastSurge = 1.0;
public:
    void addStrategy(SurgeStrategy* s) { strategies.push_back(s); }
    void addObserver(SurgeObserver* o) { observers.push_back(o); }

    double calculateSurge(const PricingContext& ctx) {
        // TODO: Multiply all strategy multipliers together
        //       Cap at 3.0x
        //       If changed by > 0.5 from lastSurge, notify observers
        //       Update lastSurge and return
        return 1.0;
    }

    double calculateFare(const PricingContext& ctx) {
        // TODO: Return baseFare * calculateSurge(ctx)
        return ctx.baseFare;
    }
};

// ─── Test Entry Points ────────────────────────────────────────────────────────

double calculateSurge(const PricingContext& ctx) {
    PricingEngine engine;
    engine.addStrategy(new DemandSurge());
    engine.addStrategy(new WeatherSurge());
    engine.addStrategy(new TimeSurge());
    return engine.calculateSurge(ctx);
}

double calculateFare(double baseFare, const PricingContext& ctx) {
    PricingContext c = ctx;
    c.baseFare = baseFare;
    PricingEngine engine;
    engine.addStrategy(new DemandSurge());
    engine.addStrategy(new WeatherSurge());
    engine.addStrategy(new TimeSurge());
    return engine.calculateFare(c);
}

int main() {
    cout << "Ride Surge Pricing — implement the TODO methods above." << endl;
    return 0;
}
