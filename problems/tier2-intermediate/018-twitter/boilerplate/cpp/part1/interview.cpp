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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Design and implement a Twitter class that:
//   1. Lets users post tweets (postTweet)
//   2. Lets users follow/unfollow others (follow, unfollow)
//   3. Returns the 10 most recent tweets in a user's feed (getNewsFeed)
//
// The feed includes the user's own tweets AND tweets from users they follow,
// ordered by most recent first. Return at most 10 tweet IDs.
//
// Think about:
//   - How do you store user-to-tweets and user-to-following relationships?
//   - How do you order tweets by recency?
//   - What happens when a user unfollows someone?
//   - Can a user follow themselves?
//
// Entry points (must exist for tests):
//   Twitter()
//   void postTweet(int userId, int tweetId)
//   vector<int> getNewsFeed(int userId)
//   void follow(int followerId, int followeeId)
//   void unfollow(int followerId, int followeeId)
//
// ─────────────────────────────────────────────────────────────────────────────


