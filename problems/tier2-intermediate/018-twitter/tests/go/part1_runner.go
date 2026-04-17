package main

import "fmt"

func part1Tests() int {
	passed := 0
	failed := 0

	test := func(name string, fn func()) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("FAIL", name)
					failed++
				}
			}()
			fn()
			fmt.Println("PASS", name)
			passed++
		}()
	}

	// Test 1: Post a tweet and retrieve it in the feed
	test("test_post_and_get_feed", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 101)
		feed := tw.GetNewsFeed(1)
		if len(feed) != 1 {
			panic("expected 1 tweet in feed")
		}
		if feed[0] != 101 {
			panic("expected tweet 101")
		}
	})

	// Test 2: Follow a user and see their tweets in the feed
	test("test_follow_and_feed", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 101)
		tw.PostTweet(2, 201)
		tw.Follow(1, 2)
		feed := tw.GetNewsFeed(1)
		if len(feed) != 2 {
			panic("expected 2 tweets in feed")
		}
		// Most recent first: tweet 201 was posted after 101
		if feed[0] != 201 {
			panic("expected 201 first")
		}
		if feed[1] != 101 {
			panic("expected 101 second")
		}
	})

	// Test 3: Unfollow removes tweets from feed
	test("test_unfollow_removes_from_feed", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 101)
		tw.PostTweet(2, 201)
		tw.Follow(1, 2)
		tw.Unfollow(1, 2)
		feed := tw.GetNewsFeed(1)
		if len(feed) != 1 {
			panic("expected 1 tweet after unfollow")
		}
		if feed[0] != 101 {
			panic("expected only user 1's own tweet")
		}
	})

	// Test 4: Feed returns at most 10 tweets
	test("test_feed_max_10", func() {
		tw := NewTwitter()
		for i := 1; i <= 15; i++ {
			tw.PostTweet(1, i)
		}
		feed := tw.GetNewsFeed(1)
		if len(feed) != 10 {
			panic("expected exactly 10 tweets")
		}
		if feed[0] != 15 {
			panic("expected most recent first (15)")
		}
		if feed[9] != 6 {
			panic("expected 6 as 10th tweet")
		}
	})

	// Test 5: Self-follow is a no-op
	test("test_self_follow_noop", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 101)
		tw.Follow(1, 1) // should be ignored
		feed := tw.GetNewsFeed(1)
		if len(feed) != 1 {
			panic("expected 1 tweet after self-follow")
		}
		if feed[0] != 101 {
			panic("expected tweet 101")
		}
	})

	// Test 6: Empty feed for user with no tweets and no follows
	test("test_empty_feed", func() {
		tw := NewTwitter()
		feed := tw.GetNewsFeed(1)
		if len(feed) != 0 {
			panic("expected empty feed")
		}
	})

	// Test 7: Unfollow someone not followed is a no-op
	test("test_unfollow_noop", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 101)
		tw.Unfollow(1, 2) // never followed 2
		feed := tw.GetNewsFeed(1)
		if len(feed) != 1 {
			panic("expected 1 tweet")
		}
		if feed[0] != 101 {
			panic("expected tweet 101")
		}
	})

	// Test 8: Feed merges tweets from multiple followed users in order
	test("test_multi_follow_feed_order", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 10) // timestamp 0
		tw.PostTweet(2, 20) // timestamp 1
		tw.PostTweet(3, 30) // timestamp 2
		tw.PostTweet(2, 21) // timestamp 3
		tw.Follow(1, 2)
		tw.Follow(1, 3)
		feed := tw.GetNewsFeed(1)
		if len(feed) != 4 {
			panic("expected 4 tweets")
		}
		if feed[0] != 21 {
			panic("expected 21 most recent")
		}
		if feed[1] != 30 {
			panic("expected 30 second")
		}
		if feed[2] != 20 {
			panic("expected 20 third")
		}
		if feed[3] != 10 {
			panic("expected 10 oldest")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
