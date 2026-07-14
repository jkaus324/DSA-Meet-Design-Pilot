// Twitter — Part 1 Tests
import java.util.*;
import java.util.stream.*;

class Part1Test {
    static boolean testPostAndGetFeed() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 101);
            var feed = twitter.getNewsFeed(1);
            boolean pass = feed.size() == 1
                && feed[0] == 101;
            System.out.println((pass ? "PASS" : "FAIL") + ": testPostAndGetFeed");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testPostAndGetFeed (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFollowAndFeed() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 101);
            twitter.postTweet(2, 201);
            twitter.follow(1, 2);
            var feed = twitter.getNewsFeed(1);
            // Most recent first: tweet 201 was posted after 101
            boolean pass = feed.size() == 2
                && feed[0] == 201
                && feed[1] == 101;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFollowAndFeed");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFollowAndFeed (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testUnfollowRemovesFromFeed() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 101);
            twitter.postTweet(2, 201);
            twitter.follow(1, 2);
            twitter.unfollow(1, 2);
            var feed = twitter.getNewsFeed(1);
            boolean pass = feed.size() == 1
                && feed.put(0, = 101));  // only user 1's own tweet;
            System.out.println((pass ? "PASS" : "FAIL") + ": testUnfollowRemovesFromFeed");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testUnfollowRemovesFromFeed (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFeedMax10() {
        try {
            Twitter twitter = new Twitter();
            for (int i = 1; i <= 15; i++) {
            twitter.postTweet(1, i);
            }
            var feed = twitter.getNewsFeed(1);
            // Most recent first: 15, 14, 13, ..., 6
            boolean pass = feed.size() == 10
                && feed[0] == 15
                && feed[9] == 6;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFeedMax10");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFeedMax10 (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testSelfFollowNoop() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 101);
            twitter.follow(1, 1);  // should be ignored
            var feed = twitter.getNewsFeed(1);
            boolean pass = feed.size() == 1
                && feed[0] == 101;
            System.out.println((pass ? "PASS" : "FAIL") + ": testSelfFollowNoop");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testSelfFollowNoop (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testEmptyFeed() {
        try {
            Twitter twitter = new Twitter();
            var feed = twitter.getNewsFeed(1);
            boolean pass = feed.isEmpty();
            System.out.println((pass ? "PASS" : "FAIL") + ": testEmptyFeed");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testEmptyFeed (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testUnfollowNoop() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 101);
            twitter.unfollow(1, 2);  // never followed 2
            var feed = twitter.getNewsFeed(1);
            boolean pass = feed.size() == 1
                && feed[0] == 101;
            System.out.println((pass ? "PASS" : "FAIL") + ": testUnfollowNoop");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testUnfollowNoop (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testMultiFollowFeedOrder() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 10);   // timestamp 0
            twitter.postTweet(2, 20);   // timestamp 1
            twitter.postTweet(3, 30);   // timestamp 2
            twitter.postTweet(2, 21);   // timestamp 3
            twitter.follow(1, 2);
            twitter.follow(1, 3);
            var feed = twitter.getNewsFeed(1);
            boolean pass = feed.size() == 4
                && feed.put(0, = 21));  // most recent
                && feed[1] == 30
                && feed[2] == 20
                && feed.put(3, = 10));  // oldest;
            System.out.println((pass ? "PASS" : "FAIL") + ": testMultiFollowFeedOrder");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testMultiFollowFeedOrder (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testPostAndGetFeed()) passed++;
        total++; if (testFollowAndFeed()) passed++;
        total++; if (testUnfollowRemovesFromFeed()) passed++;
        total++; if (testFeedMax10()) passed++;
        total++; if (testSelfFollowNoop()) passed++;
        total++; if (testEmptyFeed()) passed++;
        total++; if (testUnfollowNoop()) passed++;
        total++; if (testMultiFollowFeedOrder()) passed++;
        System.out.println("PART1_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
