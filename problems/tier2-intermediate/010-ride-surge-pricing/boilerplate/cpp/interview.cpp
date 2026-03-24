#include <iostream>
#include <vector>
#include <string>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Design and implement a Surge Pricing Engine that:
//   1. Calculates surge multiplier based on supply/demand and conditions
//   2. Lets multiple surge factors (time, weather, demand) combine independently
//   3. Notifies relevant parties when surge changes significantly
//
// Think about:
//   - How do you combine multiple surge factors without one giant if-else chain?
//   - How would you add a "special event" surge factor with zero changes to existing code?
//   - Who needs to be notified when surge changes? How do they subscribe?
//
// Entry points:
//   double calculateSurge(const PricingContext& ctx);
//   double calculateFare(const RideRequest& req, const PricingContext& ctx);
//
// ─────────────────────────────────────────────────────────────────────────────

