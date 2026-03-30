# Problem 018 — Simplified Twitter System

**Tier:** 2 (Intermediate) | **Patterns:** Observer, Factory, Singleton | **DSA:** HashMap, HashSet, Heap, LinkedList
**Companies:** AngelOne | **Time:** 50 minutes

---

## Problem Statement

You are building a simplified version of Twitter. The system must support users posting tweets, following and unfollowing other users, and retrieving a personalized news feed.

Each user has a unique integer `userId`. Each tweet has a unique integer `tweetId`. Tweets are ordered by recency — the most recently posted tweet appears first.

**Your task:** Design and implement a `Twitter` class that supports the core social media operations described below, and then optimize the news feed generation using a k-way merge algorithm.

---

## Before You Code

> Read this section carefully. This is where the design thinking happens.

**Ask yourself:**
1. How do you store the relationship between users and their tweets efficiently?
2. When generating a news feed, you need tweets from potentially many followed users — how do you avoid sorting *all* tweets every time?
3. What data structure lets you efficiently extract the top K elements from multiple sorted streams?

**Naive approach:** Collect all tweets from user + followed users into one list, sort by timestamp, return top 10. This is O(N log N) where N is total tweets across all followed users.

**Pattern approach:** Each user's tweet list is already sorted by time (newest first). Treat each user's tweets as a sorted stream. Use a **max-heap** (priority queue) to perform a **k-way merge** — push the newest tweet from each user, pop the max, push that user's next tweet. This gives O(10 * log K) where K is the number of followed users.

**The key insight:** A news feed is a merge of K sorted lists, limited to 10 results. This is the classic "merge K sorted lists" problem disguised as a system design question.

---

## Data Structures

```cpp
struct Tweet {
    int tweetId;
    int timestamp;  // auto-assigned, higher = more recent
};
```

---

## Part 1

**Base requirement — Core Twitter functionality**

Implement a `Twitter` class that supports the following operations:

| Operation | Description |
|-----------|-------------|
| `Twitter()` | Initialize the system |
| `postTweet(userId, tweetId)` | User posts a new tweet. If the user does not exist, auto-create them. |
| `getNewsFeed(userId)` | Return the 10 most recent tweet IDs from the user's feed. The feed includes the user's own tweets and tweets from users they follow, ordered by most recent first. If fewer than 10 tweets exist, return all of them. |
| `follow(followerId, followeeId)` | `followerId` starts following `followeeId`. A user cannot follow themselves. |
| `unfollow(followerId, followeeId)` | `followerId` stops following `followeeId`. Unfollowing a user you don't follow is a no-op. |

**Constraints:**
- `1 <= userId, tweetId <= 10^4`
- Each call to `postTweet` uses a unique `tweetId`
- At most 3 * 10^4 calls total to `postTweet`, `getNewsFeed`, `follow`, `unfollow`

**Entry points (tests will call these):**
```cpp
Twitter();
void postTweet(int userId, int tweetId);
vector<int> getNewsFeed(int userId);
void follow(int followerId, int followeeId);
void unfollow(int followerId, int followeeId);
```

**Design goal:** The naive approach (collect all tweets, sort, return top 10) is acceptable for Part 1. Focus on correct data modeling.

---

## Part 2

**Extension — Optimized news feed with k-way merge**

The product team reports that power users follow thousands of accounts. The naive `getNewsFeed` is too slow.

Optimize `getNewsFeed` using a **min-heap** based k-way merge:

1. For the user and each user they follow, take the most recent tweet as the head of that user's stream.
2. Push all stream heads into a max-heap (priority queue ordered by timestamp).
3. Pop the top element (most recent overall). If that user has more tweets, push their next tweet.
4. Repeat until you have 10 tweets or the heap is empty.

**Complexity target:**
- `getNewsFeed`: O(10 * log K) where K = number of followed users
- `postTweet`: O(1)
- `follow` / `unfollow`: O(1)

**Entry points:** Same as Part 1, but with optimized implementation.

```cpp
Twitter();
void postTweet(int userId, int tweetId);
vector<int> getNewsFeed(int userId);
void follow(int followerId, int followeeId);
void unfollow(int followerId, int followeeId);
```

**Discussion points for the interviewer:**
- Why does each user's tweet list need to be in reverse chronological order?
- What happens if a user unfollows someone — do their tweets disappear from future feeds?
- How would you extend this to support "liked" tweets or retweets in the feed?

---

## Running Tests

```bash
./run-tests.sh 018-twitter cpp
```
