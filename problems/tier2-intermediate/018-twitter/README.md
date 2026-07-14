# Problem 018 — Simplified Twitter System

**Tier:** 2 (Intermediate) | **Pattern:** Observer + Strategy | **DSA:** HashMap + Heap (Merge K Sorted)
**Companies:** Twitter, Amazon, Flipkart, AngelOne | **Time:** 60 minutes

---

## Problem Statement

You are building a simplified Twitter. Users post tweets, follow and unfollow other users, and retrieve a personalized news feed. The feed returns the 10 most recent tweets from the user and everyone they follow, in reverse chronological order. The naive approach collects and sorts all tweets — Part 2 optimizes this to O(10 log K) using a k-way heap merge.

**Constraints:**
- `1 <= userId, tweetId <= 10^4`
- At most 3 * 10^4 total calls across all operations
- Feed returns at most 10 tweet IDs; if fewer exist, return all
- A user cannot follow themselves; unfollowing someone you don't follow is a no-op
- Timestamp is auto-assigned; higher timestamp = more recent

---

## Base Requirement — Core Social Media Operations

Implement a `Twitter` class supporting posting, following, unfollowing, and feed generation. Users are auto-created on first interaction. The feed must include the user's own tweets.

**Example:**
```
Twitter tw
tw.postTweet(1, 101)   // user 1 posts tweet 101
tw.postTweet(2, 201)   // user 2 posts tweet 201
tw.postTweet(1, 102)   // user 1 posts tweet 102

tw.getNewsFeed(1)  →  [102, 101]   // own tweets, newest first, no follows yet

tw.follow(1, 2)
tw.getNewsFeed(1)  →  [102, 201, 101]   // user 2's tweet inserted in order

tw.unfollow(1, 2)
tw.getNewsFeed(1)  →  [102, 101]   // back to own tweets only
```

**Public methods:**
- `Twitter()`
- `void postTweet(int userId, int tweetId)`
- `vector<int> getNewsFeed(int userId)`
- `void follow(int followerId, int followeeId)`
- `void unfollow(int followerId, int followeeId)`

---

## Extension 1 — Optimized Feed with K-Way Heap Merge

The naive approach (collect all tweets, sort, take top 10) is O(N log N) where N is total tweets. Optimize `getNewsFeed` using a max-heap k-way merge.

**Algorithm:**
1. For each followed user (including self), push their most recent tweet onto a max-heap keyed by timestamp.
2. Pop the top (most recent overall). If that user has more tweets, push their next one.
3. Repeat until 10 tweets collected or heap empty.

**Complexity targets:**
- `getNewsFeed`: O(10 * log K), K = number of followed users
- `postTweet`: O(1)
- `follow` / `unfollow`: O(1)

**Example:**
```
// User 1 follows 500 users, each with 1000 tweets
tw.getNewsFeed(1)
// Heap starts with 501 entries (1 per user's latest tweet)
// Pops 10 times, each pop does at most log(501) comparisons
// Total: ~10 * 9 ≈ 90 comparisons vs. 500,000 sort comparisons
```

**Public methods:** Same interface as Part 1; only the internal implementation of `getNewsFeed` changes.

---

## Running Tests

```bash
./run-tests.sh 018-twitter cpp
```
