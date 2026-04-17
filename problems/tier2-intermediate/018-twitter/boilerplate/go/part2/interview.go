package main

// Tweet holds a tweet's ID and the global timestamp it was posted.
type Tweet struct {
	TweetId   int
	Timestamp int
}

// Twitter — optimize the news feed using a k-way merge algorithm.
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
//   NewTwitter() *Twitter
//   (*Twitter).PostTweet(userId, tweetId int)
//   (*Twitter).GetNewsFeed(userId int) []int
//   (*Twitter).Follow(followerId, followeeId int)
//   (*Twitter).Unfollow(followerId, followeeId int)

type Twitter struct {
}

func NewTwitter() *Twitter {
	return &Twitter{}
}

func (t *Twitter) PostTweet(userId, tweetId int) {
}

func (t *Twitter) GetNewsFeed(userId int) []int {
	return nil
}

func (t *Twitter) Follow(followerId, followeeId int) {
}

func (t *Twitter) Unfollow(followerId, followeeId int) {
}
