# AI Code Review — Amazon Locker System

## Context for the AI
You are reviewing a solution to "Amazon Locker System" — an LLD interview problem testing Strategy + State patterns with HashMap + Queue as the DSA core. This was asked at Amazon.

The candidate was given a multi-part problem:
- Part 1: Core locker allocation (smallest-fit) and retrieval by pickup code
- Part 2: Code expiry with configurable TTL, and Observer-based notification channels

## Review Criteria

### 1. Pattern Correctness
- Is the LockerAllocationStrategy interface cleanly defined so that new allocation rules (e.g., random, round-robin) require no changes to LockerSystem?
- Does SmallestFitStrategy correctly implement the SMALL→MEDIUM→LARGE fallback?
- Is the NotificationChannel interface a clean Observer — does LockerSystem iterate observers without knowing their concrete types?
- Are notifications fired at the correct times (deposit and expiry, not retrieval)?

### 2. Open/Closed Principle
- Can a new allocation strategy (e.g., PreferLargestFit) be added without modifying LockerSystem?
- Can a new notification channel (e.g., SMS) be registered without changing any existing code?

### 3. C++ Quality
- Use of `map<LockerSize, queue<string>>` for per-size available locker queues
- Memory management: does LockerSystem delete the strategy in its destructor?
- Use of `static LockerSystem* g_system` — does initLockerSystem() correctly delete and recreate?

### 4. Extension Handling
- How naturally did the base allocation logic accommodate code expiry in Part 2?
- Is the expiry check (currentTime - depositTime > expiryHours * 3600) correct for edge cases (exactly at expiry boundary)?
- What had to change for Part 2? What from Part 1 survived unchanged?

### 5. Interview Readiness
- Could the candidate explain the smallest-fit allocation strategy and when it falls back to a larger locker?
- What follow-up questions would expose weak understanding?
  - "What happens if two delivery agents deposit simultaneously — is your code thread-safe?"
  - "How would you support a 'preferred locker' feature where a customer requests a specific location?"
  - "What if checkExpired is called very frequently — how would you optimize it?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->
