#include <iostream>
#include <vector>
#include <string>
#include <set>
#include <unordered_map>
#include <algorithm>
#include <climits>
using namespace std;

// --- Data Model -------------------------------------------------------------

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
    int getPendingCount() const { return (int)upRequests.size() + (int)downRequests.size(); }

    void addRequest(int floor, Direction direction) {
        // If floor == currentFloor and elevator is IDLE, open doors immediately
        if (floor == currentFloor && state == ElevatorState::IDLE) {
            state = ElevatorState::DOOR_OPEN;
            return;
        }

        // Route to appropriate set based on requested direction
        if (direction == Direction::UP) {
            upRequests.insert(floor);
        } else if (direction == Direction::DOWN) {
            downRequests.insert(floor);
        } else {
            // Direction::NONE — fallback to position-based routing
            if (floor > currentFloor) upRequests.insert(floor);
            else downRequests.insert(floor);
        }

        // If IDLE, pick direction toward nearest request
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
            case ElevatorState::IDLE:
                break;

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
                    // Direction::NONE (was idle, door opened at current floor)
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
        int bestIdx = 0;
        int bestScore = INT_MAX;
        const int PENALTY = 10000;

        for (int i = 0; i < (int)elevators.size(); i++) {
            int dist = abs(elevators[i]->getCurrentFloor() - requestFloor);
            int score = dist;
            ElevatorState st = elevators[i]->getState();
            Direction dir = elevators[i]->getCurrentDirection();

            if (st == ElevatorState::IDLE || dir == Direction::NONE) {
                // Idle elevator: just use distance, no penalty
                score = dist;
            } else if (dir == Direction::UP && requestDirection == Direction::UP
                       && elevators[i]->getCurrentFloor() <= requestFloor) {
                // Moving up, request is up and ahead: good match
                score = dist;
            } else if (dir == Direction::DOWN && requestDirection == Direction::DOWN
                       && elevators[i]->getCurrentFloor() >= requestFloor) {
                // Moving down, request is down and ahead: good match
                score = dist;
            } else {
                // Wrong direction or request is behind: penalize
                score = dist + PENALTY;
            }

            if (score < bestScore) {
                bestScore = score;
                bestIdx = i;
            }
        }
        return bestIdx;
    }
};

// --- LeastLoaded Strategy ---------------------------------------------------

class LeastLoaded : public DispatchStrategy {
public:
    int selectElevator(const vector<Elevator*>& elevators,
                       int requestFloor,
                       Direction requestDirection) override {
        int bestIdx = 0;
        int bestCount = INT_MAX;
        for (int i = 0; i < (int)elevators.size(); i++) {
            int cnt = elevators[i]->getPendingCount();
            if (cnt < bestCount) {
                bestCount = cnt;
                bestIdx = i;
            }
        }
        return bestIdx;
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
        elevators.push_back(new Elevator(id));
    }

    void setDispatchStrategy(DispatchStrategy* s) {
        strategy = s;
    }

    Elevator* getElevator(int index) const {
        if (index < 0 || index >= (int)elevators.size()) return nullptr;
        return elevators[index];
    }

    int getElevatorCount() const { return (int)elevators.size(); }

    void addRequest(int floor, Direction direction) {
        if (elevators.empty()) return;
        int idx = 0;
        if (strategy) {
            idx = strategy->selectElevator(elevators, floor, direction);
        }
        elevators[idx]->addRequest(floor, direction);
    }

    void step() {
        for (auto* e : elevators) {
            e->step();
        }
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Elevator System -- run tests to verify implementation." << endl;
    return 0;
}
#endif
