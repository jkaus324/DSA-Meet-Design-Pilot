#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <cmath>
using namespace std;

// --- Data Model (given -- do not modify) ------------------------------------

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

// --- Spot Factory (from Part 1) ---------------------------------------------

class SpotFactory {
public:
    static ParkingSpot createSpot(const string& spotId, int floor, SpotSize size) {
        return {spotId, floor, size, false, ""};
    }
};

// --- Helpers (from Part 1) --------------------------------------------------

SpotSize getMinSpotSize(VehicleType type) {
    switch (type) {
        case VehicleType::MOTORCYCLE: return SpotSize::SMALL;
        case VehicleType::CAR: return SpotSize::MEDIUM;
        case VehicleType::TRUCK: return SpotSize::LARGE;
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
    double calculateFee(long durationSeconds) override {
        // TODO: Return the fixed fee regardless of duration
        return 0.0;
    }
};

// --- Hourly Strategy --------------------------------------------------------

class Hourly : public PricingStrategy {
    double ratePerHour;
public:
    Hourly(double rate) : ratePerHour(rate) {}
    double calculateFee(long durationSeconds) override {
        // TODO: Calculate hours = ceil(durationSeconds / 3600.0)
        // TODO: Return ratePerHour * hours
        return 0.0;
    }
};

// --- Tiered Strategy --------------------------------------------------------

class Tiered : public PricingStrategy {
    double baseRate;    // first hour
    double midRate;     // per hour for hours 1-3
    double highRate;    // per hour beyond 3 hours
public:
    Tiered(double base, double mid, double high)
        : baseRate(base), midRate(mid), highRate(high) {}
    double calculateFee(long durationSeconds) override {
        // TODO: hours = ceil(durationSeconds / 3600.0)
        // TODO: if hours <= 1, return baseRate
        // TODO: if hours <= 3, return baseRate + midRate * (hours - 1)
        // TODO: if hours > 3, return baseRate + midRate * 2 + highRate * (hours - 3)
        return 0.0;
    }
};

// --- Parking Lot (extended from Part 1) -------------------------------------

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
        // TODO: Store the strategy pointer
    }

    void addGate(const string& gateId, GateType type) {
        // TODO: Create a Gate and add it to the gates vector
    }

    vector<string> getGates(GateType type) const {
        // TODO: Return all gate IDs that match the given type
        return {};
    }

    Ticket* parkVehicle(const Vehicle& vehicle, long entryTime, const string& gateId = "") {
        SpotSize minSize = getMinSpotSize(vehicle.type);
        // TODO: Find nearest compatible spot (same as Part 1)
        // TODO: Create Ticket with entryGateId = gateId
        // TODO: Store in activeTickets, return pointer
        return nullptr;
    }

    double unparkVehicle(const string& ticketId, long exitTime, const string& gateId = "") {
        // TODO: Find ticket in activeTickets
        // TODO: Set exitGateId = gateId
        // TODO: Free the spot
        // TODO: Calculate duration = exitTime - entryTime
        // TODO: Use strategy->calculateFee(duration) if strategy is set
        // TODO: Remove ticket, return fee
        return -1.0;
    }

    int getAvailableSpots(SpotSize size) const {
        int count = 0;
        for (auto& floor : floors)
            for (auto& spot : floor)
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
    cout << "Parking Lot Part 2 -- implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
