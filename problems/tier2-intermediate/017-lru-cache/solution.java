// LRU Cache — Solution (Java)
import java.util.*;

class LruOp {
    public String kind;
    public int i1, i2, i3, i4;

    public LruOp(String kind, int i1, int i2, int i3, int i4) {
        this.kind = kind; this.i1 = i1; this.i2 = i2; this.i3 = i3; this.i4 = i4;
    }
}

enum EvictionReason { CAPACITY, TTL_EXPIRED, EXPLICIT_DELETE }

class CapturedEvent {
    public int key;
    public int value;
    public EvictionReason reason;
    public CapturedEvent(int k, int v, EvictionReason r) { key = k; value = v; reason = r; }
}

interface EvictionListener {
    void onEviction(int key, int value, EvictionReason reason);
}

class LruNode {
    int key, value;
    long expiresAt;
    LruNode prev, next;
    LruNode(int k, int v, long exp) { key = k; value = v; expiresAt = exp; }
}

class LRUCache {
    int capacity;
    Map<Integer, LruNode> cache = new HashMap<>();
    LruNode head, tail;
    List<EvictionListener> listeners = new ArrayList<>();

    public LRUCache(int cap) {
        this.capacity = cap;
        head = new LruNode(0, 0, 0);
        tail = new LruNode(0, 0, 0);
        head.next = tail;
        tail.prev = head;
    }

    private void addToFront(LruNode node) {
        node.next = head.next;
        node.prev = head;
        head.next.prev = node;
        head.next = node;
    }

    private void removeNode(LruNode node) {
        node.prev.next = node.next;
        node.next.prev = node.prev;
    }

    private void moveToFront(LruNode node) { removeNode(node); addToFront(node); }

    private boolean isExpired(LruNode node, long currentTime) {
        return node.expiresAt > 0 && currentTime >= node.expiresAt;
    }

    private void notifyListeners(int key, int value, EvictionReason reason) {
        // Snapshot to avoid concurrent modification
        List<EvictionListener> copy = new ArrayList<>(listeners);
        for (EvictionListener l : copy) l.onEviction(key, value, reason);
    }

    private void evictNode(LruNode node, EvictionReason reason) {
        removeNode(node);
        cache.remove(node.key);
        notifyListeners(node.key, node.value, reason);
    }

    private void evictLRU() {
        if (head.next == tail) return;
        LruNode lru = tail.prev;
        evictNode(lru, EvictionReason.CAPACITY);
    }

    public int get(int key) {
        LruNode node = cache.get(key);
        if (node == null) return -1;
        moveToFront(node);
        return node.value;
    }

    public void put(int key, int value) {
        LruNode existing = cache.get(key);
        if (existing != null) {
            existing.value = value;
            moveToFront(existing);
            return;
        }
        if (cache.size() >= capacity) evictLRU();
        LruNode newNode = new LruNode(key, value, 0);
        addToFront(newNode);
        cache.put(key, newNode);
    }

    public int getT(int key, long currentTime) {
        LruNode node = cache.get(key);
        if (node == null) return -1;
        if (isExpired(node, currentTime)) {
            evictNode(node, EvictionReason.TTL_EXPIRED);
            return -1;
        }
        moveToFront(node);
        return node.value;
    }

    public void putT(int key, int value, long currentTime, int ttl) {
        long expiresAt = (ttl > 0) ? currentTime + ttl : 0;
        LruNode existing = cache.get(key);
        if (existing != null) {
            if (isExpired(existing, currentTime)) {
                evictNode(existing, EvictionReason.TTL_EXPIRED);
            } else {
                existing.value = value;
                existing.expiresAt = expiresAt;
                moveToFront(existing);
                return;
            }
        }
        if (cache.size() >= capacity) evictLRU();
        LruNode newNode = new LruNode(key, value, expiresAt);
        addToFront(newNode);
        cache.put(key, newNode);
    }

    public boolean deleteKey(int key) {
        LruNode node = cache.get(key);
        if (node == null) return false;
        evictNode(node, EvictionReason.EXPLICIT_DELETE);
        return true;
    }

    public int size() { return cache.size(); }

    public void addEvictionListener(EvictionListener listener) { listeners.add(listener); }

    public void removeEvictionListener(EvictionListener listener) {
        listeners.remove(listener);
    }
}

class CapturingListener implements EvictionListener {
    public List<CapturedEvent> events = new ArrayList<>();
    @Override
    public void onEviction(int key, int value, EvictionReason reason) {
        events.add(new CapturedEvent(key, value, reason));
    }
}

public class Solution {
    public static List<String> lru_simulate(List<LruOp> ops) {
        List<String> out = new ArrayList<>();
        LRUCache cache = null;
        List<CapturingListener> listeners = new ArrayList<>();

        for (LruOp op : ops) {
            String k = op.kind;
            if ("new".equals(k)) {
                cache = new LRUCache(op.i1);
                listeners = new ArrayList<>();
                out.add("ok");
            } else if ("put".equals(k)) {
                cache.put(op.i1, op.i2);
                out.add("ok");
            } else if ("put_t".equals(k)) {
                cache.putT(op.i1, op.i2, (long) op.i3, op.i4);
                out.add("ok");
            } else if ("get".equals(k)) {
                out.add(Integer.toString(cache.get(op.i1)));
            } else if ("get_t".equals(k)) {
                out.add(Integer.toString(cache.getT(op.i1, (long) op.i2)));
            } else if ("delete".equals(k)) {
                out.add(cache.deleteKey(op.i1) ? "ok" : "fail");
            } else if ("size".equals(k)) {
                out.add(Integer.toString(cache.size()));
            } else if ("add_listener".equals(k)) {
                while (listeners.size() <= op.i1) listeners.add(new CapturingListener());
                cache.addEvictionListener(listeners.get(op.i1));
                out.add("ok");
            } else if ("remove_listener".equals(k)) {
                if (op.i1 < listeners.size()) cache.removeEvictionListener(listeners.get(op.i1));
                out.add("ok");
            } else if ("events_count".equals(k)) {
                out.add(op.i1 < listeners.size() ? Integer.toString(listeners.get(op.i1).events.size()) : "0");
            } else if ("event_key".equals(k)) {
                if (op.i1 < listeners.size() && op.i2 < listeners.get(op.i1).events.size())
                    out.add(Integer.toString(listeners.get(op.i1).events.get(op.i2).key));
                else out.add("");
            } else if ("event_value".equals(k)) {
                if (op.i1 < listeners.size() && op.i2 < listeners.get(op.i1).events.size())
                    out.add(Integer.toString(listeners.get(op.i1).events.get(op.i2).value));
                else out.add("");
            } else if ("event_reason".equals(k)) {
                if (op.i1 < listeners.size() && op.i2 < listeners.get(op.i1).events.size())
                    out.add(listeners.get(op.i1).events.get(op.i2).reason.name());
                else out.add("");
            } else {
                out.add("unknown:" + k);
            }
        }
        return out;
    }
}
