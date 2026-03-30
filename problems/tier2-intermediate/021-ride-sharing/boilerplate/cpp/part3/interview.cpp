#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct User {
    string id;
    string name;
    int ridesOffered;
    int ridesTaken;
};

struct Vehicle {
    string id;
    string ownerId;
    string model;
    string regNumber;
};

struct Ride {
    string id;
    string driverId;
    string vehicleId;
    string origin;
    string destination;
    int totalSeats;
    int availableSeats;
    bool active;
};

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Extend your RideService to support ending rides and printing statistics.
//
// Think about:
//   - How does ending a ride free the vehicle for future rides?
//   - What if someone tries to end a ride that's already ended?
//   - How do you track per-user statistics efficiently?
//
// Entry points (must exist for tests):
//   All Part 1 and Part 2 entry points PLUS:
//   void RideService::endRide(const string& rideId);
//   vector<pair<string, pair<int, int>>> RideService::getRideStats() const;
//   void RideService::printRideStats() const;
//
// You also need:
//   RideSelectionStrategy interface + concrete strategies from Part 2
//
// ─────────────────────────────────────────────────────────────────────────────


