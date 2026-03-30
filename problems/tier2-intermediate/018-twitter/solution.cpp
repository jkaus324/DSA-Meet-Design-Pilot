#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <queue>
#include <algorithm>
#include <tuple>
using namespace std;

// ─── Data Model ─────────────────────────────────────────────────────────────

struct Tweet {
    int tweetId;
    int timestamp;
};

// ─── Optimized Twitter System (K-Way Merge) ────────────────────────────────

class Twitter {
    int time;
    unordered_map<int, vector<Tweet>> tweets;
    unordered_map<int, unordered_set<int>> follows;

public:
    Twitter() : time(0) {}

    void postTweet(int userId, int tweetId) {
        tweets[userId].push_back(Tweet{tweetId, time++});
    }

    vector<int> getNewsFeed(int userId) {
        // Step 1: Build set of relevant users (self + followees)
        unordered_set<int> users;
        users.insert(userId);
        if (follows.count(userId)) {
            for (int followee : follows[userId]) {
                users.insert(followee);
            }
        }

        // Step 2: Max-heap by timestamp
        // Tuple: (timestamp, tweetId, userId, index into that user's tweet list)
        auto cmp = [](const tuple<int,int,int,int>& a, const tuple<int,int,int,int>& b) {
            return get<0>(a) < get<0>(b);
        };
        priority_queue<tuple<int,int,int,int>, vector<tuple<int,int,int,int>>, decltype(cmp)> pq(cmp);

        // Step 3: Seed the heap with each user's most recent tweet
        for (int uid : users) {
            if (tweets.count(uid) && !tweets[uid].empty()) {
                int idx = (int)tweets[uid].size() - 1;
                pq.push({tweets[uid][idx].timestamp, tweets[uid][idx].tweetId, uid, idx});
            }
        }

        // Step 4: K-way merge — extract up to 10 most recent tweets
        vector<int> result;
        while (!pq.empty() && result.size() < 10) {
            auto [ts, tid, uid, idx] = pq.top();
            pq.pop();
            result.push_back(tid);
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

#ifndef RUNNING_TESTS
int main() {
    cout << "Twitter System — run tests to verify." << endl;
    return 0;
}
#endif
