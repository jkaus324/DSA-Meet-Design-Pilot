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

class SpotFactory {
public:
    static ParkingSpot createSpot(const string& spotId, int floor, SpotSize size) {
        // TODO: Return a ParkingSpot with isOccupied=false and empty vehicleLicensePlate
        return {"", 0, SpotSize::SMALL, false, ""};
    }
};

// --- Compatibility Helper ---------------------------------------------------

SpotSize getMinSpotSize(VehicleType type) {
    // TODO: Return SMALL for MOTORCYCLE, MEDIUM for CAR, LARGE for TRUCK
    return SpotSize::LARGE;
}

bool isCompatible(SpotSize spotSize, SpotSize minRequired) {
    // TODO: Return true if (int)spotSize >= (int)minRequired
    return false;
}

// --- Parking Lot ------------------------------------------------------------

class ParkingLot {
    vector<vector<ParkingSpot>> floors;
    unordered_map<string, Ticket> activeTickets;
    int nextTicketId;

public:
    ParkingLot(int numFloors) : nextTicketId(1) {
        floors.resize(numFloors);
    }

    void addSpot(int floor, SpotSize size) {
        if (floor < 0 || floor >= (int)floors.size()) return;
        string spotId = "F" + to_string(floor) + "S" + to_string(floors[floor].size());
        // TODO: Use SpotFactory::createSpot to create the spot and add to floors[floor]
    }

    Ticket* parkVehicle(const Vehicle& vehicle, long entryTime) {
        SpotSize minSize = getMinSpotSize(vehicle.type);
        // TODO: Iterate floors in order (0, 1, 2, ...)
        //   For each floor, iterate spots in index order
        //   Find the first unoccupied spot where isCompatible(spot.size, minSize) is true
        //   Mark it occupied, set vehicleLicensePlate
        //   Create a Ticket with id "T" + nextTicketId++, store in activeTickets
        //   Return pointer to the stored ticket
        // TODO: Return nullptr if no compatible spot found
        return nullptr;
    }

    double unparkVehicle(const string& ticketId, long exitTime) {
        // TODO: Find the ticket in activeTickets
        // TODO: Free the corresponding spot (isOccupied=false, clear plate)
        // TODO: Calculate duration = exitTime - entryTime
        // TODO: For Part 1, return duration as fee (1.0 per second)
        // TODO: Remove ticket from activeTickets
        // TODO: Return -1.0 if ticket not found
        return -1.0;
    }

    int getAvailableSpots(SpotSize size) const {
        // TODO: Count all unoccupied spots across all floors matching the exact size
        return 0;
    }

    int getAvailableSpotsByFloor(int floor, SpotSize size) const {
        // TODO: Count unoccupied spots on the given floor matching the exact size
        return 0;
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Parking Lot System -- implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
