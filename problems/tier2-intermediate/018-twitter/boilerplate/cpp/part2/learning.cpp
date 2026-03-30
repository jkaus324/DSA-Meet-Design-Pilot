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

class Twitter {
    int time;
    unordered_map<int, vector<Tweet>> tweets;
    unordered_map<int, unordered_set<int>> follows;

public:
    Twitter() : time(0) {}

    void postTweet(int userId, int tweetId) {
        // TODO: Append Tweet{tweetId, time++} to tweets[userId]
    }

    vector<int> getNewsFeed(int userId) {
        // Step 1: Build set of relevant users
        // TODO: Create a set containing userId and all users in follows[userId]

        // Step 2: Define max-heap comparator
        // The heap stores tuples of (timestamp, tweetId, userId, index)
        auto cmp = [](const tuple<int,int,int,int>& a, const tuple<int,int,int,int>& b) {
            return get<0>(a) < get<0>(b);  // max-heap by timestamp
        };
        priority_queue<tuple<int,int,int,int>, vector<tuple<int,int,int,int>>, decltype(cmp)> pq(cmp);

        // Step 3: Seed the heap
        // TODO: For each relevant user, if they have tweets, push their most
        //       recent tweet (last element) into the heap as
        //       (timestamp, tweetId, userId, index)

        // Step 4: K-way merge
        vector<int> result;
        while (!pq.empty() && result.size() < 10) {
            // TODO: Pop the top element
            // TODO: Add its tweetId to result
            // TODO: If that user has an older tweet (index > 0), push it
        }
        return result;
    }

    void follow(int followerId, int followeeId) {
        // TODO: If followerId != followeeId, insert into follows[followerId]
    }

    void unfollow(int followerId, int followeeId) {
        // TODO: Erase followeeId from follows[followerId]
    }
};

#ifndef RUNNING_TESTS
int main() {
    cout << "Optimized Twitter — implement the TODO methods above, then run tests." << endl;
    return 0;
}
#endif
