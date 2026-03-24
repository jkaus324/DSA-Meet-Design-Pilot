#include <iostream>
#include <vector>
#include <string>
using namespace std;

struct RideRequest { string userId, pickup, dropoff, rideType; };
struct Driver { string id; double rating; string rideType; bool available; };
struct PricingContext {
    double baseFare;
    int    availableDrivers;
    int    activeRideRequests;
    string timeOfDay;
    string weather;
};

// ─── NEW in Extension 1 ──────────────────────────────────────────────────────
//
// The ops team wants DYNAMIC SURGE NOTIFICATIONS:
//   - When the surge multiplier changes by more than 0.5x, notify stakeholders
//   - Stakeholders: drivers (to encourage them to go online), ops dashboard
//   - Drivers only receive notifications relevant to their ride type
//
// Think about:
//   - How do you combine the Strategy pattern (surge calculation) with
//     the Observer pattern (surge notifications)?
//   - Does the surge engine become the subject? Or is there a separate notifier?
//   - How do you filter driver notifications by ride type?
//
// Entry points:
//   double calculateSurge(const PricingContext& ctx);
//   double calculateFare(const RideRequest& req, const PricingContext& ctx);
//   void registerSurgeObserver(/* your observer type */);
//   void notifySurgeChange(double oldMultiplier, double newMultiplier,
//                          const string& rideType);
//
// ─────────────────────────────────────────────────────────────────────────────


