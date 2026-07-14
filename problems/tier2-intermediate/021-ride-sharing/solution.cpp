#include <iostream>
#include <memory>
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

// ─── Ops simulator (used by spec-based tests) ──────────────────────────────
//
// Drives one RideService through a sequence of ops. ride_*  ops store rideId
// in slot i1; later ops can refer to the slot.
//
// Op fields:
//   "new"                                                              -> "ok"
//   "add_user"   s1=name                                                -> "ok"
//   "add_veh"    s1=user s2=model s3=reg                                -> "ok"
//   "offer"      s1=user s2=origin s3=dest s4=reg i1=seats i2=slot       -> rideId or ""
//   "ride_active" i2=slot                                                -> "yes"/"no"
//   "ride_origin" i2=slot                                                -> origin string
//   "ride_dest"   i2=slot                                                -> dest string
//   "ride_total"  i2=slot                                                -> int
//   "ride_avail"  i2=slot                                                -> int
//   "ride_driver" i2=slot                                                -> driver name
//   "select_mv"   s1=passenger s2=origin s3=dest i1=seats i2=outSlot      -> rideId or ""
//   "select_pv"   s1=passenger s2=origin s3=dest s4=preferredModel i1=seats i2=outSlot -> rideId or ""
//   "end"         i2=slot                                                 -> "ok"
//   "end_id"      s1=rideId                                               -> "ok"
//   "user_offered" s1=name                                                -> int
//   "user_taken"   s1=name                                                -> int
//   "has_user"    s1=name                                                 -> "yes"/"no"
//   "has_vehicle" s1=reg                                                  -> "yes"/"no"
//   "has_ride"    i2=slot                                                 -> "yes"/"no"

struct RideOp {
    string kind;
    string s1;
    string s2;
    string s3;
    string s4;
    int    i1;
    int    i2;
};

vector<string> ride_simulate(vector<RideOp> ops) {
    vector<string> out;
    unique_ptr<RideService> svc(new RideService());
    vector<string> rideSlots(32, "");
    for (const auto& op : ops) {
        const string& k = op.kind;
        if (k == "new") {
            svc.reset(new RideService());
            for (auto& s : rideSlots) s.clear();
            out.push_back("ok");
        } else if (k == "add_user") {
            svc->addUser(op.s1); out.push_back("ok");
        } else if (k == "add_veh") {
            svc->addVehicle(op.s1, op.s2, op.s3); out.push_back("ok");
        } else if (k == "offer") {
            string rid = svc->offerRide(op.s1, op.s2, op.s3, op.i1, op.s4);
            if (op.i2 >= 0 && op.i2 < (int)rideSlots.size()) rideSlots[op.i2] = rid;
            out.push_back(rid);
        } else if (k == "ride_active") {
            const string& rid = rideSlots[op.i2];
            out.push_back(svc->hasRide(rid) && svc->getRide(rid).active ? "yes" : "no");
        } else if (k == "ride_origin") {
            out.push_back(svc->hasRide(rideSlots[op.i2]) ? svc->getRide(rideSlots[op.i2]).origin : "");
        } else if (k == "ride_dest") {
            out.push_back(svc->hasRide(rideSlots[op.i2]) ? svc->getRide(rideSlots[op.i2]).destination : "");
        } else if (k == "ride_total") {
            out.push_back(to_string(svc->hasRide(rideSlots[op.i2]) ? svc->getRide(rideSlots[op.i2]).totalSeats : -1));
        } else if (k == "ride_avail") {
            out.push_back(to_string(svc->hasRide(rideSlots[op.i2]) ? svc->getRide(rideSlots[op.i2]).availableSeats : -1));
        } else if (k == "ride_driver") {
            out.push_back(svc->hasRide(rideSlots[op.i2]) ? svc->getRide(rideSlots[op.i2]).driverId : "");
        } else if (k == "select_mv") {
            MostVacantStrategy mv;
            string rid = svc->selectRide(op.s1, op.s2, op.s3, op.i1, &mv);
            if (op.i2 >= 0 && op.i2 < (int)rideSlots.size()) rideSlots[op.i2] = rid;
            out.push_back(rid);
        } else if (k == "select_pv") {
            PreferredVehicleStrategy pv(svc->getVehicles());
            string rid = svc->selectRide(op.s1, op.s2, op.s3, op.i1, &pv, op.s4);
            if (op.i2 >= 0 && op.i2 < (int)rideSlots.size()) rideSlots[op.i2] = rid;
            out.push_back(rid);
        } else if (k == "end") {
            svc->endRide(rideSlots[op.i2]); out.push_back("ok");
        } else if (k == "end_id") {
            svc->endRide(op.s1); out.push_back("ok");
        } else if (k == "user_offered") {
            out.push_back(svc->hasUser(op.s1) ? to_string(svc->getUser(op.s1).ridesOffered) : "0");
        } else if (k == "user_taken") {
            out.push_back(svc->hasUser(op.s1) ? to_string(svc->getUser(op.s1).ridesTaken) : "0");
        } else if (k == "has_user") {
            out.push_back(svc->hasUser(op.s1) ? "yes" : "no");
        } else if (k == "has_vehicle") {
            out.push_back(svc->hasVehicle(op.s1) ? "yes" : "no");
        } else if (k == "has_ride") {
            out.push_back(svc->hasRide(rideSlots[op.i2]) ? "yes" : "no");
        } else {
            out.push_back("unknown:" + k);
        }
    }
    return out;
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Ride Sharing — Full Solution" << endl;
    return 0;
}
#endif
