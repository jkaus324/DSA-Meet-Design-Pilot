# AI Code Review — Simplified Twitter System

## Context for the AI
You are reviewing a solution to "Simplified Twitter System" — an LLD interview problem testing Observer + Strategy patterns with HashMap + Heap (Merge K Sorted) as the DSA core. This was asked at Twitter, Amazon, Flipkart, AngelOne.

The candidate was given a multi-part problem:
- Part 1: Core social media operations — postTweet, getNewsFeed (naive), follow, unfollow
- Part 2: Optimized getNewsFeed using K-way heap merge — O(10 log K) instead of O(N log N)

## Review Criteria

### 1. Pattern Correctness
- Is the heap entry a tuple containing (timestamp, tweetId, userId, index) to enable efficient "next tweet" retrieval?
- Is the max-heap correctly seeded with each relevant user's most recent tweet?
- After popping the top tweet, does the implementation push the same user's previous tweet (index - 1)?
- Are users auto-created on first interaction (no explicit user registration)?

### 2. The Junction Moment
- Does the candidate understand WHY storing tweets as a vector (append-only, index-addressable) is the key design choice?
- Returning a `List<Tweet>` (sorted set) from tweet storage would destroy the k-way merge — does the candidate articulate this?
- Could the candidate explain the O(10 log K) vs O(N log N) trade-off verbally?

### 3. C++ Quality
- Heap comparator: lambda or struct? Is it correctly defined for a max-heap?
- Use of `unordered_set<int>` for follow relationships — O(1) follow/unfollow
- Use of `unordered_map<int, vector<Tweet>>` — correct choice for O(1) tweet append

### 4. Extension Handling
- How naturally did the Part 1 naive approach refactor to the Part 2 heap merge?
- Did the data structure choice in Part 1 (vector vs. deque vs. list) affect the difficulty of Part 2?
- What changed between Part 1 and Part 2? Only getNewsFeed — or did data structures change too?

### 5. Interview Readiness
- Could the candidate trace through the heap merge algorithm step by step with an example?
- What follow-up questions would expose weak understanding?
  - "What if a user follows 10,000 others but none have posted — how does your heap handle it?"
  - "How would you add a 'retweet' feature — what changes in the data model?"
  - "Could you make getNewsFeed lazy — only compute the next tweet when asked?"
- Rate: Hire / Lean Hire / Lean No Hire / No Hire

## My Solution

```cpp
// PASTE YOUR SOLUTION HERE
```

## My Approach
<!-- Describe your thought process in 2-3 sentences -->

## Specific Questions
<!-- Ask the AI anything specific about your solution -->
