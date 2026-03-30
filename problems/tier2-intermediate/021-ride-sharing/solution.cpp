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
        Ride* best = nullptr;
        for (auto* ride : candidates) {
            if (ride->availableSeats >= seatsNeeded) {
                if (!best || ride->availableSeats > best->availableSeats) {
                    best = ride;
                }
            }
        }
        return best;
    }
};

class PreferredVehicleStrategy : public RideSelectionStrategy {
    unordered_map<string, Vehicle>& vehicleStore;
public:
    PreferredVehicleStrategy(unordered_map<string, Vehicle>& vs) : vehicleStore(vs) {}

    Ride* select(vector<Ride*>& candidates, int seatsNeeded,
                 const string& preference) override {
        for (auto* ride : candidates) {
            if (ride->availableSeats >= seatsNeeded) {
                auto it = vehicleStore.find(ride->vehicleId);
                if (it != vehicleStore.end() && it->second.model == preference) {
                    return ride;
                }
            }
        }
        return nullptr;
    }
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
        if (users.find(name) != users.end()) return;
        users[name] = User{name, name, 0, 0};
    }

    void addVehicle(const string& userName, const string& model, const string& regNumber) {
        if (users.find(userName) == users.end()) return;
        vehicles[regNumber] = Vehicle{regNumber, userName, model, regNumber};
    }

    string offerRide(const string& userName, const string& origin,
                     const string& dest, int seats, const string& vehicleRegNumber) {
        if (users.find(userName) == users.end()) return "";
        if (vehicles.find(vehicleRegNumber) == vehicles.end()) return "";
        if (vehicles[vehicleRegNumber].ownerId != userName) return "";
        if (activeVehicles.find(vehicleRegNumber) != activeVehicles.end()) return "";

        string rideId = "RIDE-" + to_string(++rideCounter);
        rides[rideId] = Ride{rideId, userName, vehicleRegNumber, origin, dest, seats, seats, true};
        activeVehicles[vehicleRegNumber] = rideId;
        users[userName].ridesOffered++;
        return rideId;
    }

    string selectRide(const string& passengerName, const string& origin,
                      const string& dest, int seats,
                      RideSelectionStrategy* strategy, const string& preference = "") {
        if (users.find(passengerName) == users.end()) return "";

        vector<Ride*> candidates;
        for (auto& [id, ride] : rides) {
            if (ride.active && ride.origin == origin && ride.destination == dest
                && ride.availableSeats >= seats && ride.driverId != passengerName) {
                candidates.push_back(&ride);
            }
        }

        Ride* selected = strategy->select(candidates, seats, preference);
        if (selected) {
            selected->availableSeats -= seats;
            users[passengerName].ridesTaken++;
            return selected->id;
        }
        return "";
    }

    void endRide(const string& rideId) {
        auto it = rides.find(rideId);
        if (it == rides.end()) return;
        if (!it->second.active) return;
        it->second.active = false;
        activeVehicles.erase(it->second.vehicleId);
    }

    vector<pair<string, pair<int, int>>> getRideStats() const {
        vector<pair<string, pair<int, int>>> stats;
        for (const auto& [name, user] : users) {
            stats.push_back({name, {user.ridesOffered, user.ridesTaken}});
        }
        return stats;
    }

    void printRideStats() const {
        for (const auto& [name, user] : users) {
            cout << "User: " << name << " — Offered: " << user.ridesOffered
                 << ", Taken: " << user.ridesTaken << endl;
        }
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
    cout << "Ride Sharing — Full Solution" << endl;
    return 0;
}
#endif
