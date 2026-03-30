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

// ─── RideService ─────────────────────────────────────────────────────────────
// HINT: Use unordered_map for O(1) lookups:
//   - users keyed by name
//   - vehicles keyed by regNumber
//   - rides keyed by rideId
// HINT: Track active vehicles with a separate map: regNumber → rideId
//   This makes "is this vehicle in an active ride?" an O(1) check

class RideService {
    // HINT: Declare your maps here
    // HINT: Use an int counter for generating ride IDs like "RIDE-1", "RIDE-2"

public:
    // HINT: addUser creates a User with ridesOffered=0, ridesTaken=0
    // HINT: Check if user already exists before adding
    void addUser(const string& name);

    // HINT: addVehicle creates a Vehicle linked to the user
    // HINT: Validate that the user exists before adding
    void addVehicle(const string& userName, const string& model, const string& regNumber);

    // HINT: offerRide should:
    //   1. Validate user exists
    //   2. Validate vehicle exists and belongs to user
    //   3. Check vehicle doesn't have an active ride (use activeVehicles map)
    //   4. Create ride, mark vehicle as active, increment ridesOffered
    //   5. Return rideId or "" on failure
    string offerRide(const string& userName, const string& origin,
                     const string& dest, int seats, const string& vehicleRegNumber);
};

