// Part 2 Tests — Optimized News Feed with K-Way Merge
// Tests that the optimized implementation produces the same correct results
// The interface is identical to Part 1; the optimization is internal.

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part2_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Basic k-way merge — two users interleaved
    try {
        Twitter twitter;
        twitter.postTweet(1, 10);  // ts 0
        twitter.postTweet(2, 20);  // ts 1
        twitter.postTweet(1, 11);  // ts 2
        twitter.postTweet(2, 21);  // ts 3
        twitter.follow(1, 2);
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 4);
        assert(feed[0] == 21);
        assert(feed[1] == 11);
        assert(feed[2] == 20);
        assert(feed[3] == 10);
        cout << "PASS test_kway_merge_interleaved" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_kway_merge_interleaved" << endl;
        failed++;
    }

    // Test 2: K-way merge with many users — top 10 from 5 users
    try {
        Twitter twitter;
        // Each user posts 3 tweets
        for (int u = 1; u <= 5; u++) {
            for (int t = 0; t < 3; t++) {
                twitter.postTweet(u, u * 100 + t);
            }
        }
        // User 1 follows everyone
        for (int u = 2; u <= 5; u++) {
            twitter.follow(1, u);
        }
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 10);
        // Most recent 10 tweets: user 5's last tweet was most recent,
        // then user 5's second, etc.
        // Total 15 tweets, we get top 10 by timestamp
        assert(feed[0] == 502);  // user 5, tweet 2 (most recent overall)
        assert(feed[1] == 501);
        assert(feed[2] == 500);
        assert(feed[3] == 402);
        assert(feed[4] == 401);
        cout << "PASS test_kway_merge_many_users" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_kway_merge_many_users" << endl;
        failed++;
    }

    // Test 3: Feed after follow then unfollow
    try {
        Twitter twitter;
        twitter.postTweet(1, 10);
        twitter.postTweet(2, 20);
        twitter.postTweet(3, 30);
        twitter.follow(1, 2);
        twitter.follow(1, 3);
        auto feed1 = twitter.getNewsFeed(1);
        assert(feed1.size() == 3);

        twitter.unfollow(1, 3);
        auto feed2 = twitter.getNewsFeed(1);
        assert(feed2.size() == 2);
        assert(feed2[0] == 20);
        assert(feed2[1] == 10);
        cout << "PASS test_kway_after_unfollow" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_kway_after_unfollow" << endl;
        failed++;
    }

    // Test 4: Only user's own tweets when following nobody
    try {
        Twitter twitter;
        twitter.postTweet(1, 10);
        twitter.postTweet(1, 11);
        twitter.postTweet(2, 20);  // user 2 exists but user 1 doesn't follow
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 2);
        assert(feed[0] == 11);
        assert(feed[1] == 10);
        cout << "PASS test_own_tweets_only" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_own_tweets_only" << endl;
        failed++;
    }

    // Test 5: Exactly 10 tweets when more are available
    try {
        Twitter twitter;
        for (int i = 1; i <= 20; i++) {
            twitter.postTweet(1, i);
        }
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 10);
        for (int i = 0; i < 10; i++) {
            assert(feed[i] == 20 - i);  // 20, 19, 18, ..., 11
        }
        cout << "PASS test_exactly_10" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_exactly_10" << endl;
        failed++;
    }

    // Test 6: Follow after posting — new tweets appear in feed
    try {
        Twitter twitter;
        twitter.postTweet(2, 200);
        twitter.postTweet(2, 201);
        // User 1 hasn't followed 2 yet
        auto feed1 = twitter.getNewsFeed(1);
        assert(feed1.empty());
        // Now follow
        twitter.follow(1, 2);
        auto feed2 = twitter.getNewsFeed(1);
        assert(feed2.size() == 2);
        assert(feed2[0] == 201);
        assert(feed2[1] == 200);
        cout << "PASS test_follow_after_posting" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_follow_after_posting" << endl;
        failed++;
    }

    // Test 7: Stress — merge from 10 users, 5 tweets each, get top 10
    try {
        Twitter twitter;
        for (int u = 1; u <= 10; u++) {
            for (int t = 0; t < 5; t++) {
                twitter.postTweet(u, u * 1000 + t);
            }
        }
        for (int u = 2; u <= 10; u++) {
            twitter.follow(1, u);
        }
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 10);
        // The 10 most recent tweets are the last tweet from each of the 10 users
        // posted in order: user1 t0..t4, user2 t0..t4, ... user10 t0..t4
        // Most recent: user10 t4, user10 t3, user10 t2, user10 t1, user10 t0,
        //              user9 t4, user9 t3, user9 t2, user9 t1, user9 t0
        assert(feed[0] == 10004);  // user 10, tweet 4
        assert(feed[9] == 9000);   // user 9, tweet 0
        cout << "PASS test_stress_10_users" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_stress_10_users" << endl;
        failed++;
    }

    cout << "PART2_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
