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

// --- Include your Part 1 ParkingLot logic here ------------------------------

// --- Pricing Strategy Interface ---------------------------------------------
// HINT: The strategy receives duration in seconds and returns a fee (double).
//       The parking lot calls strategy->calculateFee(duration) during unpark.

// class PricingStrategy {
// public:
//     virtual double calculateFee(long durationSeconds) = 0;
//     virtual ~PricingStrategy() = default;
// };

// --- FlatRate ---------------------------------------------------------------
// HINT: Constructor takes a fixed fee. calculateFee ignores duration.

// --- Hourly -----------------------------------------------------------------
// HINT: Constructor takes ratePerHour. Use ceil(seconds / 3600.0) for hours.

// --- Tiered -----------------------------------------------------------------
// HINT: Constructor takes baseRate (first hour), midRate (hours 1-3),
//       highRate (hours 3+). Calculate: base + mid*(hours-1) for 1-3h,
//       base + mid*2 + high*(hours-3) for 3+h.

// --- Gate Management --------------------------------------------------------
// HINT: Store gates in a vector<Gate>. parkVehicle and unparkVehicle
//       accept a gateId parameter. Record it on the Ticket.

// --- Test Entry Points (must exist for tests to compile) --------------------
// Your solution must provide:
//
//   void setPricingStrategy(PricingStrategy* strategy);
//   void addGate(const string& gateId, GateType type);
//   Ticket* parkVehicle(const Vehicle& vehicle, long entryTime, const string& gateId);
//   double unparkVehicle(const string& ticketId, long exitTime, const string& gateId);
//   vector<string> getGates(GateType type);
//
// -------------------------------------------------------------------------


