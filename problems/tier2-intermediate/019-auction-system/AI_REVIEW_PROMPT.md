# AI Code Review — Online Auction System

## Context for the AI
You are reviewing a solution to "Online Auction System" — an LLD interview problem testing Strategy + Observer + State patterns with Priority Queue + HashMap as the DSA core. This was asked at Flipkart, Amazon, eBay.

The candidate was given a multi-part problem:
- Part 1: Core auction — registerUser, createAuction, placeBid, getWinningBid
- Part 2: Auction lifecycle — closeAuction, getAuctionStatus (OPEN/CLOSED/NO_SALE), state guards
- Part 3: Three auction strategies — ASCENDING (default), SEALED (hidden bids), BUYNOW (auto-close at 1.5x)

## Review Criteria

### 1. Pattern Correctness
- Does AuctionStrategy encapsulate all three varying behaviors: acceptBid, getVisibleWinningBid, shouldAutoClose?
- Is the strategy stored per-auction (in strategies map) rather than per-system?
- Does createStrategy act as a factory — decoupling strategy creation from AuctionSystem?
- Does placeBid delegate all bid validation to the strategy, with no if-else on strategy type?

### 2. Open/Closed Principle
- Can a new auction type (e.g., Dutch Auction — price decreases until someone accepts) be added without modifying AuctionSystem?
- Does the design survive all three strategy variants without structural changes?

### 3. C++ Quality
- Memory management: does AuctionSystem delete all strategy pointers in its destructor?
- Use of `unordered_map<int, AuctionStrategy*>` for per-auction strategies
- Is BUYNOW auto-close implemented without a special case in AuctionSystem::placeBid?

### 4. Extension Handling
- How naturally did the base design (Part 1) absorb state management in Part 2?
- Were strategy-specific behaviors cleanly separated — did adding SEALED/BUYNOW require zero AuctionSystem changes?
- What had to change between parts? What survived unchanged?

### 5. Interview Readiness
- Could the candidate explain why shouldAutoClose belongs in the strategy rather than in AuctionSystem?
- What follow-up questions would expose weak understanding?
  - "How would you add a reserve price that must be met for an auction to close as SOLD?"
  - "How would you notify all bidders when an auction closes?"
  - "In SEALED bid, how do you prevent a bidder from placing multiple bids?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->
