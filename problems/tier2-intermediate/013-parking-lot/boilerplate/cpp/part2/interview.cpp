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

// --- Your Design Starts Here ------------------------------------------------
//
// Extend your Part 1 ParkingLot to support:
//   1. Pluggable pricing strategies: FlatRate, Hourly, Tiered
//   2. Entry/exit gate registration and tracking
//   3. Gate IDs recorded on tickets
//
// Think about:
//   - How do you define a pricing interface so new strategies can be
//     added without modifying the lot?
//   - How does the Tiered strategy calculate fees across brackets?
//   - Should gates validate their type (entry gates for parking,
//     exit gates for unparking)?
//
// Entry points (must exist for tests):
//   void setPricingStrategy(PricingStrategy* strategy);
//   void addGate(const string& gateId, GateType type);
//   Ticket* parkVehicle(const Vehicle& vehicle, long entryTime, const string& gateId);
//   double unparkVehicle(const string& ticketId, long exitTime, const string& gateId);
//   vector<string> getGates(GateType type);
//
// -------------------------------------------------------------------------


