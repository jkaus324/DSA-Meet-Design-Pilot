#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <queue>
#include <algorithm>
#include <tuple>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct Tweet {
    int tweetId;
    int timestamp;
};

// ─── Optimized Twitter System ───────────────────────────────────────────────
// HINT: The key optimization is in getNewsFeed.
// HINT: Use a priority_queue (max-heap) to merge K sorted tweet streams.
// HINT: Each heap entry needs: timestamp, tweetId, userId, and an index
//       pointing to the next tweet from that user.

class Twitter {
    // HINT: Same data structures as Part 1 for storage
    // HINT: The optimization is purely in how you READ, not how you WRITE

public:
    Twitter() {
        // TODO: Initialize the system
    }

    void postTweet(int userId, int tweetId) {
        // TODO: Same as Part 1 — store tweet with auto-incrementing timestamp
    }

    vector<int> getNewsFeed(int userId) {
        // TODO: Build a set of all relevant users (self + following)
        // HINT: For each user with tweets, push their MOST RECENT tweet into a max-heap
        // HINT: Heap element = (timestamp, tweetId, userId, indexInTweetList)
        // HINT: Pop the max, add to result, push that user's NEXT tweet
        // HINT: Stop after 10 results or when heap is empty
        return {};
    }

    void follow(int followerId, int followeeId) {
        // TODO: Same as Part 1
    }

    void unfollow(int followerId, int followeeId) {
        // TODO: Same as Part 1
    }
};
