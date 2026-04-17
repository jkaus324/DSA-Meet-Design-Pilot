package main

import "container/heap"

// Tweet holds a tweet's ID and the global timestamp it was posted.
type Tweet struct {
	TweetId   int
	Timestamp int
}

// heapEntry is one element pushed into the max-heap during k-way merge.
type heapEntry struct {
	timestamp int
	tweetId   int
	userId    int
	index     int // index into that user's tweet slice (pointing to the next tweet to consider)
}

// tweetHeap implements heap.Interface for a max-heap ordered by timestamp.
type tweetHeap []heapEntry

func (h tweetHeap) Len() int            { return len(h) }
func (h tweetHeap) Less(i, j int) bool  { return h[i].timestamp > h[j].timestamp } // max-heap
func (h tweetHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *tweetHeap) Push(x interface{}) { *h = append(*h, x.(heapEntry)) }
func (h *tweetHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Twitter manages users, their tweets, and follow relationships.
type Twitter struct {
	time    int
	tweets  map[int][]Tweet      // userId -> list of tweets (oldest first)
	follows map[int]map[int]bool // userId -> set of followee IDs
}

func NewTwitter() *Twitter {
	return &Twitter{
		tweets:  make(map[int][]Tweet),
		follows: make(map[int]map[int]bool),
	}
}

func (t *Twitter) PostTweet(userId, tweetId int) {
	// TODO: Append Tweet{TweetId: tweetId, Timestamp: t.time} to t.tweets[userId]
	// TODO: Increment t.time
}

func (t *Twitter) GetNewsFeed(userId int) []int {
	// Step 1: Build set of relevant users (self + following)
	// TODO: Create a map/set containing userId and all users in t.follows[userId]

	// Step 2: Seed the max-heap with the most recent tweet of each relevant user
	h := &tweetHeap{}
	heap.Init(h)
	// TODO: For each relevant user, if they have tweets, push their last tweet as
	//       heapEntry{timestamp, tweetId, userId, index = len-1}

	// Step 3: K-way merge
	result := []int{}
	for h.Len() > 0 && len(result) < 10 {
		// TODO: Pop the top entry
		// TODO: Append its tweetId to result
		// TODO: If that user has an older tweet (entry.index > 0), push entry with index-1
	}
	return result
}

func (t *Twitter) Follow(followerId, followeeId int) {
	// TODO: If followerId != followeeId, add followeeId to t.follows[followerId]
	// TODO: Initialize the inner map if it doesn't exist
}

func (t *Twitter) Unfollow(followerId, followeeId int) {
	// TODO: Delete followeeId from t.follows[followerId]
}
