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

// ─── Selection Strategy Interface ────────────────────────────────────────────
// HINT: Each strategy picks one ride from a list of candidates.
// HINT: The preference string is used by PreferredVehicle to match model name.

class RideSelectionStrategy {
public:
    virtual Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                         const string& preference) = 0;
    virtual ~RideSelectionStrategy() = default;
};

// ─── Concrete Strategies ─────────────────────────────────────────────────────
// TODO: Implement MostVacantStrategy
//   HINT: Iterate candidates, find the one with max availableSeats >= seatsNeeded

// TODO: Implement PreferredVehicleStrategy
//   HINT: Needs access to vehicle data to resolve vehicleId → model
//   HINT: Pass a reference to the vehicles map in the constructor

// ─── RideService ─────────────────────────────────────────────────────────────
// HINT: selectRide should:
//   1. Filter all active rides matching origin + destination
//   2. Exclude rides where driverId == passengerName
//   3. Pass candidates to strategy->select()
//   4. If selected, decrement availableSeats and increment passenger's ridesTaken
//   5. Return rideId or "" if no match

// class RideService { ... };

