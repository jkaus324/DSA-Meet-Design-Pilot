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
// Design and implement a RideService that:
//   1. Registers users by name (name is the unique ID)
//   2. Registers vehicles belonging to a user
//   3. Allows users to offer rides (origin, dest, seats, vehicle)
//   4. Prevents a vehicle from having multiple active rides
//
// Think about:
//   - What data structures give O(1) lookup for users, vehicles, rides?
//   - How do you check if a vehicle already has an active ride?
//   - What happens if a user tries to offer a ride with someone else's vehicle?
//
// Entry points (must exist for tests):
//   void RideService::addUser(const string& name);
//   void RideService::addVehicle(const string& userName, const string& model,
//                                const string& regNumber);
//   string RideService::offerRide(const string& userName, const string& origin,
//                                 const string& dest, int seats,
//                                 const string& vehicleRegNumber);
//
// ─────────────────────────────────────────────────────────────────────────────


