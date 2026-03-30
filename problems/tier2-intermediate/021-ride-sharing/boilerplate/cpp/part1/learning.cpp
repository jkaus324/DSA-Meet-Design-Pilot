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

class RideService {
    unordered_map<string, User> users;
    unordered_map<string, Vehicle> vehicles;          // keyed by regNumber
    unordered_map<string, Ride> rides;                // keyed by rideId
    unordered_map<string, string> activeVehicles;     // regNumber → rideId
    int rideCounter;

public:
    RideService() : rideCounter(0) {}

    void addUser(const string& name) {
        // TODO: Check if user already exists (users.find)
        // TODO: If not, create User{name, name, 0, 0} and store in users map
    }

    void addVehicle(const string& userName, const string& model, const string& regNumber) {
        // TODO: Check if user exists in users map
        // TODO: If user exists, create Vehicle{regNumber, userName, model, regNumber}
        // TODO: Store in vehicles map keyed by regNumber
    }

    string offerRide(const string& userName, const string& origin,
                     const string& dest, int seats, const string& vehicleRegNumber) {
        // TODO: Validate user exists — return "" if not
        // TODO: Validate vehicle exists — return "" if not
        // TODO: Validate vehicle belongs to this user (vehicle.ownerId == userName)
        // TODO: Check activeVehicles map — return "" if vehicle already has active ride
        // TODO: Generate rideId: "RIDE-" + to_string(++rideCounter)
        // TODO: Create Ride{rideId, userName, vehicleRegNumber, origin, dest, seats, seats, true}
        // TODO: Store ride in rides map
        // TODO: Mark vehicle as active: activeVehicles[vehicleRegNumber] = rideId
        // TODO: Increment users[userName].ridesOffered
        // TODO: Return rideId
        return "";
    }

    // Accessors for testing
    bool hasUser(const string& name) const {
        return users.find(name) != users.end();
    }

    bool hasVehicle(const string& regNumber) const {
        return vehicles.find(regNumber) != vehicles.end();
    }

    bool hasRide(const string& rideId) const {
        return rides.find(rideId) != rides.end();
    }

    const User& getUser(const string& name) const {
        return users.at(name);
    }

    const Ride& getRide(const string& rideId) const {
        return rides.at(rideId);
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Ride Sharing Part 1 — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
