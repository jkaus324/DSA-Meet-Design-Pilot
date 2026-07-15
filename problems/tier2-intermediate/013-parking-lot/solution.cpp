#include <iostream>
#include <memory>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <cmath>
using namespace std;

// --- Data Model -------------------------------------------------------------

enum class VehicleType {
    MOTORCYCLE,
    CAR,
    TRUCK
};

enum class SpotSize {
    SMALL,
    MEDIUM,
    LARGE
};

enum class GateType {
    ENTRY,
    EXIT
};

struct Vehicle {
    string licensePlate;
    VehicleType type;
};

struct ParkingSpot {
    string spotId;
    int floor;
    SpotSize size;
    bool isOccupied;
    string vehicleLicensePlate;
};

struct Ticket {
    string ticketId;
    string licensePlate;
    string spotId;
    int floor;
    long entryTime;
    string entryGateId;
    string exitGateId;
};

struct Gate {
    string gateId;
    GateType type;
};

// --- Spot Factory -----------------------------------------------------------

class SpotFactory {
public:
    static ParkingSpot createSpot(const string& spotId, int floor, SpotSize size) {
        return {spotId, floor, size, false, ""};
    }
};

// --- Compatibility Helpers --------------------------------------------------

SpotSize getMinSpotSize(VehicleType type) {
    switch (type) {
        case VehicleType::MOTORCYCLE: return SpotSize::SMALL;
        case VehicleType::CAR:        return SpotSize::MEDIUM;
        case VehicleType::TRUCK:      return SpotSize::LARGE;
    }
    return SpotSize::LARGE;
}

bool isCompatible(SpotSize spotSize, SpotSize minRequired) {
    return (int)spotSize >= (int)minRequired;
}

// --- Pricing Strategy Interface ---------------------------------------------

class PricingStrategy {
public:
    virtual double calculateFee(long durationSeconds) = 0;
    virtual ~PricingStrategy() = default;
};

// --- FlatRate Strategy ------------------------------------------------------

class FlatRate : public PricingStrategy {
    double fee;
public:
    FlatRate(double fee) : fee(fee) {}
    double calculateFee(long /*durationSeconds*/) override {
        return fee;
    }
};

// --- Hourly Strategy --------------------------------------------------------

class Hourly : public PricingStrategy {
    double ratePerHour;
public:
    Hourly(double rate) : ratePerHour(rate) {}
    double calculateFee(long durationSeconds) override {
        double hours = ceil((double)durationSeconds / 3600.0);
        return ratePerHour * hours;
    }
};

// --- Tiered Strategy --------------------------------------------------------

class Tiered : public PricingStrategy {
    double baseRate;
    double midRate;
    double highRate;
public:
    Tiered(double base, double mid, double high)
        : baseRate(base), midRate(mid), highRate(high) {}
    double calculateFee(long durationSeconds) override {
        double hours = ceil((double)durationSeconds / 3600.0);
        if (hours <= 1) return baseRate;
        if (hours <= 3) return baseRate + midRate * (hours - 1);
        return baseRate + midRate * 2 + highRate * (hours - 3);
    }
};

// --- Parking Lot ------------------------------------------------------------

class ParkingLot {
    vector<vector<ParkingSpot>> floors;
    unordered_map<string, Ticket> activeTickets;
    vector<Gate> gates;
    PricingStrategy* strategy;
    int nextTicketId;

public:
    ParkingLot(int numFloors) : strategy(nullptr), nextTicketId(1) {
        floors.resize(numFloors);
    }

    void addSpot(int floor, SpotSize size) {
        if (floor < 0 || floor >= (int)floors.size()) return;
        string spotId = "F" + to_string(floor) + "S" + to_string(floors[floor].size());
        floors[floor].push_back(SpotFactory::createSpot(spotId, floor, size));
    }

    void setPricingStrategy(PricingStrategy* s) {
        strategy = s;
    }

    void addGate(const string& gateId, GateType type) {
        gates.push_back({gateId, type});
    }

    vector<string> getGates(GateType type) const {
        vector<string> result;
        for (auto& g : gates) {
            if (g.type == type) result.push_back(g.gateId);
        }
        return result;
    }

    Ticket* parkVehicle(const Vehicle& vehicle, long entryTime, const string& gateId = "") {
        SpotSize minSize = getMinSpotSize(vehicle.type);
        for (int f = 0; f < (int)floors.size(); f++) {
            for (int s = 0; s < (int)floors[f].size(); s++) {
                ParkingSpot& spot = floors[f][s];
                if (!spot.isOccupied && isCompatible(spot.size, minSize)) {
                    spot.isOccupied = true;
                    spot.vehicleLicensePlate = vehicle.licensePlate;
                    string tid = "T" + to_string(nextTicketId++);
                    Ticket ticket;
                    ticket.ticketId = tid;
                    ticket.licensePlate = vehicle.licensePlate;
                    ticket.spotId = spot.spotId;
                    ticket.floor = f;
                    ticket.entryTime = entryTime;
                    ticket.entryGateId = gateId;
                    ticket.exitGateId = "";
                    activeTickets[tid] = ticket;
                    return &activeTickets[tid];
                }
            }
        }
        return nullptr;
    }

    double unparkVehicle(const string& ticketId, long exitTime, const string& gateId = "") {
        auto it = activeTickets.find(ticketId);
        if (it == activeTickets.end()) return -1.0;

        Ticket& ticket = it->second;
        ticket.exitGateId = gateId;

        // Free the spot
        for (auto& floorSpots : floors) {
            for (auto& spot : floorSpots) {
                if (spot.spotId == ticket.spotId && spot.isOccupied) {
                    spot.isOccupied = false;
                    spot.vehicleLicensePlate = "";
                    break;
                }
            }
        }

        long duration = exitTime - ticket.entryTime;
        double fee;
        if (strategy) {
            fee = strategy->calculateFee(duration);
        } else {
            fee = (double)duration; // 1.0 per second for Part 1
        }

        activeTickets.erase(it);
        return fee;
    }

    int getAvailableSpots(SpotSize size) const {
        int count = 0;
        for (auto& floorSpots : floors)
            for (auto& spot : floorSpots)
                if (!spot.isOccupied && spot.size == size) count++;
        return count;
    }

    int getAvailableSpotsByFloor(int floor, SpotSize size) const {
        if (floor < 0 || floor >= (int)floors.size()) return 0;
        int count = 0;
        for (auto& spot : floors[floor])
            if (!spot.isOccupied && spot.size == size) count++;
        return count;
    }
};

// --- Ops simulator (used by spec-based tests) -------------------------------
//
// Drives one ParkingLot through a sequence of operations.
//
// Op fields used:
//   "new"            i1=numFloors                                 -> "ok"
//   "add_spot"       i1=floor, s1=size("S"|"M"|"L")                -> "ok"
//   "add_gate"       s1=gateId, s2="entry"|"exit"                  -> "ok"
//   "gates_count"    s1="entry"|"exit"                             -> int as string
//   "gate_at"        s1="entry"|"exit", i1=index                   -> gate id
//   "set_pricing"    s1="flat"|"hourly"|"tiered"  i1=base i2=mid i3=high -> "ok"
//   "park"           s1=licensePlate, s2=type("M"|"C"|"T"), i1=entryTime, s3=gate
//                    => stores last ticketId in slot i2 (default 0)
//                    -> ticket id ("" on failure)
//   "ticket_at"      i1=slot                                       -> ticket id stored
//   "ticket_floor"   i1=slot                                       -> int as string
//   "ticket_entry"   i1=slot                                       -> entryGateId
//   "ticket_spot_id" i1=slot                                       -> spot id
//   "unpark"         i1=slot, i2=exitTime, s1=gateId               -> fee with 2 decimals or "-1"
//   "unpark_id"      s1=ticketId, i2=exitTime, s2=gateId           -> fee or "-1"
//   "available"      s1=size                                       -> int as string
//   "available_floor" i1=floor, s1=size                            -> int as string

struct ParkOp {
    string kind;
    string s1;
    string s2;
    string s3;
    int    i1;
    int    i2;
    int    i3;
};

struct TicketSnap {
    string id;
    int    floor;
    string spotId;
    string entryGate;
};

static SpotSize size_from(const string& s) {
    if (s == "S" || s == "small")  return SpotSize::SMALL;
    if (s == "M" || s == "medium") return SpotSize::MEDIUM;
    return SpotSize::LARGE;
}
static VehicleType vtype_from(const string& s) {
    if (s == "M" || s == "moto" || s == "motorcycle") return VehicleType::MOTORCYCLE;
    if (s == "C" || s == "car")                       return VehicleType::CAR;
    return VehicleType::TRUCK;
}
static GateType gate_from(const string& s) {
    return s == "entry" ? GateType::ENTRY : GateType::EXIT;
}
static string fee_to_str(double f) {
    if (f < 0) return "-1";
    char buf[32];
    snprintf(buf, sizeof(buf), "%.2f", f);
    return string(buf);
}

vector<string> parking_simulate(vector<ParkOp> ops) {
    vector<string> out;
    unique_ptr<ParkingLot> lot;
    vector<unique_ptr<PricingStrategy>> strategies; // keep alive
    vector<string> tickets(16, "");
    vector<TicketSnap> snaps(16);
    for (const auto& op : ops) {
        const string& k = op.kind;
        if (k == "new") {
            lot.reset(new ParkingLot(op.i1));
            strategies.clear();
            for (auto& t : tickets) t.clear();
            out.push_back("ok");
        } else if (k == "add_spot") {
            lot->addSpot(op.i1, size_from(op.s1));
            out.push_back("ok");
        } else if (k == "add_gate") {
            lot->addGate(op.s1, gate_from(op.s2));
            out.push_back("ok");
        } else if (k == "gates_count") {
            out.push_back(to_string((int)lot->getGates(gate_from(op.s1)).size()));
        } else if (k == "gate_at") {
            auto g = lot->getGates(gate_from(op.s1));
            out.push_back(op.i1 >= 0 && op.i1 < (int)g.size() ? g[op.i1] : "");
        } else if (k == "set_pricing") {
            PricingStrategy* p = nullptr;
            if (op.s1 == "flat")        { p = new FlatRate((double)op.i1); }
            else if (op.s1 == "hourly") { p = new Hourly((double)op.i1); }
            else if (op.s1 == "tiered") { p = new Tiered((double)op.i1, (double)op.i2, (double)op.i3); }
            if (p) {
                strategies.emplace_back(p);
                lot->setPricingStrategy(p);
            }
            out.push_back("ok");
        } else if (k == "park") {
            Vehicle v{op.s1, vtype_from(op.s2)};
            Ticket* t = lot->parkVehicle(v, (long)op.i1, op.s3);
            if (op.i2 >= 0 && op.i2 < (int)tickets.size()) {
                tickets[op.i2] = t ? t->ticketId : "";
                snaps[op.i2] = t ? TicketSnap{t->ticketId, t->floor, t->spotId, t->entryGateId}
                                 : TicketSnap{"", -1, "", ""};
            }
            out.push_back(t ? t->ticketId : "");
        } else if (k == "ticket_at") {
            out.push_back(op.i1 >= 0 && op.i1 < (int)tickets.size() ? tickets[op.i1] : "");
        } else if (k == "ticket_floor") {
            out.push_back(op.i1 >= 0 && op.i1 < (int)snaps.size() ? to_string(snaps[op.i1].floor) : "-1");
        } else if (k == "ticket_spot_id") {
            out.push_back(op.i1 >= 0 && op.i1 < (int)snaps.size() ? snaps[op.i1].spotId : "");
        } else if (k == "ticket_entry") {
            out.push_back(op.i1 >= 0 && op.i1 < (int)snaps.size() ? snaps[op.i1].entryGate : "");
        } else if (k == "unpark") {
            string tid = op.i1 >= 0 && op.i1 < (int)tickets.size() ? tickets[op.i1] : "";
            double fee = lot->unparkVehicle(tid, (long)op.i2, op.s1);
            out.push_back(fee_to_str(fee));
        } else if (k == "unpark_id") {
            double fee = lot->unparkVehicle(op.s1, (long)op.i2, op.s2);
            out.push_back(fee_to_str(fee));
        } else if (k == "available") {
            out.push_back(to_string(lot->getAvailableSpots(size_from(op.s1))));
        } else if (k == "available_floor") {
            out.push_back(to_string(lot->getAvailableSpotsByFloor(op.i1, size_from(op.s1))));
        } else {
            out.push_back("unknown:" + k);
        }
    }
    return out;
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Parking Lot System -- run tests to verify." << endl;
    return 0;
}
#endif
