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

// --- Include your Part 1 Elevator class here --------------------------------
// HINT: Your Elevator class needs getId(), getCurrentFloor(), getState(),
//       getCurrentDirection(), and getPendingCount() as public getters
//       so that dispatch strategies can inspect elevator state.

// --- Dispatch Strategy Interface --------------------------------------------
// HINT: The strategy receives a list of all elevators plus the request
//       details, and returns the index of the best elevator.

// class DispatchStrategy {
// public:
//     virtual int selectElevator(const vector<Elevator*>& elevators,
//                                int requestFloor,
//                                Direction requestDirection) = 0;
//     virtual ~DispatchStrategy() = default;
// };

// --- NearestFirst Strategy --------------------------------------------------
// HINT: Calculate distance = abs(elevator.floor - requestFloor).
//       Penalize elevators moving in the wrong direction (add a large offset).
//       Prefer idle elevators or those moving toward the request.

// --- LeastLoaded Strategy ---------------------------------------------------
// HINT: Simply pick the elevator with the smallest getPendingCount().

// --- Elevator System --------------------------------------------------------
// HINT: Holds a vector<Elevator*> and a DispatchStrategy*.
//       addRequest() uses the strategy to pick an elevator, then calls
//       that elevator's addRequest(). step() calls step() on all elevators.

// --- Test Entry Points (must exist for tests to compile) --------------------
// Your solution must provide these via an ElevatorSystem instance:
//
//   void addElevator(int id);
//   void setDispatchStrategy(DispatchStrategy* strategy);
//   void addRequest(int floor, Direction direction);
//   void step();
//   Elevator* getElevator(int index);
//   int getElevatorCount();
//
// -------------------------------------------------------------------------


