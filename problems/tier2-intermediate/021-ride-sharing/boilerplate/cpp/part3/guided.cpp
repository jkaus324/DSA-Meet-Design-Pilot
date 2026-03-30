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

// ─── Selection Strategy Interface (from Part 2) ─────────────────────────────

class RideSelectionStrategy {
public:
    virtual Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                         const string& preference) = 0;
    virtual ~RideSelectionStrategy() = default;
};

// TODO: Implement MostVacantStrategy and PreferredVehicleStrategy (same as Part 2)

// ─── RideService ─────────────────────────────────────────────────────────────
// HINT: endRide should:
//   1. Validate ride exists
//   2. Check if ride is still active (no-op if already ended)
//   3. Set ride.active = false
//   4. Remove vehicle from activeVehicles map (frees it for future rides)

// HINT: getRideStats should iterate the users map and return
//   a vector of {name, {ridesOffered, ridesTaken}} pairs

// HINT: printRideStats should format as:
//   "User: NAME — Offered: X, Taken: Y"

// class RideService { ... };

