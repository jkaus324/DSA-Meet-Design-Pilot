package main

import "container/heap"

// ─── Data Model ──────────────────────────────────────────────────────────────

type Tweet struct {
	TweetId   int
	Timestamp int
}

// heapEntry is one element pushed into the max-heap during k-way merge.
type heapEntry struct {
	timestamp int
	tweetId   int
	userId    int
	index     int
}

type tweetHeap []heapEntry

func (h tweetHeap) Len() int            { return len(h) }
func (h tweetHeap) Less(i, j int) bool  { return h[i].timestamp > h[j].timestamp }
func (h tweetHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *tweetHeap) Push(x interface{}) { *h = append(*h, x.(heapEntry)) }
func (h *tweetHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// ─── Twitter ─────────────────────────────────────────────────────────────────

type Twitter struct {
	time    int
	tweets  map[int][]Tweet
	follows map[int]map[int]bool
}

func NewTwitter() *Twitter {
	return &Twitter{
		tweets:  make(map[int][]Tweet),
		follows: make(map[int]map[int]bool),
	}
}

func (t *Twitter) PostTweet(userId, tweetId int) {
	t.tweets[userId] = append(t.tweets[userId], Tweet{TweetId: tweetId, Timestamp: t.time})
	t.time++
}

func (t *Twitter) GetNewsFeed(userId int) []int {
	relevant := map[int]bool{userId: true}
	for f := range t.follows[userId] {
		relevant[f] = true
	}

	h := &tweetHeap{}
	heap.Init(h)
	for uid := range relevant {
		feed, ok := t.tweets[uid]
		if !ok || len(feed) == 0 {
			continue
		}
		idx := len(feed) - 1
		heap.Push(h, heapEntry{
			timestamp: feed[idx].Timestamp,
			tweetId:   feed[idx].TweetId,
			userId:    uid,
			index:     idx,
		})
	}

	result := []int{}
	for h.Len() > 0 && len(result) < 10 {
		top := heap.Pop(h).(heapEntry)
		result = append(result, top.tweetId)
		if top.index > 0 {
			next := top.index - 1
			feed := t.tweets[top.userId]
			heap.Push(h, heapEntry{
				timestamp: feed[next].Timestamp,
				tweetId:   feed[next].TweetId,
				userId:    top.userId,
				index:     next,
			})
		}
	}
	return result
}

func (t *Twitter) Follow(followerId, followeeId int) {
	if followerId == followeeId {
		return
	}
	if t.follows[followerId] == nil {
		t.follows[followerId] = make(map[int]bool)
	}
	t.follows[followerId][followeeId] = true
}

func (t *Twitter) Unfollow(followerId, followeeId int) {
	if t.follows[followerId] != nil {
		delete(t.follows[followerId], followeeId)
	}
}
