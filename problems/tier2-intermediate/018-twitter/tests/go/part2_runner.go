package main

import "fmt"

func part2Tests() int {
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

	// Test 1: Basic k-way merge — two users interleaved
	test("test_kway_merge_interleaved", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 10) // ts 0
		tw.PostTweet(2, 20) // ts 1
		tw.PostTweet(1, 11) // ts 2
		tw.PostTweet(2, 21) // ts 3
		tw.Follow(1, 2)
		feed := tw.GetNewsFeed(1)
		if len(feed) != 4 {
			panic("expected 4 tweets")
		}
		if feed[0] != 21 {
			panic("expected 21 first")
		}
		if feed[1] != 11 {
			panic("expected 11 second")
		}
		if feed[2] != 20 {
			panic("expected 20 third")
		}
		if feed[3] != 10 {
			panic("expected 10 fourth")
		}
	})

	// Test 2: K-way merge with many users — top 10 from 5 users
	test("test_kway_merge_many_users", func() {
		tw := NewTwitter()
		for u := 1; u <= 5; u++ {
			for t := 0; t < 3; t++ {
				tw.PostTweet(u, u*100+t)
			}
		}
		for u := 2; u <= 5; u++ {
			tw.Follow(1, u)
		}
		feed := tw.GetNewsFeed(1)
		if len(feed) != 10 {
			panic("expected 10 tweets")
		}
		if feed[0] != 502 {
			panic("expected 502 first")
		}
		if feed[1] != 501 {
			panic("expected 501 second")
		}
		if feed[2] != 500 {
			panic("expected 500 third")
		}
		if feed[3] != 402 {
			panic("expected 402 fourth")
		}
		if feed[4] != 401 {
			panic("expected 401 fifth")
		}
	})

	// Test 3: Feed after follow then unfollow
	test("test_kway_after_unfollow", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 10)
		tw.PostTweet(2, 20)
		tw.PostTweet(3, 30)
		tw.Follow(1, 2)
		tw.Follow(1, 3)
		feed1 := tw.GetNewsFeed(1)
		if len(feed1) != 3 {
			panic("expected 3 tweets before unfollow")
		}
		tw.Unfollow(1, 3)
		feed2 := tw.GetNewsFeed(1)
		if len(feed2) != 2 {
			panic("expected 2 tweets after unfollow")
		}
		if feed2[0] != 20 {
			panic("expected 20 first")
		}
		if feed2[1] != 10 {
			panic("expected 10 second")
		}
	})

	// Test 4: Only user's own tweets when following nobody
	test("test_own_tweets_only", func() {
		tw := NewTwitter()
		tw.PostTweet(1, 10)
		tw.PostTweet(1, 11)
		tw.PostTweet(2, 20) // user 2 exists but user 1 doesn't follow
		feed := tw.GetNewsFeed(1)
		if len(feed) != 2 {
			panic("expected 2 tweets")
		}
		if feed[0] != 11 {
			panic("expected 11 first")
		}
		if feed[1] != 10 {
			panic("expected 10 second")
		}
	})

	// Test 5: Exactly 10 tweets when more are available
	test("test_exactly_10", func() {
		tw := NewTwitter()
		for i := 1; i <= 20; i++ {
			tw.PostTweet(1, i)
		}
		feed := tw.GetNewsFeed(1)
		if len(feed) != 10 {
			panic("expected exactly 10 tweets")
		}
		for i := 0; i < 10; i++ {
			if feed[i] != 20-i {
				panic(fmt.Sprintf("expected %d at index %d", 20-i, i))
			}
		}
	})

	// Test 6: Follow after posting — new tweets appear in feed
	test("test_follow_after_posting", func() {
		tw := NewTwitter()
		tw.PostTweet(2, 200)
		tw.PostTweet(2, 201)
		feed1 := tw.GetNewsFeed(1)
		if len(feed1) != 0 {
			panic("expected empty feed before follow")
		}
		tw.Follow(1, 2)
		feed2 := tw.GetNewsFeed(1)
		if len(feed2) != 2 {
			panic("expected 2 tweets after follow")
		}
		if feed2[0] != 201 {
			panic("expected 201 first")
		}
		if feed2[1] != 200 {
			panic("expected 200 second")
		}
	})

	// Test 7: Stress — merge from 10 users, 5 tweets each, get top 10
	test("test_stress_10_users", func() {
		tw := NewTwitter()
		for u := 1; u <= 10; u++ {
			for t := 0; t < 5; t++ {
				tw.PostTweet(u, u*1000+t)
			}
		}
		for u := 2; u <= 10; u++ {
			tw.Follow(1, u)
		}
		feed := tw.GetNewsFeed(1)
		if len(feed) != 10 {
			panic("expected 10 tweets")
		}
		if feed[0] != 10004 {
			panic("expected 10004 first")
		}
		if feed[9] != 9000 {
			panic("expected 9000 tenth")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
