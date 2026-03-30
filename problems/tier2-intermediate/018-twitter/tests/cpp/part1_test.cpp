// Part 1 Tests — Core Twitter System
// Tests posting tweets, follow/unfollow, and basic news feed generation

#include "solution.cpp"
#include <cassert>
#include <iostream>
using namespace std;

int part1_tests() {
    int passed = 0;
    int failed = 0;

    // Test 1: Post a tweet and retrieve it in the feed
    try {
        Twitter twitter;
        twitter.postTweet(1, 101);
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 1);
        assert(feed[0] == 101);
        cout << "PASS test_post_and_get_feed" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_post_and_get_feed" << endl;
        failed++;
    }

    // Test 2: Follow a user and see their tweets in the feed
    try {
        Twitter twitter;
        twitter.postTweet(1, 101);
        twitter.postTweet(2, 201);
        twitter.follow(1, 2);
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 2);
        // Most recent first: tweet 201 was posted after 101
        assert(feed[0] == 201);
        assert(feed[1] == 101);
        cout << "PASS test_follow_and_feed" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_follow_and_feed" << endl;
        failed++;
    }

    // Test 3: Unfollow removes tweets from feed
    try {
        Twitter twitter;
        twitter.postTweet(1, 101);
        twitter.postTweet(2, 201);
        twitter.follow(1, 2);
        twitter.unfollow(1, 2);
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 1);
        assert(feed[0] == 101);  // only user 1's own tweet
        cout << "PASS test_unfollow_removes_from_feed" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_unfollow_removes_from_feed" << endl;
        failed++;
    }

    // Test 4: Feed returns at most 10 tweets
    try {
        Twitter twitter;
        for (int i = 1; i <= 15; i++) {
            twitter.postTweet(1, i);
        }
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 10);
        // Most recent first: 15, 14, 13, ..., 6
        assert(feed[0] == 15);
        assert(feed[9] == 6);
        cout << "PASS test_feed_max_10" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_feed_max_10" << endl;
        failed++;
    }

    // Test 5: Self-follow is a no-op
    try {
        Twitter twitter;
        twitter.postTweet(1, 101);
        twitter.follow(1, 1);  // should be ignored
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 1);
        assert(feed[0] == 101);
        cout << "PASS test_self_follow_noop" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_self_follow_noop" << endl;
        failed++;
    }

    // Test 6: Empty feed for user with no tweets and no follows
    try {
        Twitter twitter;
        auto feed = twitter.getNewsFeed(1);
        assert(feed.empty());
        cout << "PASS test_empty_feed" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_empty_feed" << endl;
        failed++;
    }

    // Test 7: Unfollow someone not followed is a no-op
    try {
        Twitter twitter;
        twitter.postTweet(1, 101);
        twitter.unfollow(1, 2);  // never followed 2
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 1);
        assert(feed[0] == 101);
        cout << "PASS test_unfollow_noop" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_unfollow_noop" << endl;
        failed++;
    }

    // Test 8: Feed merges tweets from multiple followed users in order
    try {
        Twitter twitter;
        twitter.postTweet(1, 10);   // timestamp 0
        twitter.postTweet(2, 20);   // timestamp 1
        twitter.postTweet(3, 30);   // timestamp 2
        twitter.postTweet(2, 21);   // timestamp 3
        twitter.follow(1, 2);
        twitter.follow(1, 3);
        auto feed = twitter.getNewsFeed(1);
        assert(feed.size() == 4);
        assert(feed[0] == 21);  // most recent
        assert(feed[1] == 30);
        assert(feed[2] == 20);
        assert(feed[3] == 10);  // oldest
        cout << "PASS test_multi_follow_feed_order" << endl;
        passed++;
    } catch (...) {
        cout << "FAIL test_multi_follow_feed_order" << endl;
        failed++;
    }

    cout << "PART1_SUMMARY " << passed << "/" << (passed + failed) << endl;
    return failed;
}
