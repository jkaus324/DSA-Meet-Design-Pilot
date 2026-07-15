// Twitter — post tweets, follow/unfollow, k-way merged news feed (Go port).
package main

import (
	"sort"
	"strconv"
)

type TwitterOp struct {
	kind string
	i1   int
	i2   int
}

type tweetRec struct {
	tid int
	ts  int
}

type twitter struct {
	time    int
	tweets  map[int][]tweetRec
	follows map[int]map[int]bool
}

func newTwitter() *twitter {
	return &twitter{
		time:    0,
		tweets:  map[int][]tweetRec{},
		follows: map[int]map[int]bool{},
	}
}

func (t *twitter) postTweet(userID, tweetID int) {
	t.tweets[userID] = append(t.tweets[userID], tweetRec{tweetID, t.time})
	t.time++
}

func (t *twitter) getNewsFeed(userID int) []int {
	users := map[int]bool{userID: true}
	for f := range t.follows[userID] {
		users[f] = true
	}
	// Collect the candidate tweets. Timestamps are globally unique (time is a
	// monotonically increasing counter), so a global sort by descending
	// timestamp reproduces the k-way merge ordering exactly.
	cand := []tweetRec{}
	for uid := range users {
		cand = append(cand, t.tweets[uid]...)
	}
	sort.Slice(cand, func(i, j int) bool { return cand[i].ts > cand[j].ts })
	result := []int{}
	for i := 0; i < len(cand) && len(result) < 10; i++ {
		result = append(result, cand[i].tid)
	}
	return result
}

func (t *twitter) follow(followerID, followeeID int) {
	if followerID == followeeID {
		return
	}
	if t.follows[followerID] == nil {
		t.follows[followerID] = map[int]bool{}
	}
	t.follows[followerID][followeeID] = true
}

func (t *twitter) unfollow(followerID, followeeID int) {
	if s := t.follows[followerID]; s != nil {
		delete(s, followeeID)
	}
}

func twitter_simulate(ops []TwitterOp) []string {
	out := []string{}
	tw := newTwitter()
	for _, op := range ops {
		switch op.kind {
		case "new":
			tw = newTwitter()
			out = append(out, "ok")
		case "post":
			tw.postTweet(op.i1, op.i2)
			out = append(out, "ok")
		case "follow":
			tw.follow(op.i1, op.i2)
			out = append(out, "ok")
		case "unfollow":
			tw.unfollow(op.i1, op.i2)
			out = append(out, "ok")
		case "feed_size":
			out = append(out, strconv.Itoa(len(tw.getNewsFeed(op.i1))))
		case "feed_at":
			f := tw.getNewsFeed(op.i1)
			if op.i2 >= 0 && op.i2 < len(f) {
				out = append(out, strconv.Itoa(f[op.i2]))
			} else {
				out = append(out, "-1")
			}
		default:
			out = append(out, "unknown:"+op.kind)
		}
	}
	return out
}
