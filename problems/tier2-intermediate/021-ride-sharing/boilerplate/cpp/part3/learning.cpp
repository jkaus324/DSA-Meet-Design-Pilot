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

class RideSelectionStrategy {
public:
    virtual Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                         const string& preference) = 0;
    virtual ~RideSelectionStrategy() = default;
};

class MostVacantStrategy : public RideSelectionStrategy {
public:
    Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                 const string& preference) override {
        // TODO: Find ride with most availableSeats >= seatsNeeded
        return nullptr;
    }
};

class PreferredVehicleStrategy : public RideSelectionStrategy {
    unordered_map<string, Vehicle>& vehicleStore;
public:
    PreferredVehicleStrategy(unordered_map<string, Vehicle>& vs) : vehicleStore(vs) {}

    Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                 const string& preference) override {
        // TODO: Find first ride whose vehicle model matches preference
        return nullptr;
    }
};

// ─── RideService ─────────────────────────────────────────────────────────────

class RideService {
    unordered_map<string, User> users;
    unordered_map<string, Vehicle> vehicles;
    unordered_map<string, Ride> rides;
    unordered_map<string, string> activeVehicles;
    int rideCounter;

public:
    RideService() : rideCounter(0) {}

    void addUser(const string& name) {
        // TODO: Same as Part 1
    }

    void addVehicle(const string& userName, const string& model, const string& regNumber) {
        // TODO: Same as Part 1
    }

    string offerRide(const string& userName, const string& origin,
                     const string& dest, int seats, const string& vehicleRegNumber) {
        // TODO: Same as Part 1
        return "";
    }

    string selectRide(const string& passengerName, const string& origin,
                      const string& dest, int seats,
                      RideSelectionStrategy* strategy, const string& preference = "") {
        // TODO: Same as Part 2
        return "";
    }

    void endRide(const string& rideId) {
        // TODO: Check if ride exists in rides map — return if not

        // TODO: Check if ride is still active — return if already ended (no-op)

        // TODO: Set ride.active = false

        // TODO: Remove vehicle from activeVehicles map
        //   Use: activeVehicles.erase(ride.vehicleId)
    }

    vector<pair<string, pair<int, int>>> getRideStats() const {
        // TODO: Build a vector of {userName, {ridesOffered, ridesTaken}}
        // Iterate through users map and collect stats
        return {};
    }

    void printRideStats() const {
        // TODO: For each user, print:
        //   "User: NAME — Offered: X, Taken: Y"
        // Use cout for output
    }

    // Expose vehicles map for PreferredVehicleStrategy
    unordered_map<string, Vehicle>& getVehicles() { return vehicles; }

    // Accessors for testing
    bool hasUser(const string& name) const { return users.find(name) != users.end(); }
    bool hasVehicle(const string& reg) const { return vehicles.find(reg) != vehicles.end(); }
    bool hasRide(const string& id) const { return rides.find(id) != rides.end(); }
    const User& getUser(const string& name) const { return users.at(name); }
    const Ride& getRide(const string& id) const { return rides.at(id); }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Ride Sharing Part 3 — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
