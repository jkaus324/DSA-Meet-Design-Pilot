// Part 2 Tests — Ride-Sharing: Pluggable Selection Strategies
// Tests MostVacant and PreferredVehicle ride selection strategies

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

// Helper to set up a standard test scenario
RideService setupTestService() {
    RideService service;
    service.addUser("Rohan");
    service.addUser("Deepa");
    service.addUser("Amit");
    service.addUser("Priya");

    service.addVehicle("Rohan", "Swift", "KA-01-1234");
    service.addVehicle("Deepa", "XUV", "KA-02-5678");
    service.addVehicle("Amit", "Swift", "KA-03-9012");

    // Rohan offers Bangalore→Mysore, 3 seats, Swift
    service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
    // Deepa offers Bangalore→Mysore, 5 seats, XUV
    service.offerRide("Deepa", "Bangalore", "Mysore", 5, "KA-02-5678");
    // Amit offers Bangalore→Chennai, 2 seats, Swift
    service.offerRide("Amit", "Bangalore", "Chennai", 2, "KA-03-9012");

    return service;
}

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: MostVacant selects ride with most available seats
    try {
        RideService service = setupTestService();
        MostVacantStrategy strategy;
        string rideId = service.selectRide("Priya", "Bangalore", "Mysore", 1, &strategy);
        assert(!rideId.empty());
        const Ride& ride = service.getRide(rideId);
        // Deepa's ride has 5 seats (most vacant)
        assert(ride.driverId == "Deepa");
        cout << "PASS test_most_vacant_strategy" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_most_vacant_strategy" << endl;
        failed++;
    }

    // Test 2: PreferredVehicle selects ride with matching model
    try {
        RideService service = setupTestService();
        PreferredVehicleStrategy strategy(service.getVehicles());
        string rideId = service.selectRide("Priya", "Bangalore", "Mysore", 1, &strategy, "Swift");
        assert(!rideId.empty());
        const Ride& ride = service.getRide(rideId);
        // Rohan has a Swift on this route
        assert(ride.driverId == "Rohan");
        cout << "PASS test_preferred_vehicle_strategy" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_preferred_vehicle_strategy" << endl;
        failed++;
    }

    // Test 3: No match for route returns empty
    try {
        RideService service = setupTestService();
        MostVacantStrategy strategy;
        string rideId = service.selectRide("Priya", "Delhi", "Mumbai", 1, &strategy);
        assert(rideId.empty());
        cout << "PASS test_no_matching_route" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_no_matching_route" << endl;
        failed++;
    }

    // Test 4: Selecting ride decrements available seats
    try {
        RideService service = setupTestService();
        MostVacantStrategy strategy;
        string rideId = service.selectRide("Priya", "Bangalore", "Mysore", 2, &strategy);
        assert(!rideId.empty());
        const Ride& ride = service.getRide(rideId);
        // Deepa's ride: 5 total, now 3 available
        assert(ride.availableSeats == 3);
        cout << "PASS test_seats_decremented" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_seats_decremented" << endl;
        failed++;
    }

    // Test 5: Selecting ride increments passenger's ridesTaken
    try {
        RideService service = setupTestService();
        MostVacantStrategy strategy;
        assert(service.getUser("Priya").ridesTaken == 0);
        service.selectRide("Priya", "Bangalore", "Mysore", 1, &strategy);
        assert(service.getUser("Priya").ridesTaken == 1);
        cout << "PASS test_rides_taken_incremented" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_rides_taken_incremented" << endl;
        failed++;
    }

    // Test 6: Cannot select own offered ride
    try {
        RideService service = setupTestService();
        MostVacantStrategy strategy;
        // Deepa tries to select a Bangalore→Mysore ride, but she offered one
        // Only Rohan's ride should be a candidate for Deepa
        string rideId = service.selectRide("Deepa", "Bangalore", "Mysore", 1, &strategy);
        if (!rideId.empty()) {
            const Ride& ride = service.getRide(rideId);
            assert(ride.driverId != "Deepa");  // must not select own ride
        }
        cout << "PASS test_cannot_select_own_ride" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_cannot_select_own_ride" << endl;
        failed++;
    }

    // Test 7: Not enough seats returns empty
    try {
        RideService service = setupTestService();
        MostVacantStrategy strategy;
        // Request 10 seats — no ride has that many
        string rideId = service.selectRide("Priya", "Bangalore", "Mysore", 10, &strategy);
        assert(rideId.empty());
        cout << "PASS test_not_enough_seats" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_not_enough_seats" << endl;
        failed++;
    }

    // Test 8: PreferredVehicle with no matching model returns empty
    try {
        RideService service = setupTestService();
        PreferredVehicleStrategy strategy(service.getVehicles());
        string rideId = service.selectRide("Priya", "Bangalore", "Mysore", 1, &strategy, "BMW");
        assert(rideId.empty());
        cout << "PASS test_no_matching_vehicle_model" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_no_matching_vehicle_model" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
