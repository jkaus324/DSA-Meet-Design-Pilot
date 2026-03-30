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
};

// --- Spot Factory -----------------------------------------------------------
// HINT: Encapsulate spot creation. The factory takes an ID, floor, and size,
//       and returns a ParkingSpot with isOccupied=false.

// class SpotFactory {
// public:
//     static ParkingSpot createSpot(const string& spotId, int floor, SpotSize size);
// };

// --- Compatibility Helper ---------------------------------------------------
// HINT: Map VehicleType to minimum SpotSize. A vehicle can park in any spot
//       whose size >= its minimum. Use the enum's underlying int for comparison:
//       (int)SpotSize::SMALL=0, MEDIUM=1, LARGE=2.

// --- Parking Lot ------------------------------------------------------------
// HINT: Use vector<vector<ParkingSpot>> for floors. Scan floor 0 first,
//       then floor 1, etc. Within a floor, scan spots in index order.
// HINT: Use unordered_map<string, Ticket> to track active tickets by ID.
// HINT: Generate ticket IDs incrementally: "T1", "T2", etc.

// class ParkingLot {
// private:
//     vector<vector<ParkingSpot>> floors;
//     unordered_map<string, Ticket> activeTickets;
//     int nextTicketId;
// public:
//     ParkingLot(int numFloors);
//     void addSpot(int floor, SpotSize size);
//     Ticket* parkVehicle(const Vehicle& vehicle, long entryTime);
//     double unparkVehicle(const string& ticketId, long exitTime);
//     int getAvailableSpots(SpotSize size);
//     int getAvailableSpotsByFloor(int floor, SpotSize size);
// };

// --- Test Entry Points (must exist for tests to compile) --------------------
// Your solution must provide these via a ParkingLot instance:
//
//   Ticket* parkVehicle(const Vehicle& vehicle, long entryTime);
//   double unparkVehicle(const string& ticketId, long exitTime);
//   int getAvailableSpots(SpotSize size);
//   int getAvailableSpotsByFloor(int floor, SpotSize size);
//
// -------------------------------------------------------------------------


