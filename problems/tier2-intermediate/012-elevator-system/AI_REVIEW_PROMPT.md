# AI Code Review — Elevator System

## Context for the AI
You are reviewing a solution to "Elevator System" — an LLD interview problem testing State + Strategy patterns with a Priority Queue + State Machine as the DSA core. This was asked at Amazon, Microsoft, Flipkart, Adobe.

The candidate was given a multi-part problem:
- Part 1: Single elevator with SCAN ordering — serve all floors in the current direction, then reverse
- Part 2: Multiple elevators managed by a pluggable DispatchStrategy (NearestFirst, LeastLoaded)

## Review Criteria

### 1. Pattern Correctness
- Does the Elevator implement a proper state machine with IDLE → MOVING → DOOR_OPEN transitions?
- Is the SCAN algorithm correctly implemented using separate up/down request sets?
- Does the DispatchStrategy interface allow new strategies to be added without touching ElevatorSystem?
- Are strategy objects cleanly separated from the Elevator and ElevatorSystem classes?

### 2. Open/Closed Principle
- Can a new dispatch strategy (e.g., RoundRobin) be added without modifying ElevatorSystem?
- Does the design survive the multi-elevator extension naturally?

### 3. C++ Quality
- Memory management: does ElevatorSystem properly delete elevator pointers on destruction?
- Use of `std::set<int>` for SCAN ordering vs. `std::priority_queue` — which is more appropriate and why?
- Const correctness on getters (getCurrentFloor, getState, etc.)

### 4. Extension Handling
- How cleanly does the single-elevator SCAN logic extend when wrapped in ElevatorSystem?
- Did the candidate design Elevator to be independently testable (no coupling to ElevatorSystem)?
- What had to change when adding DispatchStrategy? What survived?

### 5. Interview Readiness
- Could the candidate explain the SCAN algorithm and the state transition diagram verbally?
- What follow-up questions would expose weak understanding?
  - "What happens if a DOOR_OPEN step is skipped — does the system still work?"
  - "How would you add a priority emergency floor request?"
  - "How does NearestFirst handle an elevator moving in the wrong direction?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->
