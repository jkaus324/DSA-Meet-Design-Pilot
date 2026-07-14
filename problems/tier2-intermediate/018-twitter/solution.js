"use strict";
// Twitter — post tweets, follow/unfollow, k-way merged news feed.

class TwitterOp {
  constructor(kind, i1 = 0, i2 = 0) {
    this.kind = kind;
    this.i1 = i1;
    this.i2 = i2;
  }
}

// Min-heap whose elements are tuples [negTs, tid, uid, idx], compared
// lexicographically exactly like Python's heapq on tuples.
function tupleLess(a, b) {
  for (let i = 0; i < a.length; i++) {
    if (a[i] < b[i]) return true;
    if (a[i] > b[i]) return false;
  }
  return false;
}

class MinHeap {
  constructor() {
    this.h = [];
  }
  get size() {
    return this.h.length;
  }
  push(item) {
    const h = this.h;
    h.push(item);
    let i = h.length - 1;
    while (i > 0) {
      const parent = (i - 1) >> 1;
      if (tupleLess(h[i], h[parent])) {
        [h[i], h[parent]] = [h[parent], h[i]];
        i = parent;
      } else break;
    }
  }
  pop() {
    const h = this.h;
    const top = h[0];
    const last = h.pop();
    if (h.length > 0) {
      h[0] = last;
      let i = 0;
      const n = h.length;
      while (true) {
        const l = 2 * i + 1;
        const r = 2 * i + 2;
        let smallest = i;
        if (l < n && tupleLess(h[l], h[smallest])) smallest = l;
        if (r < n && tupleLess(h[r], h[smallest])) smallest = r;
        if (smallest === i) break;
        [h[i], h[smallest]] = [h[smallest], h[i]];
        i = smallest;
      }
    }
    return top;
  }
}

class Twitter {
  constructor() {
    this.time = 0;
    this.tweets = new Map(); // userId -> list of [tweetId, timestamp]
    this.follows = new Map(); // userId -> Set of followee userIds
  }

  post_tweet(userId, tweetId) {
    if (!this.tweets.has(userId)) this.tweets.set(userId, []);
    this.tweets.get(userId).push([tweetId, this.time]);
    this.time += 1;
  }

  get_news_feed(userId) {
    const users = new Set([userId]);
    if (this.follows.has(userId)) {
      for (const u of this.follows.get(userId)) users.add(u);
    }
    const heap = new MinHeap();
    for (const uid of users) {
      const userTweets = this.tweets.get(uid);
      if (userTweets && userTweets.length > 0) {
        const idx = userTweets.length - 1;
        const [tid, ts] = userTweets[idx];
        heap.push([-ts, tid, uid, idx]);
      }
    }
    const result = [];
    while (heap.size > 0 && result.length < 10) {
      const [, tid, uid, idx] = heap.pop();
      result.push(tid);
      if (idx > 0) {
        const nextIdx = idx - 1;
        const [ntid, nts] = this.tweets.get(uid)[nextIdx];
        heap.push([-nts, ntid, uid, nextIdx]);
      }
    }
    return result;
  }

  follow(followerId, followeeId) {
    if (followerId !== followeeId) {
      if (!this.follows.has(followerId)) this.follows.set(followerId, new Set());
      this.follows.get(followerId).add(followeeId);
    }
  }

  unfollow(followerId, followeeId) {
    if (this.follows.has(followerId)) {
      this.follows.get(followerId).delete(followeeId);
    }
  }
}

function twitter_simulate(ops) {
  const out = [];
  let tw = new Twitter();
  for (const op of ops) {
    const k = op.kind;
    if (k === "new") {
      tw = new Twitter();
      out.push("ok");
    } else if (k === "post") {
      tw.post_tweet(op.i1, op.i2);
      out.push("ok");
    } else if (k === "follow") {
      tw.follow(op.i1, op.i2);
      out.push("ok");
    } else if (k === "unfollow") {
      tw.unfollow(op.i1, op.i2);
      out.push("ok");
    } else if (k === "feed_size") {
      out.push(String(tw.get_news_feed(op.i1).length));
    } else if (k === "feed_at") {
      const f = tw.get_news_feed(op.i1);
      if (op.i2 >= 0 && op.i2 < f.length) {
        out.push(String(f[op.i2]));
      } else {
        out.push("-1");
      }
    } else {
      out.push("unknown:" + k);
    }
  }
  return out;
}

module.exports = { TwitterOp, twitter_simulate };
