package main

import "sort"

// Tweet holds a tweet's ID and the global timestamp it was posted.
type Tweet struct {
	TweetId   int
	Timestamp int
}

// Twitter manages users, their tweets, and follow relationships.
type Twitter struct {
	time    int
	tweets  map[int][]Tweet        // userId -> list of tweets
	follows map[int]map[int]bool   // userId -> set of followee IDs
}

func NewTwitter() *Twitter {
	return &Twitter{
		time:    0,
		tweets:  make(map[int][]Tweet),
		follows: make(map[int]map[int]bool),
	}
}

func (t *Twitter) PostTweet(userId, tweetId int) {
	// TODO: Append Tweet{TweetId: tweetId, Timestamp: t.time} to t.tweets[userId]
	// TODO: Increment t.time after appending
}

func (t *Twitter) GetNewsFeed(userId int) []int {
	// TODO: Collect all tweets from this user
	// TODO: Collect all tweets from users in t.follows[userId]
	// TODO: Sort collected tweets by Timestamp descending (most recent first)
	// TODO: Return the TweetIds of the first 10 (or fewer) tweets
	_ = sort.Slice // use sort.Slice for sorting
	return nil
}

func (t *Twitter) Follow(followerId, followeeId int) {
	// TODO: If followerId != followeeId, add followeeId to t.follows[followerId]
	// TODO: Initialize the inner map if it doesn't exist yet
}

func (t *Twitter) Unfollow(followerId, followeeId int) {
	// TODO: Delete followeeId from t.follows[followerId] (safe even if not present)
}
