// Part 2 Tests — Maintenance Mode
// Tests operator maintenance mode: enter, exit, restock

#include <cassert>
#include <iostream>
#include <string>
using namespace std;

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: valid PIN enters maintenance mode
    try {
        reset();
        enterMaintenance("1234");
        assert(getState() == "Maintenance");
        cout << "PASS test_enter_maintenance_valid_pin" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_enter_maintenance_valid_pin" << endl;
        failed++;
    }

    // Test 2: invalid PIN does not enter maintenance
    try {
        reset();
        enterMaintenance("wrong");
        assert(getState() == "Idle"); // should stay idle
        cout << "PASS test_enter_maintenance_invalid_pin" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_enter_maintenance_invalid_pin" << endl;
        failed++;
    }

    // Test 3: in maintenance, user actions are blocked
    try {
        reset();
        enterMaintenance("1234");
        selectItem("Cola"); // should print warning, not crash
        assert(getState() == "Maintenance"); // stays in maintenance
        cout << "PASS test_user_blocked_in_maintenance" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_user_blocked_in_maintenance" << endl;
        failed++;
    }

    // Test 4: exit maintenance returns to Idle
    try {
        reset();
        enterMaintenance("1234");
        exitMaintenance("1234");
        assert(getState() == "Idle");
        cout << "PASS test_exit_maintenance" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_exit_maintenance" << endl;
        failed++;
    }

    // Test 5: restock only works in maintenance
    try {
        reset();
        restock("Cola", 10); // should be blocked outside maintenance
        enterMaintenance("1234");
        restock("Cola", 5); // should work in maintenance
        cout << "PASS test_restock_in_maintenance" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_restock_in_maintenance" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
