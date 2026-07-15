// Twitter — Solution (Java)
import java.util.*;

class TwitterOp {
    public String kind;
    public int i1, i2;
    public TwitterOp(String kind, int i1, int i2) { this.kind = kind; this.i1 = i1; this.i2 = i2; }
}

class Tweet {
    public int tweetId;
    public int timestamp;
    public Tweet(int t, int ts) { tweetId = t; timestamp = ts; }
}

class Twitter {
    int time = 0;
    Map<Integer, List<Tweet>> tweets = new HashMap<>();
    Map<Integer, Set<Integer>> follows = new HashMap<>();

    public void postTweet(int userId, int tweetId) {
        tweets.computeIfAbsent(userId, k -> new ArrayList<>()).add(new Tweet(tweetId, time++));
    }

    public List<Integer> getNewsFeed(int userId) {
        Set<Integer> users = new HashSet<>();
        users.add(userId);
        Set<Integer> f = follows.get(userId);
        if (f != null) users.addAll(f);

        // Heap entries: int[]{timestamp, tweetId, userId, index}
        PriorityQueue<int[]> pq = new PriorityQueue<>((a, b) -> Integer.compare(b[0], a[0]));
        for (int uid : users) {
            List<Tweet> ts = tweets.get(uid);
            if (ts != null && !ts.isEmpty()) {
                int idx = ts.size() - 1;
                pq.offer(new int[]{ts.get(idx).timestamp, ts.get(idx).tweetId, uid, idx});
            }
        }

        List<Integer> result = new ArrayList<>();
        while (!pq.isEmpty() && result.size() < 10) {
            int[] top = pq.poll();
            int idx = top[3];
            int uid = top[2];
            result.add(top[1]);
            if (idx > 0) {
                int next = idx - 1;
                Tweet nt = tweets.get(uid).get(next);
                pq.offer(new int[]{nt.timestamp, nt.tweetId, uid, next});
            }
        }
        return result;
    }

    public void follow(int followerId, int followeeId) {
        if (followerId != followeeId) {
            follows.computeIfAbsent(followerId, k -> new HashSet<>()).add(followeeId);
        }
    }

    public void unfollow(int followerId, int followeeId) {
        Set<Integer> set = follows.get(followerId);
        if (set != null) set.remove(followeeId);
    }
}

public class Solution {
    public static List<String> twitter_simulate(List<TwitterOp> ops) {
        List<String> out = new ArrayList<>();
        Twitter tw = new Twitter();
        for (TwitterOp op : ops) {
            String k = op.kind;
            if ("new".equals(k)) {
                tw = new Twitter();
                out.add("ok");
            } else if ("post".equals(k)) {
                tw.postTweet(op.i1, op.i2);
                out.add("ok");
            } else if ("follow".equals(k)) {
                tw.follow(op.i1, op.i2);
                out.add("ok");
            } else if ("unfollow".equals(k)) {
                tw.unfollow(op.i1, op.i2);
                out.add("ok");
            } else if ("feed_size".equals(k)) {
                out.add(Integer.toString(tw.getNewsFeed(op.i1).size()));
            } else if ("feed_at".equals(k)) {
                List<Integer> f = tw.getNewsFeed(op.i1);
                out.add(op.i2 >= 0 && op.i2 < f.size() ? Integer.toString(f.get(op.i2)) : "-1");
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}
