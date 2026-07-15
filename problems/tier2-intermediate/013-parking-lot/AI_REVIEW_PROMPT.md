# AI Code Review — Parking Lot System

## Context for the AI
You are reviewing a solution to "Parking Lot System" — an LLD interview problem testing Strategy + Factory patterns with HashMap + heap-style nearest-first allocation as the DSA core. This was asked at Amazon, Flipkart, Walmart, Salesforce.

The candidate was given a multi-part problem:
- Part 1: Multi-floor parking with vehicle-spot matching — smallest compatible spot, nearest floor first; default fee of 1.0/second
- Part 2: Pluggable pricing strategies (FlatRate, Hourly, Tiered) and gate tracking with entry/exit association

## Review Criteria

### 1. Pattern Correctness
- Is the PricingStrategy interface pure abstract with a single calculateFee(long) method?
- Are FlatRate, Hourly, and Tiered correctly implemented and fully decoupled from ParkingLot?
- Is SpotFactory used to centralize spot creation, or is construction scattered?
- Does ParkingLot correctly delegate fee calculation to the strategy without any if-else on type?

### 2. Open/Closed Principle
- Can a new pricing strategy (e.g., WeekendSurcharge) be added without modifying ParkingLot?
- Does the design survive both extension requirements without structural changes?

### 3. C++ Quality
- Memory management: who owns PricingStrategy* — does ParkingLot delete it?
- Const correctness on getAvailableSpots, getAvailableSpotsByFloor, getGates
- STL usage: vector<vector<ParkingSpot>> for floor-ordered storage vs. alternatives

### 4. Extension Handling
- How naturally did the base design absorb the pricing strategy in Part 2?
- Did the ticket struct accommodate entryGateId/exitGateId without redesign?
- What had to change when adding gate management? What survived?

### 5. Interview Readiness
- Could the candidate explain the spot compatibility matrix and nearest-first allocation verbally?
- What follow-up questions would expose weak understanding?
  - "How would you support reserved spots that only accept a specific license plate?"
  - "How would you implement a monthly pass that bypasses fee calculation?"
  - "What if two vehicles arrive simultaneously — is your design thread-safe?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->
