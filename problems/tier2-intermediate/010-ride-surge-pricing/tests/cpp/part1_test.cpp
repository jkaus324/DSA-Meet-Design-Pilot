// Part 1 Tests — Ride Surge Pricing Engine
// Tests surge calculation and fare computation

#include "solution.cpp"
#include <cassert>
#include <iostream>
#include <cmath>
using namespace std;

bool approxEqual(double a, double b, double eps = 0.01) {
    return fabs(a - b) < eps;
}

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: no surge when supply > demand
    try {
        PricingContext ctx = {10.0, 20, 5, "morning", "clear"};
        double surge = calculateSurge(ctx);
        assert(surge >= 1.0); // always at least 1.0x
        assert(surge <= 1.5); // low demand shouldn't spike
        cout << "PASS test_no_surge_normal_conditions" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_no_surge_normal_conditions" << endl;
        failed++;
    }

    // Test 2: high demand ratio causes surge
    try {
        PricingContext ctx = {10.0, 2, 10, "evening", "clear"}; // 5:1 demand ratio
        double surge = calculateSurge(ctx);
        assert(surge > 1.0); // must surge
        cout << "PASS test_high_demand_causes_surge" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_high_demand_causes_surge" << endl;
        failed++;
    }

    // Test 3: storm weather increases surge
    try {
        PricingContext ctx = {10.0, 10, 10, "morning", "storm"};
        double surge = calculateSurge(ctx);
        assert(surge > 1.0);
        cout << "PASS test_storm_increases_surge" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_storm_increases_surge" << endl;
        failed++;
    }

    // Test 4: fare = baseFare * surgeMultiplier
    try {
        PricingContext ctx = {100.0, 5, 5, "morning", "clear"};
        RideRequest req = {"u1", "A", "B", "economy"};
        double fare = calculateFare(req, ctx);
        double surge = calculateSurge(ctx);
        assert(approxEqual(fare, 100.0 * surge));
        cout << "PASS test_fare_calculation" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_fare_calculation" << endl;
        failed++;
    }

    // Test 5: surge is never below 1.0
    try {
        PricingContext ctx = {50.0, 100, 1, "morning", "clear"}; // plenty of drivers
        double surge = calculateSurge(ctx);
        assert(surge >= 1.0);
        cout << "PASS test_surge_minimum_one" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_surge_minimum_one" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
