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

// ─── Your Design Starts Here ─────────────────────────────────────────────────
//
// Optimize the Twitter news feed using a k-way merge algorithm.
//
// Each user's tweets are stored in chronological order (a sorted stream).
// Instead of collecting ALL tweets and sorting, merge K streams using a heap:
//
//   1. Push the most recent tweet from each relevant user into a max-heap
//   2. Pop the top (most recent overall)
//   3. Push that user's next tweet into the heap
//   4. Repeat until you have 10 results or heap is empty
//
// Think about:
//   - What goes into each heap entry? (timestamp, tweetId, userId, index)
//   - Why is this O(10 * log K) instead of O(N log N)?
//   - How do you track "the next tweet" for each user in the heap?
//
// Entry points (must exist for tests — same as Part 1):
//   Twitter()
//   void postTweet(int userId, int tweetId)
//   vector<int> getNewsFeed(int userId)
//   void follow(int followerId, int followeeId)
//   void unfollow(int followerId, int followeeId)
//
// ─────────────────────────────────────────────────────────────────────────────


