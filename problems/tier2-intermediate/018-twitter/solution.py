"""Twitter — post tweets, follow/unfollow, k-way merged news feed."""

import heapq


class TwitterOp:
    def __init__(self, kind, i1=0, i2=0):
        self.kind = kind
        self.i1 = i1
        self.i2 = i2


class Twitter:
    def __init__(self):
        self.time = 0
        self.tweets = {}   # userId -> list[(tweetId, timestamp)]
        self.follows = {}  # userId -> set of followee userIds

    def post_tweet(self, user_id, tweet_id):
        self.tweets.setdefault(user_id, []).append((tweet_id, self.time))
        self.time += 1

    def get_news_feed(self, user_id):
        users = {user_id}
        if user_id in self.follows:
            users.update(self.follows[user_id])
        # Max-heap by timestamp via negated timestamp.
        heap = []
        for uid in users:
            user_tweets = self.tweets.get(uid)
            if user_tweets:
                idx = len(user_tweets) - 1
                tid, ts = user_tweets[idx]
                heapq.heappush(heap, (-ts, tid, uid, idx))
        result = []
        while heap and len(result) < 10:
            neg_ts, tid, uid, idx = heapq.heappop(heap)
            result.append(tid)
            if idx > 0:
                next_idx = idx - 1
                ntid, nts = self.tweets[uid][next_idx]
                heapq.heappush(heap, (-nts, ntid, uid, next_idx))
        return result

    def follow(self, follower_id, followee_id):
        if follower_id != followee_id:
            self.follows.setdefault(follower_id, set()).add(followee_id)

    def unfollow(self, follower_id, followee_id):
        if follower_id in self.follows:
            self.follows[follower_id].discard(followee_id)


def twitter_simulate(ops):
    out = []
    tw = Twitter()
    for op in ops:
        k = op.kind
        if k == "new":
            tw = Twitter()
            out.append("ok")
        elif k == "post":
            tw.post_tweet(op.i1, op.i2)
            out.append("ok")
        elif k == "follow":
            tw.follow(op.i1, op.i2)
            out.append("ok")
        elif k == "unfollow":
            tw.unfollow(op.i1, op.i2)
            out.append("ok")
        elif k == "feed_size":
            out.append(str(len(tw.get_news_feed(op.i1))))
        elif k == "feed_at":
            f = tw.get_news_feed(op.i1)
            if 0 <= op.i2 < len(f):
                out.append(str(f[op.i2]))
            else:
                out.append("-1")
        else:
            out.append("unknown:" + k)
    return out
