#include <iostream>
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

#ifndef RUNNING_TESTS
int main() {
    cout << "Parking Lot System -- run tests to verify." << endl;
    return 0;
}
#endif
