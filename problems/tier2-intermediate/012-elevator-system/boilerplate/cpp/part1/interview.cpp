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

// --- Your Design Starts Here ------------------------------------------------
//
// Design and implement a single Elevator that:
//   1. Starts at floor 0 in IDLE state
//   2. Accepts external requests (floor + direction) and internal requests
//   3. Processes requests in SCAN order: serve all floors in the current
//      direction before reversing
//
// step() behavior:
//   - IDLE + requests exist: pick direction, start moving
//   - MOVING_UP / MOVING_DOWN: move one floor; if current floor has a
//     pending request, transition to DOOR_OPEN
//   - DOOR_OPEN: close doors, resume moving or go IDLE if no requests
//
// Think about:
//   - What data structure efficiently tracks pending floors per direction?
//   - How do you decide when to reverse direction?
//   - What happens if a request arrives for the current floor while idle?
//
// Entry points (must exist for tests):
//   void addRequest(int floor, Direction direction);
//   void step();
//   int getCurrentFloor();
//   ElevatorState getState();
//
// -------------------------------------------------------------------------


