# Problem 012 — Elevator System

**Tier:** 2 (Intermediate) | **Pattern:** State + Strategy + Command | **DSA:** Queue + PriorityQueue + HashMap
**Companies:** Adobe | **Time:** 45 minutes

---

## Problem Statement

You're designing an elevator control system for a commercial building. The system must manage elevator states, process floor requests in an efficient order, and scale to multiple elevators with pluggable dispatching strategies.

**Your task:** Design and implement an `ElevatorSystem` that handles floor requests, moves elevators using the SCAN algorithm, and dispatches requests across multiple elevators using swappable strategies.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. An elevator can be IDLE, MOVING_UP, MOVING_DOWN, or DOOR_OPEN. How do you model transitions between these states cleanly? What happens if you use a chain of if-else statements instead?
2. When processing requests, a naive approach serves them in FIFO order. But real elevators use SCAN — serve all floors in one direction before reversing. What data structure lets you efficiently find the next floor in the current direction?
3. When you add multiple elevators, how do you decide which elevator handles a request? Should the dispatching logic live inside the elevator? Or should it be a separate, swappable component?

**The key insight:** The **State** pattern models elevator behavior per state (idle elevators accept requests differently than moving ones). The **Strategy** pattern decouples the dispatching algorithm from the elevator itself. A priority queue (or sorted set) per direction enables efficient SCAN ordering.

---

## Data Structures

```cpp
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
    Direction direction;  // Direction the passenger wants to go (for external requests)
};
```

---

## Part 1

**Base requirement — Single elevator with SCAN ordering**

Implement an `Elevator` that starts at floor 0 in IDLE state. It accepts external requests (someone pressed UP/DOWN on a floor) and internal requests (someone pressed a floor button inside the elevator). Requests are processed in SCAN order: the elevator serves all pending floors in its current direction before reversing.

**SCAN rule:** If the elevator is moving up, it serves all upward requests in ascending order. When no more upward requests remain, it reverses to serve downward requests in descending order. If idle, the direction is determined by the first incoming request.

**Entry points (tests will call these):**
```cpp
void addRequest(int floor, Direction direction);
void step();              // Advance one step: move one floor or open/close doors
int getCurrentFloor();
ElevatorState getState();
```

**What to implement:**
```cpp
class Elevator {
    int currentFloor;
    ElevatorState state;
    Direction currentDirection;
    // Pending floors in each direction
public:
    Elevator();
    void addRequest(int floor, Direction direction);
    void step();
    int getCurrentFloor() const;
    ElevatorState getState() const;
};
```

**Behavior per step():**
- If IDLE and requests exist: set direction based on nearest request logic, transition to MOVING_UP or MOVING_DOWN.
- If MOVING_UP or MOVING_DOWN: move one floor in that direction. If the current floor has a pending request, transition to DOOR_OPEN and remove that request.
- If DOOR_OPEN: close doors and resume moving (or go IDLE if no more requests).

**Design goal:** The state transitions must be clean. Each state should have well-defined behavior. The SCAN algorithm should use a sorted structure (e.g., `set<int>`) for each direction to efficiently determine the next stop.

---

## Part 2

**Extension — Multiple elevators with dispatch strategies**

Scale to N elevators. When an external request arrives, the system must dispatch it to the best elevator using a pluggable strategy.

| Strategy | Rule |
|----------|------|
| NearestFirst | Assign to the nearest elevator that is idle or moving toward the request floor in the correct direction |
| LeastLoaded | Assign to the elevator with the fewest pending requests |

**Design challenge:** How do you add a new dispatching strategy without modifying the `ElevatorSystem` class?

**New entry points:**
```cpp
void addElevator(int id);
void setDispatchStrategy(DispatchStrategy* strategy);
// addRequest now dispatches to the best elevator
void addRequest(int floor, Direction direction);
void step();  // Steps ALL elevators
```

**What to implement:**
```cpp
class DispatchStrategy {
public:
    virtual int selectElevator(const vector<Elevator*>& elevators,
                               int requestFloor,
                               Direction requestDirection) = 0;
    virtual ~DispatchStrategy() = default;
};

class NearestFirst : public DispatchStrategy { ... };
class LeastLoaded : public DispatchStrategy { ... };

class ElevatorSystem {
    vector<Elevator*> elevators;
    DispatchStrategy* strategy;
public:
    void addElevator(int id);
    void setDispatchStrategy(DispatchStrategy* strategy);
    void addRequest(int floor, Direction direction);
    void step();
};
```

**Hint:** NearestFirst considers both distance and direction compatibility. An elevator moving UP at floor 3 is a good candidate for an UP request at floor 5, but not for a DOWN request at floor 1.

---

## Running Tests

```bash
./run-tests.sh 012-elevator-system cpp
```
