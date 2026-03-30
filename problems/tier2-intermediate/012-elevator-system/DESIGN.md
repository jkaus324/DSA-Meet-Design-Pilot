# Design Walkthrough — Elevator System

> This file is the answer guide. Only read after you've attempted the problem.

---

## The Core Design Decision

Three concerns need clean separation:

1. **What does the elevator do in each state?** — The elevator's behavior changes depending on whether it's idle, moving, or has doors open. This is the State pattern.
2. **In what order are floors served?** — The SCAN algorithm (serve all in one direction, then reverse) requires sorted data structures per direction.
3. **Which elevator gets the request?** — With multiple elevators, the dispatch logic varies (nearest, least loaded). This is the Strategy pattern.

```
ElevatorSystem
    ├── DispatchStrategy* (swappable — NearestFirst / LeastLoaded)
    └── vector<Elevator*> (all managed elevators)

Elevator
    ├── ElevatorState (IDLE / MOVING_UP / MOVING_DOWN / DOOR_OPEN)
    ├── set<int> upRequests    (floors to serve going up, ascending order)
    └── set<int> downRequests  (floors to serve going down, descending order)

addRequest(floor, direction):
    strategy->selectElevator(elevators, floor, direction)
    → selected_elevator->addRequest(floor, direction)

step():
    for each elevator:
        if IDLE && has requests → pick direction, transition to MOVING
        if MOVING → move one floor; if floor has request → DOOR_OPEN
        if DOOR_OPEN → close doors, continue or go IDLE
```

---

## Reference Implementation

```cpp
#include <vector>
#include <string>
#include <set>
#include <unordered_map>
#include <algorithm>
#include <climits>
#include <iostream>
using namespace std;

// --- Data Structures ---

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

// --- Single Elevator ---

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

    int getPendingCount() const {
        return upRequests.size() + downRequests.size();
    }

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
        // If idle, start moving
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
                // Nothing to do
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
                // Close doors and decide next action
                if (currentDirection == Direction::UP) {
                    if (!upRequests.empty()) {
                        state = ElevatorState::MOVING_UP;
                    } else if (!downRequests.empty()) {
                        currentDirection = Direction::DOWN;
                        state = ElevatorState::MOVING_DOWN;
                    } else {
                        currentDirection = Direction::NONE;
                        state = ElevatorState::IDLE;
                    }
                } else if (currentDirection == Direction::DOWN) {
                    if (!downRequests.empty()) {
                        state = ElevatorState::MOVING_DOWN;
                    } else if (!upRequests.empty()) {
                        currentDirection = Direction::UP;
                        state = ElevatorState::MOVING_UP;
                    } else {
                        currentDirection = Direction::NONE;
                        state = ElevatorState::IDLE;
                    }
                } else {
                    // Was NONE direction (requested at current floor)
                    if (!upRequests.empty()) {
                        currentDirection = Direction::UP;
                        state = ElevatorState::MOVING_UP;
                    } else if (!downRequests.empty()) {
                        currentDirection = Direction::DOWN;
                        state = ElevatorState::MOVING_DOWN;
                    } else {
                        state = ElevatorState::IDLE;
                    }
                }
                break;
        }
    }
};

// --- Dispatch Strategy ---

class DispatchStrategy {
public:
    virtual int selectElevator(const vector<Elevator*>& elevators,
                               int requestFloor,
                               Direction requestDirection) = 0;
    virtual ~DispatchStrategy() = default;
};

class NearestFirst : public DispatchStrategy {
public:
    int selectElevator(const vector<Elevator*>& elevators,
                       int requestFloor,
                       Direction requestDirection) override {
        int bestIdx = 0;
        int bestScore = INT_MAX;
        for (int i = 0; i < (int)elevators.size(); i++) {
            auto* e = elevators[i];
            int dist = abs(e->getCurrentFloor() - requestFloor);
            int score = dist;
            // Prefer idle or same-direction elevators
            if (e->getState() == ElevatorState::IDLE) {
                score = dist;
            } else if ((e->getCurrentDirection() == Direction::UP && requestDirection == Direction::UP
                         && e->getCurrentFloor() <= requestFloor) ||
                        (e->getCurrentDirection() == Direction::DOWN && requestDirection == Direction::DOWN
                         && e->getCurrentFloor() >= requestFloor)) {
                score = dist;
            } else {
                score = dist + 1000; // penalize wrong direction
            }
            if (score < bestScore) {
                bestScore = score;
                bestIdx = i;
            }
        }
        return bestIdx;
    }
};

class LeastLoaded : public DispatchStrategy {
public:
    int selectElevator(const vector<Elevator*>& elevators,
                       int requestFloor,
                       Direction requestDirection) override {
        int bestIdx = 0;
        int minLoad = INT_MAX;
        for (int i = 0; i < (int)elevators.size(); i++) {
            int load = elevators[i]->getPendingCount();
            if (load < minLoad) {
                minLoad = load;
                bestIdx = i;
            }
        }
        return bestIdx;
    }
};

// --- Elevator System ---

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

    int getElevatorCount() const { return elevators.size(); }

    void addRequest(int floor, Direction direction) {
        if (elevators.empty()) return;
        int idx = 0;
        if (strategy && elevators.size() > 1) {
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
```

---

## What interviewers look for

1. **State transitions**: Can you model the four elevator states cleanly? Do transitions handle edge cases (no more requests, direction reversal)?
2. **SCAN algorithm**: Did you use sorted sets to efficiently track pending floors per direction? The naive approach of a single queue leads to inefficient floor-hopping.
3. **Strategy pattern**: Is the dispatch logic decoupled from the elevator system? Can you add a new strategy (e.g., ZoneBased) without touching the system class?
4. **Data structure choice**: `set<int>` for O(log n) insertion and ordered traversal. HashMap for elevator lookup by ID.

---

## Common interview follow-ups

- *"How would you handle express elevators that skip certain floors?"* -- Add a `set<int> servableFloors` to each Elevator and filter during dispatch.
- *"What if an elevator breaks down mid-operation?"* -- Add an OUT_OF_SERVICE state. Dispatch strategies should skip unavailable elevators.
- *"How do you prioritize VIP floors?"* -- Use a weighted priority queue instead of a plain set. VIP floors get higher priority within the same direction.
- *"How would you optimize for peak hours (morning rush)?"* -- A strategy that sends idle elevators to the lobby floor preemptively.
