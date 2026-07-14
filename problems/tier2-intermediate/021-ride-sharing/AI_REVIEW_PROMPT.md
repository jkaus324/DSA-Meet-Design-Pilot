# AI Code Review — Ride-Sharing Application

## Context for the AI
You are reviewing a solution to "Ride-Sharing Application" — an LLD interview problem testing Strategy + Factory patterns with HashMap and Graph (BFS) data structures. This was asked at Flipkart.

The candidate was given a multi-part problem:
- Part 1: User/vehicle/ride onboarding — addUser, addVehicle, offerRide with uniqueness and active-vehicle constraints
- Part 2: Pluggable ride selection — MostVacantStrategy and PreferredVehicleStrategy behind a RideSelectionStrategy interface
- Part 3: Ride lifecycle — endRide, getRideStats, printRideStats with per-user counters

## Review Criteria

### 1. Pattern Correctness
- Is RideSelectionStrategy a true abstract interface with no selection logic leaking into RideService?
- Does PreferredVehicleStrategy receive the vehicle store by reference (not by copy or direct DB query)?
- Can a new selection strategy (e.g., "cheapest fare") be added by creating one new class only?
- Are the strategy objects created externally and injected — not newed inside RideService?

### 2. Open/Closed Principle
- Can a new selection strategy be added without modifying RideService::selectRide?
- Does the design survive Part 3 (endRide + stats) without restructuring Part 1 or Part 2?

### 3. C++ Quality
- Memory management: who owns the strategy pointer? Is there a dangling pointer risk?
- Const correctness: are read-only methods (getRideStats, printRideStats, hasUser) marked const?
- STL usage: unordered_map for O(1) user/vehicle/ride lookups vs vector linear scans

### 4. Extension Handling
- Is activeVehicles (regNumber → rideId) maintained correctly across offerRide and endRide?
- Does selectRide correctly exclude rides where driverId == passengerName?
- When availableSeats reaches 0, is the ride still in the rides map (just non-selectable)?
- How well does the base design absorb Part 3 — did stats tracking require major rewrites?

### 5. Interview Readiness
- Could the candidate explain the PreferredVehicleStrategy vehicle-store coupling decision verbally?
- What follow-up questions would expose weak understanding?
  - "What if we want to rank rides by multiple criteria simultaneously?"
  - "What if the same passenger selects two rides concurrently?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->
