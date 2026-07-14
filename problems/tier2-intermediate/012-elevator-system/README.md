# Problem 012 — Elevator System

**Tier:** 2 (Intermediate) | **Pattern:** State + Strategy | **DSA:** Priority Queue + State Machine
**Companies:** Amazon, Microsoft, Flipkart, Adobe | **Time:** 60 minutes

---

## Problem Statement

You are designing an elevator control system for a commercial building. A single elevator processes floor requests using the SCAN algorithm — it serves all floors in the current direction before reversing. The system then scales to multiple elevators, with a pluggable strategy determining which elevator handles each incoming request.

**Constraints:**
- Floors are non-negative integers; building has at most 100 floors
- Each `step()` call advances the simulation by one unit: move one floor, or open/close doors
- An elevator in `DOOR_OPEN` state closes and resumes on the next step
- External requests specify a direction (UP or DOWN); internal requests specify only a destination floor

---

## Base Requirement — Single Elevator with SCAN Ordering

Implement an `Elevator` that starts at floor 0 in IDLE state. It accepts external requests (UP/DOWN button pressed on a floor) and internal requests (floor button inside the cabin). Requests are served in SCAN order: all pending floors in the current direction first, then reverse.

**SCAN rule:** Moving UP — serve pending floors in ascending order. When none remain, reverse to serve pending floors in descending order. If idle, the first request sets the initial direction.

**Example:**
```
elevator at floor 0, IDLE
addRequest(3, UP), addRequest(7, UP), addRequest(1, DOWN)
step() → moves to floor 1 (nearest in current UP direction would be 1, but UP means ascending: goes to 3)
// After UP pass: stops at 3 (DOOR_OPEN), then 7 (DOOR_OPEN)
// Then reverses DOWN: stops at 1 (DOOR_OPEN)
getCurrentFloor()  →  3 after first stop
getState()         →  DOOR_OPEN when at a requested floor
```

**Public methods:**
- `void addRequest(int floor, Direction direction)`
- `void step()`
- `int getCurrentFloor() const`
- `ElevatorState getState() const`

---

## Extension 1 — Multiple Elevators with Dispatch Strategies

Scale to N elevators managed by an `ElevatorSystem`. Each incoming external request is dispatched to the best elevator according to a pluggable strategy. Adding a new dispatch strategy must require zero changes to `ElevatorSystem`.

| Strategy | Rule |
|---|---|
| NearestFirst | Assign to the nearest idle elevator, or one already moving toward the request in the correct direction |
| LeastLoaded | Assign to the elevator with the fewest pending requests |

**Example:**
```
// 2 elevators: E1 at floor 0 (IDLE), E2 at floor 6 (MOVING_DOWN)
// Strategy: NearestFirst
addRequest(4, UP)
// E1 is IDLE at floor 0, distance=4
// E2 is MOVING_DOWN at floor 6, wrong direction for UP at 4
// → dispatched to E1
```

**Public methods:**
- `void addElevator(int id)`
- `void setDispatchStrategy(DispatchStrategy* strategy)`
- `void addRequest(int floor, Direction direction)`
- `void step()`  — advances all elevators one step

---

## Running Tests

```bash
./run-tests.sh 012-elevator-system cpp
```
