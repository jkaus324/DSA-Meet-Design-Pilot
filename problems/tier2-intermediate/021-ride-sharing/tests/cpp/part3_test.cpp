// Part 3 Tests — Ride-Sharing: End Rides + Statistics
// Tests ride lifecycle and per-user statistics tracking

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part3_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: endRide marks ride as inactive
    try {
        RideService service;
        service.addUser("Rohan");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        string rideId = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
        assert(service.getRide(rideId).active == true);
        service.endRide(rideId);
        assert(service.getRide(rideId).active == false);
        cout << "PASS test_end_ride_marks_inactive" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_end_ride_marks_inactive" << endl;
        failed++;
    }

    // Test 2: After endRide, vehicle can be used for new ride
    try {
        RideService service;
        service.addUser("Rohan");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        string ride1 = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
        assert(!ride1.empty());

        // Cannot offer again while active
        string ride2 = service.offerRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234");
        assert(ride2.empty());

        // End first ride
        service.endRide(ride1);

        // Now can offer again
        string ride3 = service.offerRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234");
        assert(!ride3.empty());
        cout << "PASS test_vehicle_freed_after_end" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_vehicle_freed_after_end" << endl;
        failed++;
    }

    // Test 3: Ending already-ended ride is a no-op
    try {
        RideService service;
        service.addUser("Rohan");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        string rideId = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
        service.endRide(rideId);
        service.endRide(rideId);  // should not crash
        assert(service.getRide(rideId).active == false);
        cout << "PASS test_end_ride_idempotent" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_end_ride_idempotent" << endl;
        failed++;
    }

    // Test 4: Ending nonexistent ride is a no-op
    try {
        RideService service;
        service.endRide("RIDE-999");  // should not crash
        cout << "PASS test_end_nonexistent_ride" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_end_nonexistent_ride" << endl;
        failed++;
    }

    // Test 5: getRideStats returns correct counts
    try {
        RideService service;
        service.addUser("Rohan");
        service.addUser("Deepa");
        service.addUser("Priya");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        service.addVehicle("Deepa", "XUV", "KA-02-5678");

        service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
        service.offerRide("Deepa", "Bangalore", "Mysore", 5, "KA-02-5678");

        MostVacantStrategy strategy;
        service.selectRide("Priya", "Bangalore", "Mysore", 1, &strategy);

        auto stats = service.getRideStats();
        assert(!stats.empty());

        // Find each user's stats
        bool foundRohan = false, foundDeepa = false, foundPriya = false;
        for (auto& [name, counts] : stats) {
            if (name == "Rohan") {
                assert(counts.first == 1);   // offered 1
                assert(counts.second == 0);  // taken 0
                foundRohan = true;
            }
            if (name == "Deepa") {
                assert(counts.first == 1);   // offered 1
                assert(counts.second == 0);  // taken 0
                foundDeepa = true;
            }
            if (name == "Priya") {
                assert(counts.first == 0);   // offered 0
                assert(counts.second == 1);  // taken 1
                foundPriya = true;
            }
        }
        assert(foundRohan && foundDeepa && foundPriya);
        cout << "PASS test_ride_stats" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_ride_stats" << endl;
        failed++;
    }

    // Test 6: Ended ride not selectable by future passengers
    try {
        RideService service;
        service.addUser("Rohan");
        service.addUser("Priya");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");
        string rideId = service.offerRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234");
        service.endRide(rideId);

        MostVacantStrategy strategy;
        string selected = service.selectRide("Priya", "Bangalore", "Mysore", 1, &strategy);
        assert(selected.empty());  // ended ride should not be selectable
        cout << "PASS test_ended_ride_not_selectable" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_ended_ride_not_selectable" << endl;
        failed++;
    }

    // Test 7: Full workflow — offer, select, end, re-offer
    try {
        RideService service;
        service.addUser("Rohan");
        service.addUser("Deepa");
        service.addVehicle("Rohan", "Swift", "KA-01-1234");

        // Rohan offers ride
        string ride1 = service.offerRide("Rohan", "Bangalore", "Mysore", 2, "KA-01-1234");
        assert(!ride1.empty());

        // Deepa takes ride
        MostVacantStrategy strategy;
        string selected = service.selectRide("Deepa", "Bangalore", "Mysore", 1, &strategy);
        assert(selected == ride1);
        assert(service.getRide(ride1).availableSeats == 1);

        // Rohan ends ride
        service.endRide(ride1);
        assert(service.getRide(ride1).active == false);

        // Rohan offers new ride with same vehicle
        string ride2 = service.offerRide("Rohan", "Mysore", "Bangalore", 2, "KA-01-1234");
        assert(!ride2.empty());
        assert(ride2 != ride1);

        // Check stats
        assert(service.getUser("Rohan").ridesOffered == 2);
        assert(service.getUser("Deepa").ridesTaken == 1);
        cout << "PASS test_full_workflow" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_full_workflow" << endl;
        failed++;
    }

    cout << "PART3_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
