# Design Walkthrough — Parking Lot System

> This file is the answer guide. Only read after you've attempted the problem.

---

## The Core Design Decision

Three patterns address three concerns:

1. **How are vehicles and spots created?** -- Different vehicle types map to different spot requirements. A Factory encapsulates this creation logic so the lot doesn't hard-code vehicle-to-spot mappings.
2. **How is the fee calculated?** -- Pricing varies (flat, hourly, tiered). The Strategy pattern makes pricing swappable without touching the lot.
3. **How many lot instances exist?** -- A Singleton ensures all gates, spots, and tickets share one consistent state.

```
ParkingLot (Singleton)
    ├── PricingStrategy* (swappable -- FlatRate / Hourly / Tiered)
    ├── vector<vector<ParkingSpot>> floors (spots organized by floor)
    ├── unordered_map<ticketId, Ticket> activeTickets
    └── vector<Gate> gates (entry/exit gate registry)

parkVehicle(vehicle, entryTime):
    compatibleSize = getMinSpotSize(vehicle.type)
    spot = findNearestAvailableSpot(compatibleSize)  // floor-first, then spotId
    if spot found:
        mark occupied, create Ticket, return it

unparkVehicle(ticketId, exitTime):
    ticket = activeTickets[ticketId]
    duration = exitTime - ticket.entryTime
    fee = strategy->calculateFee(duration)
    free the spot, remove ticket
    return fee
```

---

## Reference Implementation

```cpp
#include <vector>
#include <string>
#include <unordered_map>
#include <algorithm>
#include <cmath>
#include <iostream>
using namespace std;

// --- Data Structures ---

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

// --- Factory ---

class SpotFactory {
public:
    static ParkingSpot createSpot(const string& spotId, int floor, SpotSize size) {
        return {spotId, floor, size, false, ""};
    }
};

// --- Pricing Strategies ---

class PricingStrategy {
public:
    virtual double calculateFee(long durationSeconds) = 0;
    virtual ~PricingStrategy() = default;
};

class FlatRate : public PricingStrategy {
    double fee;
public:
    FlatRate(double fee) : fee(fee) {}
    double calculateFee(long durationSeconds) override {
        return fee;
    }
};

class Hourly : public PricingStrategy {
    double ratePerHour;
public:
    Hourly(double rate) : ratePerHour(rate) {}
    double calculateFee(long durationSeconds) override {
        double hours = ceil((double)durationSeconds / 3600.0);
        return ratePerHour * hours;
    }
};

class Tiered : public PricingStrategy {
    double baseRate;    // first hour
    double midRate;     // per hour for hours 1-3
    double highRate;    // per hour for 3+
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

// --- Helpers ---

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

// --- Parking Lot ---

class ParkingLot {
    vector<vector<ParkingSpot>> floors;
    unordered_map<string, Ticket> activeTickets;
    unordered_map<string, string> plateToTicket; // licensePlate -> ticketId
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
        // Find nearest compatible spot (lowest floor, then lowest spot index)
        for (int f = 0; f < (int)floors.size(); f++) {
            for (int s = 0; s < (int)floors[f].size(); s++) {
                ParkingSpot& spot = floors[f][s];
                if (!spot.isOccupied && isCompatible(spot.size, minSize)) {
                    spot.isOccupied = true;
                    spot.vehicleLicensePlate = vehicle.licensePlate;
                    string tid = "T" + to_string(nextTicketId++);
                    Ticket ticket{tid, vehicle.licensePlate, spot.spotId, f, entryTime, gateId, ""};
                    activeTickets[tid] = ticket;
                    plateToTicket[vehicle.licensePlate] = tid;
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
        for (auto& spot : floors[ticket.floor]) {
            if (spot.spotId == ticket.spotId) {
                spot.isOccupied = false;
                spot.vehicleLicensePlate = "";
                break;
            }
        }
        long duration = exitTime - ticket.entryTime;
        double fee = 0.0;
        if (strategy) {
            fee = strategy->calculateFee(duration);
        }
        plateToTicket.erase(ticket.licensePlate);
        activeTickets.erase(it);
        return fee;
    }

    int getAvailableSpots(SpotSize size) const {
        int count = 0;
        for (auto& floor : floors) {
            for (auto& spot : floor) {
                if (!spot.isOccupied && spot.size == size) count++;
            }
        }
        return count;
    }

    int getAvailableSpotsByFloor(int floor, SpotSize size) const {
        if (floor < 0 || floor >= (int)floors.size()) return 0;
        int count = 0;
        for (auto& spot : floors[floor]) {
            if (!spot.isOccupied && spot.size == size) count++;
        }
        return count;
    }
};
```

---

## What interviewers look for

1. **Factory pattern**: Did you isolate spot/vehicle creation? Can you add a new vehicle type (e.g., Bus) without modifying the lot?
2. **Strategy pattern**: Is pricing fully decoupled? Can you switch from Hourly to Tiered at runtime without changing the lot?
3. **Vehicle-spot compatibility**: Is the mapping clean? Did you use an enum ordering trick (`(int)spotSize >= (int)minRequired`) or a lookup table?
4. **Spot allocation order**: Nearest compatible spot means lowest floor first, then lowest index. Interviewers check that you handle this deterministically.

---

## Common interview follow-ups

- *"How would you handle electric vehicle charging spots?"* -- Add a `bool hasCharger` to ParkingSpot. Extend the factory. Add a preference strategy for EV vehicles.
- *"What if the lot has 10,000 spots?"* -- Maintain a free-list (queue) per spot size per floor instead of scanning all spots linearly.
- *"How do you handle concurrent entry at multiple gates?"* -- Mutex on spot allocation. Or use optimistic locking with atomic compare-and-swap on the spot's isOccupied flag.
- *"How would you add a reservation system?"* -- Add a RESERVED state to spots. Reservations hold a spot for a limited time before converting to OCCUPIED or releasing.
