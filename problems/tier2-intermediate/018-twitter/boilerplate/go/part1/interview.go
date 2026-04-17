package main

// Tweet holds a tweet's ID and the global timestamp it was posted.
type Tweet struct {
	TweetId   int
	Timestamp int
}

// Twitter — design and implement this struct so that:
//   1. Users can post tweets (PostTweet)
//   2. Users can follow/unfollow others (Follow, Unfollow)
//   3. A user's news feed returns the 10 most recent tweet IDs
//      from themselves and everyone they follow (GetNewsFeed)
//
// Think about:
//   - How do you store user-to-tweets and user-to-following relationships?
//   - How do you order tweets by recency?
//   - What happens when a user unfollows someone?
//   - Can a user follow themselves?
//
// Entry points (must exist for tests):
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
