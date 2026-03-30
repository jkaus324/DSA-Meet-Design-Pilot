// Part 2 Tests -- Multiple Elevators with Dispatch Strategies
// Tests elevator management, NearestFirst, and LeastLoaded dispatch

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Add elevators and verify count
    try {
        ElevatorSystem sys;
        sys.addElevator(1);
        sys.addElevator(2);
        assert(sys.getElevatorCount() == 2);
        assert(sys.getElevator(0) != nullptr);
        assert(sys.getElevator(1) != nullptr);
        cout << "PASS test_add_elevators" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_add_elevators" << endl;
        failed++;
    }

    // Test 2: Without strategy, request goes to first elevator
    try {
        ElevatorSystem sys;
        sys.addElevator(1);
        sys.addElevator(2);
        sys.addRequest(5, Direction::UP);
        assert(sys.getElevator(0)->getState() != ElevatorState::IDLE);
        assert(sys.getElevator(1)->getState() == ElevatorState::IDLE);
        cout << "PASS test_default_dispatch" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_default_dispatch" << endl;
        failed++;
    }

    // Test 3: NearestFirst dispatches to closer elevator
    try {
        ElevatorSystem sys;
        NearestFirst nf;
        sys.addElevator(1);
        sys.addElevator(2);
        sys.setDispatchStrategy(&nf);
        // Move elevator 0 to floor 5
        sys.getElevator(0)->addRequest(5, Direction::UP);
        for (int i = 0; i < 5; i++) sys.getElevator(0)->step();
        // Elevator 0 is at floor 5 (DOOR_OPEN), elevator 1 is at floor 0
        sys.getElevator(0)->step(); // close doors, go idle at 5
        // Request floor 2 — elevator 1 (at 0) is closer than elevator 0 (at 5)
        sys.addRequest(2, Direction::UP);
        assert(sys.getElevator(1)->getState() != ElevatorState::IDLE);
        cout << "PASS test_nearest_first_dispatch" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_nearest_first_dispatch" << endl;
        failed++;
    }

    // Test 4: LeastLoaded dispatches to elevator with fewer requests
    try {
        ElevatorSystem sys;
        LeastLoaded ll;
        sys.addElevator(1);
        sys.addElevator(2);
        sys.setDispatchStrategy(&ll);
        // Give elevator 0 three requests
        sys.getElevator(0)->addRequest(3, Direction::UP);
        sys.getElevator(0)->addRequest(5, Direction::UP);
        sys.getElevator(0)->addRequest(7, Direction::UP);
        // Elevator 0 has 3 pending, elevator 1 has 0
        sys.addRequest(4, Direction::UP);
        // Should go to elevator 1 (least loaded)
        assert(sys.getElevator(1)->getPendingCount() > 0);
        cout << "PASS test_least_loaded_dispatch" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_least_loaded_dispatch" << endl;
        failed++;
    }

    // Test 5: step() advances all elevators
    try {
        ElevatorSystem sys;
        sys.addElevator(1);
        sys.addElevator(2);
        sys.getElevator(0)->addRequest(3, Direction::UP);
        sys.getElevator(1)->addRequest(2, Direction::UP);
        sys.step(); // both move one floor
        assert(sys.getElevator(0)->getCurrentFloor() == 1);
        assert(sys.getElevator(1)->getCurrentFloor() == 1);
        cout << "PASS test_step_all_elevators" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_step_all_elevators" << endl;
        failed++;
    }

    // Test 6: NearestFirst prefers same-direction elevator
    try {
        ElevatorSystem sys;
        NearestFirst nf;
        sys.addElevator(1);
        sys.addElevator(2);
        sys.setDispatchStrategy(&nf);
        // Move elevator 0 to floor 3 going up (give it request for floor 8)
        sys.getElevator(0)->addRequest(3, Direction::UP);
        sys.getElevator(0)->addRequest(8, Direction::UP);
        for (int i = 0; i < 3; i++) sys.getElevator(0)->step();
        // Elevator 0 at floor 3, DOOR_OPEN, still has request for 8 (moving up)
        sys.getElevator(0)->step(); // close doors, MOVING_UP toward 8
        // Elevator 1 at floor 0
        // Request floor 5 UP — elevator 0 is moving up past it, should be preferred
        sys.addRequest(5, Direction::UP);
        // Elevator 0 should get it (moving up, will pass floor 5)
        assert(sys.getElevator(0)->getPendingCount() >= 2); // has floor 8 + floor 5
        cout << "PASS test_nearest_prefers_same_direction" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_nearest_prefers_same_direction" << endl;
        failed++;
    }

    // Test 7: Swapping strategy at runtime
    try {
        ElevatorSystem sys;
        NearestFirst nf;
        LeastLoaded ll;
        sys.addElevator(1);
        sys.addElevator(2);
        sys.setDispatchStrategy(&nf);
        sys.addRequest(3, Direction::UP); // dispatched via NearestFirst
        sys.setDispatchStrategy(&ll);
        // Now give elevator 0 more requests so elevator 1 is least loaded
        sys.getElevator(0)->addRequest(5, Direction::UP);
        sys.getElevator(0)->addRequest(7, Direction::UP);
        sys.addRequest(4, Direction::UP); // dispatched via LeastLoaded to elevator 1
        assert(sys.getElevator(1)->getPendingCount() > 0);
        cout << "PASS test_swap_strategy_runtime" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_swap_strategy_runtime" << endl;
        failed++;
    }

    // Test 8: getElevator with invalid index returns nullptr
    try {
        ElevatorSystem sys;
        sys.addElevator(1);
        assert(sys.getElevator(-1) == nullptr);
        assert(sys.getElevator(5) == nullptr);
        cout << "PASS test_invalid_elevator_index" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_invalid_elevator_index" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
