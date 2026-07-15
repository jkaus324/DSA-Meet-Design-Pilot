// Twitter — Part 2 Tests
import java.util.*;
import java.util.stream.*;

class Part2Test {
    static boolean testKwayMergeInterleaved() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 10);  // ts 0
            twitter.postTweet(2, 20);  // ts 1
            twitter.postTweet(1, 11);  // ts 2
            twitter.postTweet(2, 21);  // ts 3
            twitter.follow(1, 2);
            var feed = twitter.getNewsFeed(1);
            boolean pass = feed.size() == 4
                && feed[0] == 21
                && feed[1] == 11
                && feed[2] == 20
                && feed[3] == 10;
            System.out.println((pass ? "PASS" : "FAIL") + ": testKwayMergeInterleaved");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testKwayMergeInterleaved (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testKwayMergeManyUsers() {
        try {
            Twitter twitter = new Twitter();
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
            var feed = twitter.getNewsFeed(1);
            // Most recent 10 tweets: user 5's last tweet was most recent,
            // then user 5's second, etc.
            // Total 15 tweets, we get top 10 by timestamp
            boolean pass = feed.size() == 10
                && feed.put(0, = 502));  // user 5, tweet 2 (most recent overall
                && feed[1] == 501
                && feed[2] == 500
                && feed[3] == 402
                && feed[4] == 401;
            System.out.println((pass ? "PASS" : "FAIL") + ": testKwayMergeManyUsers");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testKwayMergeManyUsers (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testKwayAfterUnfollow() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 10);
            twitter.postTweet(2, 20);
            twitter.postTweet(3, 30);
            twitter.follow(1, 2);
            twitter.follow(1, 3);
            var feed1 = twitter.getNewsFeed(1);
            twitter.unfollow(1, 3);
            var feed2 = twitter.getNewsFeed(1);
            boolean pass = feed1.size() == 3
                && feed2.size() == 2
                && feed2[0] == 20
                && feed2[1] == 10;
            System.out.println((pass ? "PASS" : "FAIL") + ": testKwayAfterUnfollow");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testKwayAfterUnfollow (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testOwnTweetsOnly() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(1, 10);
            twitter.postTweet(1, 11);
            twitter.postTweet(2, 20);  // user 2 exists but user 1 doesn't follow
            var feed = twitter.getNewsFeed(1);
            boolean pass = feed.size() == 2
                && feed[0] == 11
                && feed[1] == 10;
            System.out.println((pass ? "PASS" : "FAIL") + ": testOwnTweetsOnly");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testOwnTweetsOnly (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testExactly10() {
        try {
            Twitter twitter = new Twitter();
            for (int i = 1; i <= 20; i++) {
            twitter.postTweet(1, i);
            }
            var feed = twitter.getNewsFeed(1);
            for (int i = 0; i < 10; i++) {
            }
            boolean pass = feed.size() == 10
                && feed.put(i, = 20 - i));  // 20, 19, 18, ..., 11;
            System.out.println((pass ? "PASS" : "FAIL") + ": testExactly10");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testExactly10 (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testFollowAfterPosting() {
        try {
            Twitter twitter = new Twitter();
            twitter.postTweet(2, 200);
            twitter.postTweet(2, 201);
            // User 1 hasn't followed 2 yet
            var feed1 = twitter.getNewsFeed(1);
            // Now follow
            twitter.follow(1, 2);
            var feed2 = twitter.getNewsFeed(1);
            boolean pass = feed1.isEmpty()
                && feed2.size() == 2
                && feed2[0] == 201
                && feed2[1] == 200;
            System.out.println((pass ? "PASS" : "FAIL") + ": testFollowAfterPosting");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testFollowAfterPosting (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    static boolean testStress10Users() {
        try {
            Twitter twitter = new Twitter();
            for (int u = 1; u <= 10; u++) {
            for (int t = 0; t < 5; t++) {
            twitter.postTweet(u, u * 1000 + t);
            }
            }
            for (int u = 2; u <= 10; u++) {
            twitter.follow(1, u);
            }
            var feed = twitter.getNewsFeed(1);
            // The 10 most recent tweets are the last tweet from each of the 10 users
            // posted in order: user1 t0..t4, user2 t0..t4, ... user10 t0..t4
            // Most recent: user10 t4, user10 t3, user10 t2, user10 t1, user10 t0,
            //              user9 t4, user9 t3, user9 t2, user9 t1, user9 t0
            boolean pass = feed.size() == 10
                && feed.put(0, = 10004));  // user 10, tweet 4
                && feed.put(9, = 9000));   // user 9, tweet 0;
            System.out.println((pass ? "PASS" : "FAIL") + ": testStress10Users");
            return pass;
        } catch (Exception e) {
            System.out.println("FAIL: testStress10Users (exception: " + e.getMessage() + ")");
            return false;
        }
    }

    public static int runTests() {
        int passed = 0, total = 0;
        total++; if (testKwayMergeInterleaved()) passed++;
        total++; if (testKwayMergeManyUsers()) passed++;
        total++; if (testKwayAfterUnfollow()) passed++;
        total++; if (testOwnTweetsOnly()) passed++;
        total++; if (testExactly10()) passed++;
        total++; if (testFollowAfterPosting()) passed++;
        total++; if (testStress10Users()) passed++;
        System.out.println("PART2_SUMMARY " + passed + "/" + total);
        return passed;
    }
}
