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

// --- Elevator (from Part 1 -- assume fully implemented) ---------------------

class Elevator {
    int id;
    int currentFloor;
    ElevatorState state;
    Direction currentDirection;
    set<int> upRequests;
    set<int> downRequests;

public:
    Elevator(int id = 0) : id(id), currentFloor(0), state(ElevatorState::IDLE),
                           currentDirection(Direction::NONE) {}

    int getId() const { return id; }
    int getCurrentFloor() const { return currentFloor; }
    ElevatorState getState() const { return state; }
    Direction getCurrentDirection() const { return currentDirection; }
    int getPendingCount() const { return upRequests.size() + downRequests.size(); }

    void addRequest(int floor, Direction direction) {
        if (floor == currentFloor && state == ElevatorState::IDLE) {
            state = ElevatorState::DOOR_OPEN;
            return;
        }
        if (floor > currentFloor || (floor == currentFloor && direction == Direction::UP)) {
            upRequests.insert(floor);
        } else {
            downRequests.insert(floor);
        }
        if (state == ElevatorState::IDLE) {
            if (!upRequests.empty() && (downRequests.empty() ||
                abs(*upRequests.begin() - currentFloor) <= abs(*downRequests.rbegin() - currentFloor))) {
                currentDirection = Direction::UP;
                state = ElevatorState::MOVING_UP;
            } else {
                currentDirection = Direction::DOWN;
                state = ElevatorState::MOVING_DOWN;
            }
        }
    }

    void step() {
        switch (state) {
            case ElevatorState::IDLE: break;
            case ElevatorState::MOVING_UP:
                currentFloor++;
                if (upRequests.count(currentFloor)) {
                    upRequests.erase(currentFloor);
                    state = ElevatorState::DOOR_OPEN;
                }
                break;
            case ElevatorState::MOVING_DOWN:
                currentFloor--;
                if (downRequests.count(currentFloor)) {
                    downRequests.erase(currentFloor);
                    state = ElevatorState::DOOR_OPEN;
                }
                break;
            case ElevatorState::DOOR_OPEN:
                if (currentDirection == Direction::UP) {
                    if (!upRequests.empty()) state = ElevatorState::MOVING_UP;
                    else if (!downRequests.empty()) { currentDirection = Direction::DOWN; state = ElevatorState::MOVING_DOWN; }
                    else { currentDirection = Direction::NONE; state = ElevatorState::IDLE; }
                } else if (currentDirection == Direction::DOWN) {
                    if (!downRequests.empty()) state = ElevatorState::MOVING_DOWN;
                    else if (!upRequests.empty()) { currentDirection = Direction::UP; state = ElevatorState::MOVING_UP; }
                    else { currentDirection = Direction::NONE; state = ElevatorState::IDLE; }
                } else {
                    if (!upRequests.empty()) { currentDirection = Direction::UP; state = ElevatorState::MOVING_UP; }
                    else if (!downRequests.empty()) { currentDirection = Direction::DOWN; state = ElevatorState::MOVING_DOWN; }
                    else state = ElevatorState::IDLE;
                }
                break;
        }
    }
};

// --- Dispatch Strategy Interface --------------------------------------------

class DispatchStrategy {
public:
    virtual int selectElevator(const vector<Elevator*>& elevators,
                               int requestFloor,
                               Direction requestDirection) = 0;
    virtual ~DispatchStrategy() = default;
};

// --- NearestFirst Strategy --------------------------------------------------

class NearestFirst : public DispatchStrategy {
public:
    int selectElevator(const vector<Elevator*>& elevators,
                       int requestFloor,
                       Direction requestDirection) override {
        // TODO: For each elevator, compute a score based on distance
        // TODO: Penalize elevators moving in the wrong direction (add large offset)
        // TODO: Return the index of the elevator with the lowest score
        return 0;
    }
};

// --- LeastLoaded Strategy ---------------------------------------------------

class LeastLoaded : public DispatchStrategy {
public:
    int selectElevator(const vector<Elevator*>& elevators,
                       int requestFloor,
                       Direction requestDirection) override {
        // TODO: Find the elevator with the smallest getPendingCount()
        // TODO: Return its index
        return 0;
    }
};

// --- Elevator System --------------------------------------------------------

class ElevatorSystem {
    vector<Elevator*> elevators;
    DispatchStrategy* strategy;

public:
    ElevatorSystem() : strategy(nullptr) {}

    ~ElevatorSystem() {
        for (auto* e : elevators) delete e;
    }

    void addElevator(int id) {
        // TODO: Create a new Elevator with the given id and add to the vector
    }

    void setDispatchStrategy(DispatchStrategy* s) {
        // TODO: Store the strategy pointer
    }

    Elevator* getElevator(int index) const {
        if (index < 0 || index >= (int)elevators.size()) return nullptr;
        return elevators[index];
    }

    int getElevatorCount() const { return elevators.size(); }

    void addRequest(int floor, Direction direction) {
        // TODO: If no elevators, return
        // TODO: Use the strategy to select the best elevator index
        // TODO: Call addRequest on the selected elevator
    }

    void step() {
        // TODO: Call step() on every elevator
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Elevator System Part 2 -- implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
