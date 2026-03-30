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
// Extend your Part 1 Elevator to support:
//   1. Multiple elevators managed by an ElevatorSystem
//   2. Pluggable dispatch strategies that decide which elevator handles
//      a new request
//   3. Two strategies: NearestFirst (nearest elevator in compatible
//      direction) and LeastLoaded (fewest pending requests)
//
// Think about:
//   - How do you define a strategy interface so new strategies can be
//     added without modifying the system?
//   - What information does a strategy need about each elevator to
//     make a decision?
//   - How does step() work when there are multiple elevators?
//
// Entry points (must exist for tests):
//   void addElevator(int id);
//   void setDispatchStrategy(DispatchStrategy* strategy);
//   void addRequest(int floor, Direction direction);
//   void step();
//   Elevator* getElevator(int index);
//   int getElevatorCount();
//
// -------------------------------------------------------------------------


