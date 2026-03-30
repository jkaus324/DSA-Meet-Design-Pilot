// Part 1 Tests -- Multi-floor Parking with Vehicle-Spot Matching
// Tests parking, unparking, spot compatibility, and availability counts

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Park a car in a medium spot
    try {
        ParkingLot lot(2);
        lot.addSpot(0, SpotSize::MEDIUM);
        Vehicle car{"ABC123", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 1000);
        assert(t != nullptr);
        assert(t->licensePlate == "ABC123");
        assert(t->floor == 0);
        cout << "PASS test_park_car" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_park_car" << endl;
        failed++;
    }

    // Test 2: Park a motorcycle in a small spot
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::SMALL);
        Vehicle moto{"MOTO1", VehicleType::MOTORCYCLE};
        Ticket* t = lot.parkVehicle(moto, 1000);
        assert(t != nullptr);
        assert(t->licensePlate == "MOTO1");
        cout << "PASS test_park_motorcycle" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_park_motorcycle" << endl;
        failed++;
    }

    // Test 3: Car cannot park in small spot
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::SMALL);
        Vehicle car{"CAR1", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 1000);
        assert(t == nullptr);
        cout << "PASS test_car_no_small_spot" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_car_no_small_spot" << endl;
        failed++;
    }

    // Test 4: Motorcycle can park in a larger spot when small is unavailable
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        Vehicle moto{"MOTO2", VehicleType::MOTORCYCLE};
        Ticket* t = lot.parkVehicle(moto, 1000);
        assert(t != nullptr);
        cout << "PASS test_motorcycle_in_medium_spot" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_motorcycle_in_medium_spot" << endl;
        failed++;
    }

    // Test 5: Truck requires large spot
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        lot.addSpot(0, SpotSize::LARGE);
        Vehicle truck{"TRUCK1", VehicleType::TRUCK};
        Ticket* t = lot.parkVehicle(truck, 1000);
        assert(t != nullptr);
        assert(t->spotId.find("S1") != string::npos); // second spot (large)
        cout << "PASS test_truck_large_spot" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_truck_large_spot" << endl;
        failed++;
    }

    // Test 6: Unpark returns a fee and frees the spot
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        Vehicle car{"CAR2", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 1000);
        assert(t != nullptr);
        string tid = t->ticketId;
        assert(lot.getAvailableSpots(SpotSize::MEDIUM) == 0);
        double fee = lot.unparkVehicle(tid, 4600); // 3600 seconds = 1 hour
        assert(fee >= 0);
        assert(lot.getAvailableSpots(SpotSize::MEDIUM) == 1);
        cout << "PASS test_unpark_frees_spot" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_unpark_frees_spot" << endl;
        failed++;
    }

    // Test 7: Invalid ticket returns -1
    try {
        ParkingLot lot(1);
        double fee = lot.unparkVehicle("INVALID", 5000);
        assert(fee < 0);
        cout << "PASS test_invalid_ticket" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_invalid_ticket" << endl;
        failed++;
    }

    // Test 8: Available spots count is correct
    try {
        ParkingLot lot(2);
        lot.addSpot(0, SpotSize::SMALL);
        lot.addSpot(0, SpotSize::MEDIUM);
        lot.addSpot(0, SpotSize::MEDIUM);
        lot.addSpot(1, SpotSize::LARGE);
        assert(lot.getAvailableSpots(SpotSize::SMALL) == 1);
        assert(lot.getAvailableSpots(SpotSize::MEDIUM) == 2);
        assert(lot.getAvailableSpots(SpotSize::LARGE) == 1);
        cout << "PASS test_available_spots_count" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_available_spots_count" << endl;
        failed++;
    }

    // Test 9: Available spots by floor
    try {
        ParkingLot lot(2);
        lot.addSpot(0, SpotSize::MEDIUM);
        lot.addSpot(0, SpotSize::MEDIUM);
        lot.addSpot(1, SpotSize::MEDIUM);
        assert(lot.getAvailableSpotsByFloor(0, SpotSize::MEDIUM) == 2);
        assert(lot.getAvailableSpotsByFloor(1, SpotSize::MEDIUM) == 1);
        cout << "PASS test_available_spots_by_floor" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_available_spots_by_floor" << endl;
        failed++;
    }

    // Test 10: Nearest spot allocation — prefers lower floor
    try {
        ParkingLot lot(3);
        lot.addSpot(0, SpotSize::SMALL);  // floor 0 has only small
        lot.addSpot(1, SpotSize::MEDIUM); // floor 1 has medium
        lot.addSpot(2, SpotSize::MEDIUM); // floor 2 has medium
        Vehicle car{"CAR3", VehicleType::CAR};
        Ticket* t = lot.parkVehicle(car, 1000);
        assert(t != nullptr);
        assert(t->floor == 1); // floor 0 has no compatible spot, floor 1 is nearest
        cout << "PASS test_nearest_floor_allocation" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_nearest_floor_allocation" << endl;
        failed++;
    }

    // Test 11: Full lot returns nullptr
    try {
        ParkingLot lot(1);
        lot.addSpot(0, SpotSize::MEDIUM);
        Vehicle car1{"C1", VehicleType::CAR};
        Vehicle car2{"C2", VehicleType::CAR};
        assert(lot.parkVehicle(car1, 1000) != nullptr);
        assert(lot.parkVehicle(car2, 1000) == nullptr);
        cout << "PASS test_full_lot" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_full_lot" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
