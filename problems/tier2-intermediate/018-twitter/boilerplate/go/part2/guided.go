package main

import "container/heap"

// Tweet holds a tweet's ID and the global timestamp it was posted.
type Tweet struct {
	TweetId   int
	Timestamp int
}

// Twitter — optimized with k-way merge.
// HINT: The key optimization is in GetNewsFeed.
// HINT: Use a max-heap (container/heap) to merge K sorted tweet streams.
// HINT: Each heap entry needs: timestamp, tweetId, userId, and an index
//       pointing to the next tweet from that user.
type Twitter struct {
	// HINT: Same data structures as Part 1 for storage
	// HINT: The optimization is purely in how you READ, not how you WRITE
}

func NewTwitter() *Twitter {
	// TODO: Initialize the system
	return &Twitter{}
}

func (t *Twitter) PostTweet(userId, tweetId int) {
	// TODO: Same as Part 1 — store tweet with auto-incrementing timestamp
}

func (t *Twitter) GetNewsFeed(userId int) []int {
	// HINT: Build a set of all relevant users (self + following)
	// HINT: For each user with tweets, push their MOST RECENT tweet into a max-heap
	// HINT: Heap entry = {timestamp, tweetId, userId, indexInTweetList}
	// HINT: Pop the max, add to result, push that user's NEXT tweet (index-1)
	// HINT: Stop after 10 results or when heap is empty
	_ = heap.Init // use container/heap
	return nil
}

func (t *Twitter) Follow(followerId, followeeId int) {
	// TODO: Same as Part 1
}

func (t *Twitter) Unfollow(followerId, followeeId int) {
	// TODO: Same as Part 1
}
