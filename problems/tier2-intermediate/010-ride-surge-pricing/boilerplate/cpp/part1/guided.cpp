#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
using namespace std;


// Data class (given).
struct PricingContext {
    double baseFare;
    int availableDrivers;
    int activeRideRequests;
    string timeOfDay;
    string weather;
    PricingContext(double baseFare_, int availableDrivers_, int activeRideRequests_, const string& timeOfDay_, const string& weather_)
      : baseFare(baseFare_), availableDrivers(availableDrivers_), activeRideRequests(activeRideRequests_), timeOfDay(timeOfDay_), weather(weather_) {}
};

struct RideRequest {
    string userId;
    string pickup;
    string dropoff;
    string rideType;
    RideRequest(const string& userId_, const string& pickup_, const string& dropoff_, const string& rideType_)
      : userId(userId_), pickup(pickup_), dropoff(dropoff_), rideType(rideType_) {}
};

// HINT: introduce an abstraction so new ranking rules don't change existing code.
// HINT: keep the comparator small — one rule per class.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
double calculateSurge(PricingContext ctx) {
    // TODO: write your solution
    return {};
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
double calculateFare(RideRequest req, PricingContext ctx) {
    // TODO: write your solution
    return {};
}
