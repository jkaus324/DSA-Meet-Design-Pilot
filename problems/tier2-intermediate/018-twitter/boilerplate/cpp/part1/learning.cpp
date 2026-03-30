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

// ─── Twitter System ─────────────────────────────────────────────────────────

class Twitter {
    int time;
    unordered_map<int, vector<Tweet>> tweets;        // userId -> list of tweets
    unordered_map<int, unordered_set<int>> follows;   // userId -> set of followee IDs

public:
    Twitter() : time(0) {}

    void postTweet(int userId, int tweetId) {
        // TODO: Append a new Tweet{tweetId, time++} to this user's tweet list
    }

    vector<int> getNewsFeed(int userId) {
        // TODO: Collect all tweets from this user
        // TODO: Collect all tweets from users in follows[userId]
        // TODO: Sort collected tweets by timestamp descending (most recent first)
        // TODO: Return the tweetIds of the first 10 (or fewer) tweets
        return {};
    }

    void follow(int followerId, int followeeId) {
        // TODO: If followerId != followeeId, insert followeeId into follows[followerId]
    }

    void unfollow(int followerId, int followeeId) {
        // TODO: Erase followeeId from follows[followerId] (safe even if not present)
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Twitter System — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
