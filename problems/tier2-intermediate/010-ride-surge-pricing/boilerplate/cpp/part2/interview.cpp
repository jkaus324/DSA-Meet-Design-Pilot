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

// TODO: design and implement your solution.
// Required free functions:
//   double calculateSurge(PricingContext ctx);
//   double calculateFare(RideRequest req, PricingContext ctx);

double calculateSurge(PricingContext ctx) {
    // TODO: write your solution
    return {};
}

double calculateFare(RideRequest req, PricingContext ctx) {
    // TODO: write your solution
    return {};
}
