// Part 1 Tests — Vending Machine
// Tests state transitions: idle -> select -> pay -> dispense

#include "solution.cpp"
#include <cassert>
#include <iostream>
#include <string>
using namespace std;

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: initial state is Idle
    try {
        assert(getState() == "Idle");
        cout << "PASS test_initial_state_idle" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_initial_state_idle" << endl;
        failed++;
    }

    // Test 2: select item transitions to PaymentPending
    try {
        reset(); // reset machine to clean state
        selectItem("Cola");
        string s = getState();
        assert(s == "PaymentPending" || s == "ItemSelected");
        cout << "PASS test_select_transitions_state" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_select_transitions_state" << endl;
        failed++;
    }

    // Test 3: insert enough money and dispense returns to Idle
    try {
        reset();
        selectItem("Cola");
        insertMoney(25.0);  // assume Cola costs 25
        dispense();
        assert(getState() == "Idle");
        cout << "PASS test_full_purchase_cycle" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_full_purchase_cycle" << endl;
        failed++;
    }

    // Test 4: cancel from PaymentPending returns to Idle
    try {
        reset();
        selectItem("Cola");
        insertMoney(10.0);
        cancel();
        assert(getState() == "Idle");
        cout << "PASS test_cancel_returns_idle" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_cancel_returns_idle" << endl;
        failed++;
    }

    // Test 5: pay before select does nothing harmful
    try {
        reset();
        insertMoney(50.0);  // should be ignored or print warning
        assert(getState() == "Idle"); // should stay idle
        cout << "PASS test_pay_before_select_safe" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_pay_before_select_safe" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
