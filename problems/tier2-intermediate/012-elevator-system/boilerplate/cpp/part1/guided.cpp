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
// HINT: Use two set<int> to track pending floors — one for upward stops,
//       one for downward stops. A set keeps floors sorted automatically.
// HINT: When the elevator is moving up, check if currentFloor is in
//       upRequests. When moving down, check downRequests.
// HINT: In DOOR_OPEN, check if there are more requests in the current
//       direction. If not, check the other direction. If none, go IDLE.

// class Elevator {
// private:
//     int currentFloor;
//     ElevatorState state;
//     Direction currentDirection;
//     set<int> upRequests;    // floors to visit going up
//     set<int> downRequests;  // floors to visit going down
// public:
//     Elevator();
//     void addRequest(int floor, Direction direction);
//     void step();
//     int getCurrentFloor() const;
//     ElevatorState getState() const;
// };

// --- Test Entry Points (must exist for tests to compile) --------------------
// Your solution must provide these functions:
//
//   void addRequest(int floor, Direction direction);
//   void step();
//   int getCurrentFloor();
//   ElevatorState getState();
//
// How you implement them internally is up to you.
// -------------------------------------------------------------------------


