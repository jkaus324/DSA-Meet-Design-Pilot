#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <algorithm>
using namespace std;

// ─── Data Model (given — do not modify) ─────────────────────────────────────

struct Tweet {
    int tweetId;
    int timestamp;
};

// ─── System Class ───────────────────────────────────────────────────────────
// HINT: You need two maps — one to store each user's tweets,
// and one to store who each user follows.
// HINT: Use a global counter to assign timestamps so tweets can be ordered.

class Twitter {
    // HINT: What data structure gives O(1) lookup by userId?
    // HINT: What data structure lets you quickly check "does user A follow user B"?
    // HINT: How do you ensure tweets are ordered by recency?

public:
    Twitter() {
        // TODO: Initialize the system
    }

    void postTweet(int userId, int tweetId) {
        // TODO: Record this tweet for the given user
        // HINT: Auto-create the user if they don't exist
        // HINT: Assign a timestamp so you can order tweets later
    }

    vector<int> getNewsFeed(int userId) {
        // TODO: Collect tweets from this user and all users they follow
        // HINT: Gather all relevant tweets into one collection
        // HINT: Sort by timestamp (most recent first)
        // HINT: Return at most 10 tweet IDs
        return {};
    }

    void follow(int followerId, int followeeId) {
        // TODO: Add followeeId to followerId's following set
        // HINT: A user should not be able to follow themselves
    }

    void unfollow(int followerId, int followeeId) {
        // TODO: Remove followeeId from followerId's following set
        // HINT: Unfollowing someone you don't follow should be a no-op
    }
};
