#include <iostream>
#include <vector>
#include <string>
#include <set>
#include <unordered_map>
#include <algorithm>
#include <climits>
using namespace std;

// --- Data Model (given -- do not modify) ------------------------------------

enum class ElevatorState {
    IDLE,
    MOVING_UP,
    MOVING_DOWN,
    DOOR_OPEN
};

enum class Direction {
    UP,
    DOWN,
    NONE
};

struct Request {
    int floor;
    Direction direction;
};

// --- Elevator ---------------------------------------------------------------

class Elevator {
    int currentFloor;
    ElevatorState state;
    Direction currentDirection;
    set<int> upRequests;    // floors to visit going up (sorted ascending)
    set<int> downRequests;  // floors to visit going down (sorted ascending, iterate in reverse)

public:
    Elevator() : currentFloor(0), state(ElevatorState::IDLE),
                 currentDirection(Direction::NONE) {}

    int getCurrentFloor() const { return currentFloor; }
    ElevatorState getState() const { return state; }
    Direction getCurrentDirection() const { return currentDirection; }

    int getPendingCount() const {
        return upRequests.size() + downRequests.size();
    }

    void addRequest(int floor, Direction direction) {
        // TODO: If floor == currentFloor and elevator is IDLE, transition to DOOR_OPEN
        // TODO: If floor > currentFloor, add to upRequests; otherwise add to downRequests
        // TODO: If elevator is IDLE and requests now exist, pick a direction and
        //       transition to MOVING_UP or MOVING_DOWN based on nearest request
    }

    void step() {
        switch (state) {
            case ElevatorState::IDLE:
                // TODO: Nothing to do when idle
                break;

            case ElevatorState::MOVING_UP:
                // TODO: Increment currentFloor by 1
                // TODO: If currentFloor is in upRequests, erase it and go to DOOR_OPEN
                break;

            case ElevatorState::MOVING_DOWN:
                // TODO: Decrement currentFloor by 1
                // TODO: If currentFloor is in downRequests, erase it and go to DOOR_OPEN
                break;

            case ElevatorState::DOOR_OPEN:
                // TODO: Close doors and decide next state:
                //   If moving UP: continue UP if upRequests not empty,
                //                 else switch to DOWN if downRequests not empty,
                //                 else go IDLE
                //   If moving DOWN: continue DOWN if downRequests not empty,
                //                   else switch to UP if upRequests not empty,
                //                   else go IDLE
                break;
        }
    }
};

// --- Test Entry Points ------------------------------------------------------

Elevator elevator;

void addRequest(int floor, Direction direction) {
    elevator.addRequest(floor, direction);
}

void step() {
    elevator.step();
}

int getCurrentFloor() {
    return elevator.getCurrentFloor();
}

ElevatorState getState() {
    return elevator.getState();
}

#ifndef RUNNING_TESTS
int main() {
    cout << "Elevator System -- implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
