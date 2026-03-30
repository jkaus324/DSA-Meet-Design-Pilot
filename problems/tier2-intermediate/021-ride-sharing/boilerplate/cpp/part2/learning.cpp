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

// ─── Concrete Strategies ─────────────────────────────────────────────────────

class MostVacantStrategy : public RideSelectionStrategy {
public:
    Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                 const string& preference) override {
        // TODO: Find the ride with the most availableSeats (>= seatsNeeded)
        // Iterate through candidates, track the best (most seats)
        // Return nullptr if no candidate has enough seats
        return nullptr;
    }
};

class PreferredVehicleStrategy : public RideSelectionStrategy {
    unordered_map<string, Vehicle>& vehicleStore;
public:
    PreferredVehicleStrategy(unordered_map<string, Vehicle>& vs) : vehicleStore(vs) {}

    Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                 const string& preference) override {
        // TODO: Find the first ride whose vehicle model matches preference
        // Use vehicleStore[ride->vehicleId].model to get the model
        // Only consider rides with availableSeats >= seatsNeeded
        // Return nullptr if no match found
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
        // TODO: Validate passenger exists — return "" if not

        // TODO: Build candidate list:
        //   For each ride in rides map, include if:
        //     - ride.active == true
        //     - ride.origin == origin
        //     - ride.destination == dest
        //     - ride.availableSeats >= seats
        //     - ride.driverId != passengerName (can't select own ride)

        // TODO: Call strategy->select(candidates, seats, preference)

        // TODO: If selected:
        //   - Decrement selected->availableSeats by seats
        //   - Increment users[passengerName].ridesTaken
        //   - Return selected->id

        // TODO: If no match, return ""
        return "";
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
    cout << "Ride Sharing Part 2 — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
