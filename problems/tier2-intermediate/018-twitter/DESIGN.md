# Design Walkthrough — Simplified Twitter System

> This file is the answer guide. Only read after you've attempted the problem.

---

## The Core Design Decision

Twitter's news feed is fundamentally a **merge of K sorted lists** problem. Each user's tweets are naturally sorted by time (most recent first). The challenge is merging these sorted streams efficiently.

```
Twitter (Singleton — one system instance)
    ├── UserMap: HashMap<userId, User>
    │       └── User
    │           ├── tweets: LinkedList<Tweet>  ← sorted by time (newest first)
    │           └── following: HashSet<userId>
    └── News Feed Generation
            └── K-Way Merge via Max-Heap
```

The key structural decisions:
1. **HashMap for user lookup** — O(1) access to any user's data
2. **HashSet for follow relationships** — O(1) follow/unfollow/check
3. **LinkedList (or vector) for tweets** — O(1) prepend for new tweets, natural ordering by recency
4. **Max-Heap for feed merge** — efficiently merge K sorted streams without sorting all tweets

---

## Part 1: Naive Implementation

The simplest correct solution collects all relevant tweets, sorts them, and returns the top 10.

```cpp
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <algorithm>
using namespace std;

struct Tweet {
    int tweetId;
    int timestamp;
};

class Twitter {
    int time;
    unordered_map<int, vector<Tweet>> tweets;       // userId -> list of tweets
    unordered_map<int, unordered_set<int>> follows;  // userId -> set of followees

public:
    Twitter() : time(0) {}

    void postTweet(int userId, int tweetId) {
        tweets[userId].push_back({tweetId, time++});
    }

    vector<int> getNewsFeed(int userId) {
        // Collect all relevant tweets
        vector<Tweet> feed;
        // User's own tweets
        for (auto& t : tweets[userId]) feed.push_back(t);
        // Followed users' tweets
        for (int followee : follows[userId]) {
            for (auto& t : tweets[followee]) feed.push_back(t);
        }
        // Sort by timestamp descending
        sort(feed.begin(), feed.end(), [](const Tweet& a, const Tweet& b) {
            return a.timestamp > b.timestamp;
        });
        // Return top 10
        vector<int> result;
        for (int i = 0; i < min(10, (int)feed.size()); i++) {
            result.push_back(feed[i].tweetId);
        }
        return result;
    }

    void follow(int followerId, int followeeId) {
        if (followerId != followeeId) {
            follows[followerId].insert(followeeId);
        }
    }

    void unfollow(int followerId, int followeeId) {
        follows[followerId].erase(followeeId);
    }
};
```

**Complexity:** `getNewsFeed` is O(N log N) where N = total tweets from user + followed users. Acceptable for small scale, but doesn't pass the optimization bar.

---

## Part 2: K-Way Merge with Max-Heap

The optimized approach treats each user's tweet list as a sorted stream and uses a priority queue to merge them.

```cpp
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <queue>
using namespace std;

struct Tweet {
    int tweetId;
    int timestamp;
};

class Twitter {
    int time;
    unordered_map<int, vector<Tweet>> tweets;
    unordered_map<int, unordered_set<int>> follows;

public:
    Twitter() : time(0) {}

    void postTweet(int userId, int tweetId) {
        tweets[userId].push_back({tweetId, time++});
    }

    vector<int> getNewsFeed(int userId) {
        // Max-heap: (timestamp, tweetId, userId, index-in-that-user's-list)
        auto cmp = [](const tuple<int,int,int,int>& a, const tuple<int,int,int,int>& b) {
            return get<0>(a) < get<0>(b);  // max-heap by timestamp
        };
        priority_queue<tuple<int,int,int,int>, vector<tuple<int,int,int,int>>, decltype(cmp)> pq(cmp);

        // Collect all users whose tweets should appear
        unordered_set<int> users = follows[userId];
        users.insert(userId);  // include self

        // Push the most recent tweet from each user
        for (int uid : users) {
            if (!tweets[uid].empty()) {
                int idx = tweets[uid].size() - 1;  // most recent is at the end
                pq.push({tweets[uid][idx].timestamp, tweets[uid][idx].tweetId, uid, idx});
            }
        }

        vector<int> result;
        while (!pq.empty() && result.size() < 10) {
            auto [ts, tid, uid, idx] = pq.top();
            pq.pop();
            result.push_back(tid);
            // Push this user's next (older) tweet
            if (idx > 0) {
                int nextIdx = idx - 1;
                pq.push({tweets[uid][nextIdx].timestamp, tweets[uid][nextIdx].tweetId, uid, nextIdx});
            }
        }
        return result;
    }

    void follow(int followerId, int followeeId) {
        if (followerId != followeeId) {
            follows[followerId].insert(followeeId);
        }
    }

    void unfollow(int followerId, int followeeId) {
        follows[followerId].erase(followeeId);
    }
};
```

**Complexity Analysis:**
- `postTweet`: O(1) — just append
- `follow/unfollow`: O(1) — HashSet insert/erase
- `getNewsFeed`: O(K + 10 * log K) where K = number of followed users
  - K to push initial elements into the heap
  - At most 10 pops, each O(log K) for heap rebalance
  - This is dramatically better than O(N log N) for power users

---

## Pattern Insights

### Why K-Way Merge?
This is the same algorithm used in:
- **Merge K sorted linked lists** (LeetCode 23)
- **Database query merging** in distributed systems
- **External sort** when data doesn't fit in memory

The Twitter feed is a real-world application of this classic algorithm.

### Why HashMap + HashSet?
- **HashMap** for user-to-tweets mapping gives O(1) lookup per user
- **HashSet** for follow relationships gives O(1) follow/unfollow
- Together, they make the "collect relevant users" step O(1) per user

### Observer Pattern Connection
Although not explicitly implemented here, Twitter's architecture naturally maps to Observer:
- Users are **observers** (subscribers)
- Tweet posters are **subjects** (publishers)
- The follow/unfollow mechanism is subscribe/unsubscribe
- The news feed is the notification delivery mechanism

In a production system, you'd use Observer to push tweets to followers' feeds in real-time (fan-out on write) rather than pulling at read time (fan-out on read).

---

## What Interviewers Look For

1. **Correct data modeling** — Did you choose the right data structures for each relationship?
2. **Optimization awareness** — Can you explain why the naive approach is slow and how the heap-based approach improves it?
3. **Edge cases** — Self-follow prevention, unfollow no-op, user with no tweets, user with no followers
4. **Complexity analysis** — Can you derive O(10 * log K) and explain each component?
5. **Trade-off discussion** — Fan-out on write vs. fan-out on read, when to use each

---

## Common Interview Follow-ups

- *"What if the feed needs to include retweets?"* — Add a tweet type field; retweets reference the original tweet and carry the retweeter's timestamp
- *"How would you handle real-time updates?"* — Observer pattern: when a user posts, notify all followers (push model). Use WebSockets for live feed updates
- *"What about tweet deletion?"* — Soft delete (mark as deleted, skip in feed) vs. hard delete (remove from list, invalidate caches)
- *"How would you scale this to millions of users?"* — Shard by userId, use a distributed cache (Redis) for hot feeds, hybrid push/pull model based on follower count
