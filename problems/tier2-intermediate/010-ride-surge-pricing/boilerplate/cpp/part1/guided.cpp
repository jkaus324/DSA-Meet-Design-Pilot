#include <iostream>
#include <vector>
#include <string>
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
// HINT: Each surge factor (demand, weather, time) is an independent strategy.
// They each contribute a multiplier that combines into the final surge.

class /* YourSurgeStrategyName */ {
public:
    virtual double /* yourMultiplierMethod */(const PricingContext& ctx) = 0;
    virtual ~/* YourSurgeStrategyName */() = default;
};

// TODO: Implement concrete surge strategies:
//   - DemandSurge    (based on availableDrivers vs activeRideRequests ratio)
//   - WeatherSurge   (based on weather condition)
//   - TimeSurge      (based on timeOfDay)

// ─── Surge Observer Interface ─────────────────────────────────────────────────
// HINT: These are notified when the surge multiplier changes significantly.

class /* YourObserverName */ {
public:
    virtual void /* onSurgeChange */(double oldMultiplier, double newMultiplier) = 0;
    virtual ~/* YourObserverName */() = default;
};

// TODO: Implement observers:
//   - DriverNotifier  (tells drivers surge is high → opportunity)
//   - RiderNotifier   (warns riders about high surge)

// ─── Pricing Engine ──────────────────────────────────────────────────────────
// TODO: Implement PricingEngine that:
//   - Holds a list of surge strategies
//   - Combines their multipliers (multiply them together, cap at 3.0x)
//   - Notifies observers when surge changes by > 0.5x
//   - Has calculateSurge() and calculateFare() methods

// ─── Test Entry Points ───────────────────────────────────────────────────────
//   double calculateSurge(const PricingContext& ctx);
//   double calculateFare(double baseFare, const PricingContext& ctx);
// ─────────────────────────────────────────────────────────────────────────────

