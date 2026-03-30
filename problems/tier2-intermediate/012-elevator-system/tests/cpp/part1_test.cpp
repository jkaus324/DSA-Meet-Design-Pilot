// Part 1 Tests -- Single Elevator with SCAN Ordering
// Tests state transitions, SCAN direction handling, and request processing

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Initial state — elevator starts at floor 0, IDLE
    try {
        Elevator e;
        assert(e.getCurrentFloor() == 0);
        assert(e.getState() == ElevatorState::IDLE);
        cout << "PASS test_initial_state" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_initial_state" << endl;
        failed++;
    }

    // Test 2: Add upward request — elevator starts moving up
    try {
        Elevator e;
        e.addRequest(3, Direction::UP);
        assert(e.getState() == ElevatorState::MOVING_UP);
        cout << "PASS test_add_upward_request" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_add_upward_request" << endl;
        failed++;
    }

    // Test 3: Step moves one floor at a time
    try {
        Elevator e;
        e.addRequest(3, Direction::UP);
        e.step(); // floor 0 -> 1
        assert(e.getCurrentFloor() == 1);
        e.step(); // floor 1 -> 2
        assert(e.getCurrentFloor() == 2);
        cout << "PASS test_step_moves_one_floor" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_step_moves_one_floor" << endl;
        failed++;
    }

    // Test 4: Elevator opens doors at requested floor
    try {
        Elevator e;
        e.addRequest(2, Direction::UP);
        e.step(); // floor 0 -> 1
        e.step(); // floor 1 -> 2, door opens
        assert(e.getCurrentFloor() == 2);
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        cout << "PASS test_door_opens_at_target" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_door_opens_at_target" << endl;
        failed++;
    }

    // Test 5: After door open with no more requests, goes IDLE
    try {
        Elevator e;
        e.addRequest(1, Direction::UP);
        e.step(); // floor 0 -> 1, door opens
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        e.step(); // close doors, go idle
        assert(e.getState() == ElevatorState::IDLE);
        assert(e.getCurrentFloor() == 1);
        cout << "PASS test_idle_after_last_request" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_idle_after_last_request" << endl;
        failed++;
    }

    // Test 6: SCAN order — serve all upward requests before reversing
    try {
        Elevator e;
        e.addRequest(5, Direction::UP);
        e.addRequest(2, Direction::UP);
        // Should go up: 0->1->2(stop)->3->4->5(stop)
        e.step(); // 1
        e.step(); // 2, DOOR_OPEN
        assert(e.getCurrentFloor() == 2);
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        e.step(); // close doors, resume MOVING_UP
        assert(e.getState() == ElevatorState::MOVING_UP);
        e.step(); // 3
        e.step(); // 4
        e.step(); // 5, DOOR_OPEN
        assert(e.getCurrentFloor() == 5);
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        cout << "PASS test_scan_upward_order" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_scan_upward_order" << endl;
        failed++;
    }

    // Test 7: SCAN reversal — after upward done, serve downward requests
    try {
        Elevator e;
        e.addRequest(3, Direction::UP);
        e.addRequest(1, Direction::DOWN); // below current, goes into downRequests
        // At floor 0: upRequests={3}, downRequests={} — wait, 1 > 0 so it goes to upRequests
        // Let's use a scenario where we add a downward request after moving
        e.addRequest(3, Direction::UP);
        // Step to floor 3
        e.step(); // 1
        e.step(); // 2
        e.step(); // 3, DOOR_OPEN (first request for 3 served)
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        // Now add a downward request
        e.addRequest(1, Direction::DOWN); // 1 < 3, goes to downRequests
        e.step(); // close doors, should switch to MOVING_DOWN
        assert(e.getState() == ElevatorState::MOVING_DOWN);
        e.step(); // 2
        e.step(); // 1, DOOR_OPEN
        assert(e.getCurrentFloor() == 1);
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        cout << "PASS test_scan_reversal" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_scan_reversal" << endl;
        failed++;
    }

    // Test 8: Request at current floor while IDLE opens doors immediately
    try {
        Elevator e;
        e.addRequest(0, Direction::UP);
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        assert(e.getCurrentFloor() == 0);
        cout << "PASS test_request_at_current_floor" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_request_at_current_floor" << endl;
        failed++;
    }

    // Test 9: step() on IDLE elevator does nothing
    try {
        Elevator e;
        e.step();
        assert(e.getCurrentFloor() == 0);
        assert(e.getState() == ElevatorState::IDLE);
        e.step();
        assert(e.getCurrentFloor() == 0);
        cout << "PASS test_idle_step_noop" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_idle_step_noop" << endl;
        failed++;
    }

    // Test 10: Multiple requests in same direction served in order
    try {
        Elevator e;
        e.addRequest(5, Direction::UP);
        e.addRequest(3, Direction::UP);
        e.addRequest(1, Direction::UP);
        // Should stop at 1, 3, 5
        e.step(); // floor 1, DOOR_OPEN
        assert(e.getCurrentFloor() == 1);
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        e.step(); // close, MOVING_UP
        e.step(); // floor 2
        e.step(); // floor 3, DOOR_OPEN
        assert(e.getCurrentFloor() == 3);
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        e.step(); // close, MOVING_UP
        e.step(); // floor 4
        e.step(); // floor 5, DOOR_OPEN
        assert(e.getCurrentFloor() == 5);
        assert(e.getState() == ElevatorState::DOOR_OPEN);
        cout << "PASS test_multiple_stops_in_order" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_multiple_stops_in_order" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
