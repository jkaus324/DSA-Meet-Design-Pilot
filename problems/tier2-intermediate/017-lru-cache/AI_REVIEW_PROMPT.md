# AI Code Review — LRU Cache

## Context for the AI
You are reviewing a solution to "LRU Cache" — an LLD interview problem testing Strategy + Observer patterns with HashMap + Doubly Linked List as the DSA core. This was asked at Amazon, Google, Microsoft, Flipkart, Kutumb.

The candidate was given a multi-part problem:
- Part 1: O(1) LRU Cache — get and put with LRU eviction on capacity overflow
- Part 2: Per-entry TTL — lazy eviction on access; deleteKey; size()
- Part 3: Eviction listeners — Observer pattern for CAPACITY, TTL_EXPIRED, and EXPLICIT_DELETE events

## Review Criteria

### 1. Pattern Correctness
- Is the doubly linked list correctly implemented with sentinel head/tail nodes?
- Is addToFront/removeNode/moveToFront truly O(1)?
- Is eviction reason correctly assigned — CAPACITY for overflow, TTL_EXPIRED for lazy expiry, EXPLICIT_DELETE for deleteKey?
- Does the EvictionListener interface follow the Observer pattern cleanly?

### 2. Open/Closed Principle
- Can a new eviction listener (e.g., MetricsCollector) be added without modifying LRUCache?
- Does the design handle multiple simultaneous listeners correctly?

### 3. C++ Quality
- Memory management: are all Node* deleted properly in the destructor?
- Is the sentinel head/tail pattern used to avoid null-pointer edge cases?
- Use of `unordered_map<int, Node*>` — is this the right choice vs. `map`?
- Const correctness where applicable

### 4. Extension Handling
- Did the base design (Part 1) easily absorb TTL in Part 2? What changed?
- Is the lazy TTL eviction triggered correctly on both get() and put() with currentTime?
- Are Part 1's get(key) and Part 2's get(key, currentTime) overloads co-existing cleanly?

### 5. Interview Readiness
- Could the candidate draw the doubly linked list state after a sequence of get/put operations?
- What follow-up questions would expose weak understanding?
  - "Why use a doubly linked list instead of a singly linked list?"
  - "What happens if a listener modifies the cache during the onEviction callback?"
  - "How would you implement a proactive TTL scan instead of lazy eviction?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->
