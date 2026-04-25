package main

// Tweet holds a tweet's ID and the global timestamp it was posted.
type Tweet struct {
	TweetId   int
	Timestamp int
}

// Twitter manages users, their tweets, and follow relationships.
// HINT: You need two maps — one to store each user's tweets,
// and one to store who each user follows.
// HINT: Use a global counter to assign timestamps so tweets can be ordered.
type Twitter struct {
	// HINT: What data structure gives O(1) lookup by userId?
	// HINT: What data structure lets you quickly check "does user A follow user B"?
	// HINT: How do you ensure tweets are ordered by recency?
}

func NewTwitter() *Twitter {
	// HINT: Initialize the system
	return &Twitter{}
}

func (t *Twitter) PostTweet(userId, tweetId int) {
	// HINT: Auto-create the user's tweet list if it doesn't exist
	// HINT: Assign a timestamp so you can order tweets later
}

func (t *Twitter) GetNewsFeed(userId int) []int {
	// HINT: Collect tweets from this user and all users they follow
	// HINT: Gather all relevant tweets into one slice
	// HINT: Sort by timestamp (most recent first)
	// HINT: Return at most 10 tweet IDs
	return nil
}

func (t *Twitter) Follow(followerId, followeeId int) {
	// HINT: A user should not be able to follow themselves
}

func (t *Twitter) Unfollow(followerId, followeeId int) {
	// HINT: Unfollowing someone you don't follow should be a no-op
}
