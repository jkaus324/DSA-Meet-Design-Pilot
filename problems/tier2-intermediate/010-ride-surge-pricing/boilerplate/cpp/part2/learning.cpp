#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;

struct RideRequest { string userId, pickup, dropoff, rideType; };
struct Driver { string id; double rating; string rideType; bool available; };
struct PricingContext { double baseFare; int availableDrivers, activeRideRequests; string timeOfDay, weather; };

class SurgeStrategy {
public:
    virtual double multiplier(const PricingContext& ctx) = 0;
    virtual ~SurgeStrategy() = default;
};

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
        if (ctx.weather == "storm") return 2.0;
        if (ctx.weather == "rain") return 1.3;
        return 1.0;
    }
};

class TimeSurge : public SurgeStrategy {
public:
    double multiplier(const PricingContext& ctx) override {
        if (ctx.timeOfDay == "evening") return 1.4;
        if (ctx.timeOfDay == "morning") return 1.2;
        return 1.0;
    }
};

class SurgeObserver {
public:
    virtual void onSurgeChange(double oldMult, double newMult, const string& rideType) = 0;
    virtual ~SurgeObserver() = default;
};

class DriverObserver : public SurgeObserver {
    Driver driver;
public:
    DriverObserver(Driver d) : driver(d) {}
    void onSurgeChange(double oldMult, double newMult, const string& rideType) override {
        if (driver.rideType == rideType || rideType == "all") {
            cout << "[DRIVER " << driver.id << "] Surge changed: " << oldMult << "x -> " << newMult << "x (" << rideType << ")" << endl;
        }
    }
};

class OpsDashboardObserver : public SurgeObserver {
public:
    void onSurgeChange(double oldMult, double newMult, const string& rideType) override {
        cout << "[OPS] Surge alert for " << rideType << ": " << oldMult << "x -> " << newMult << "x" << endl;
    }
};

class SurgePricingEngine {
    vector<SurgeStrategy*> strategies;
    vector<SurgeObserver*> observers;
    double lastMultiplier = 1.0;
    const double CHANGE_THRESHOLD = 0.5;
public:
    void addStrategy(SurgeStrategy* s) { strategies.push_back(s); }
    void addObserver(SurgeObserver* o) { observers.push_back(o); }

    double calculateSurge(const PricingContext& ctx, const string& rideType = "all") {
        double mult = 1.0;
        for (auto* s : strategies) mult = max(mult, s->multiplier(ctx));
        if (abs(mult - lastMultiplier) > CHANGE_THRESHOLD) {
            for (auto* o : observers) o->onSurgeChange(lastMultiplier, mult, rideType);
        }
        lastMultiplier = mult;
        return mult;
    }
};

static SurgePricingEngine globalEngine;

double calculateSurge(const PricingContext& ctx) {
    return globalEngine.calculateSurge(ctx);
}

double calculateFare(const RideRequest& req, const PricingContext& ctx) {
    return ctx.baseFare * globalEngine.calculateSurge(ctx, req.rideType);
}

void registerSurgeObserver(SurgeObserver* obs) {
    globalEngine.addObserver(obs);
}

int main() {
    cout << "Part 2: Surge notifications — full scaffolding provided." << endl;
    return 0;
}
