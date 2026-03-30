// Part 1 Tests — Ride-Sharing: User, Vehicle, and Ride Onboarding
// Tests user registration, vehicle registration, and ride offering

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Add user and verify exists
    try {
        RideService service;
        service.addUser("Rohan");
        assert(service.hasUser("Rohan"));
        assert(!service.hasUser("Unknown"));
        cout << "PASS test_add_user" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_add_user" << endl;
        failed++;
    }

    // Test 2: Add vehicle and verify exists
    try {
        RideService service;
        service.addUser("Rohan");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        assert(service.hasVehicle("KA-01-1234"));
        cout << "PASS test_add_vehicle" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_add_vehicle" << endl;
        failed++;
    }

    // Test 3: Offer ride returns valid rideId
    try {
        RideService service;
        service.addUser("Rohan");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        string rideId = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
        assert(!rideId.empty());
        assert(service.hasRide(rideId));
        const Ride& ride = service.getRide(rideId);
        assert(ride.origin == "Bangalore");
        assert(ride.destination == "Mysore");
        assert(ride.totalSeats == 3);
        assert(ride.availableSeats == 3);
        assert(ride.active == true);
        cout << "PASS test_offer_ride" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_offer_ride" << endl;
        failed++;
    }

    // Test 4: Cannot offer ride with vehicle already in active ride
    try {
        RideService service;
        service.addUser("Rohan");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        string ride1 = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
        assert(!ride1.empty());
        string ride2 = service.offerRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234");
        assert(ride2.empty());  // should fail — vehicle already active
        cout << "PASS test_no_duplicate_active_ride_per_vehicle" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_no_duplicate_active_ride_per_vehicle" << endl;
        failed++;
    }

    // Test 5: Cannot offer ride with someone else's vehicle
    try {
        RideService service;
        service.addUser("Rohan");
        service.addUser("Deepa");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        string rideId = service.offerRide("Deepa", "Bangalore", "Mysore", 2, "KA-01-1234");
        assert(rideId.empty());  // Deepa doesn't own KA-01-1234
        cout << "PASS test_cannot_use_others_vehicle" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_cannot_use_others_vehicle" << endl;
        failed++;
    }

    // Test 6: Offering ride increments ridesOffered
    try {
        RideService service;
        service.addUser("Rohan");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        service.addVehicle("Rohan", "XUV", "KA-01-5678");
        assert(service.getUser("Rohan").ridesOffered == 0);
        service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
        assert(service.getUser("Rohan").ridesOffered == 1);
        service.offerRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-5678");
        assert(service.getUser("Rohan").ridesOffered == 2);
        cout << "PASS test_rides_offered_counter" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_rides_offered_counter" << endl;
        failed++;
    }

    // Test 7: Offer ride with nonexistent user returns empty
    try {
        RideService service;
        service.addUser("Rohan");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        string rideId = service.offerRide("Ghost", "A", "B", 2, "KA-01-1234");
        assert(rideId.empty());
        cout << "PASS test_offer_ride_invalid_user" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_offer_ride_invalid_user" << endl;
        failed++;
    }

    // Test 8: Offer ride with nonexistent vehicle returns empty
    try {
        RideService service;
        service.addUser("Rohan");
        string rideId = service.offerRide("Rohan", "A", "B", 2, "INVALID-REG");
        assert(rideId.empty());
        cout << "PASS test_offer_ride_invalid_vehicle" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_offer_ride_invalid_vehicle" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
