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
// Extend your RideService to support pluggable ride selection strategies:
//   - MostVacant: Select the ride with the most available seats
//   - PreferredVehicle: Select the ride whose vehicle model matches preference
//
// Think about:
//   - What abstraction lets you swap selection logic at runtime?
//   - How do you filter candidates (origin, dest, active, enough seats)?
//   - How does PreferredVehicleStrategy access vehicle model information?
//   - What happens when no ride matches the criteria?
//
// Entry points (must exist for tests):
//   All Part 1 entry points PLUS:
//   string RideService::selectRide(const string& passengerName,
//                                  const string& origin, const string& dest,
//                                  int seats, RideSelectionStrategy* strategy,
//                                  const string& preference);
//
// You also need:
//   RideSelectionStrategy interface
//   MostVacantStrategy
//   PreferredVehicleStrategy
//
// ─────────────────────────────────────────────────────────────────────────────


