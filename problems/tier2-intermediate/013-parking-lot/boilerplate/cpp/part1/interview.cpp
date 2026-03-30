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

// --- Your Design Starts Here ------------------------------------------------
//
// Design and implement a ParkingLot that:
//   1. Has multiple floors, each with spots of various sizes
//   2. Parks a vehicle in the nearest compatible spot (lowest floor first,
//      then lowest spot index)
//   3. Unparks a vehicle given a ticket ID and returns the parking fee
//   4. Reports available spot counts
//
// Compatibility:
//   MOTORCYCLE -> SMALL or larger
//   CAR        -> MEDIUM or larger
//   TRUCK      -> LARGE only
//
// Think about:
//   - How do you map vehicle types to minimum spot sizes?
//   - What data structure organizes spots by floor for nearest-first allocation?
//   - Should spot creation logic be in the lot, or in a separate factory?
//
// Entry points (must exist for tests):
//   Ticket* parkVehicle(const Vehicle& vehicle, long entryTime);
//   double unparkVehicle(const string& ticketId, long exitTime);
//   int getAvailableSpots(SpotSize size);
//   int getAvailableSpotsByFloor(int floor, SpotSize size);
//
// -------------------------------------------------------------------------


