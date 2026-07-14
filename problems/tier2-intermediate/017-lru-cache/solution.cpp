#include <iostream>
#include <memory>
#include <unordered_map>
#include <vector>
#include <algorithm>
using namespace std;

// ─── Data Model ─────────────────────────────────────────────────────────────

struct Node {
    int key;
    int value;
    long expiresAt;
    Node* prev;
    Node* next;
    Node(int k, int v, long exp = 0)
        : key(k), value(v), expiresAt(exp), prev(nullptr), next(nullptr) {}
};

enum class EvictionReason { CAPACITY, TTL_EXPIRED, EXPLICIT_DELETE };

// ─── Eviction Listener Interface ─────────────────────────────────────────────

class EvictionListener {
public:
    virtual void onEviction(int key, int value, EvictionReason reason) = 0;
    virtual ~EvictionListener() = default;
};

// ─── LRU Cache ──────────────────────────────────────────────────────────────

class LRUCache {
private:
    int capacity;
    unordered_map<int, Node*> cache;
    Node* head;
    Node* tail;
    vector<EvictionListener*> listeners;
    bool notifying = false;

    void addToFront(Node* node) {
        node->next = head->next;
        node->prev = head;
        head->next->prev = node;
        head->next = node;
    }

    void removeNode(Node* node) {
        node->prev->next = node->next;
        node->next->prev = node->prev;
    }

    void moveToFront(Node* node) {
        removeNode(node);
        addToFront(node);
    }

    bool isExpired(Node* node, long currentTime) {
        return node->expiresAt > 0 && currentTime >= node->expiresAt;
    }

    void notifyListeners(int key, int value, EvictionReason reason) {
        notifying = true;
        for (auto* listener : listeners) {
            listener->onEviction(key, value, reason);
        }
        notifying = false;
    }

    void evictNode(Node* node, EvictionReason reason) {
        removeNode(node);
        cache.erase(node->key);
        notifyListeners(node->key, node->value, reason);
        delete node;
    }

    void evictLRU() {
        if (head->next == tail) return;  // empty
        Node* lru = tail->prev;
        evictNode(lru, EvictionReason::CAPACITY);
    }

public:
    LRUCache(int cap) : capacity(cap) {
        head = new Node(0, 0);
        tail = new Node(0, 0);
        head->next = tail;
        tail->prev = head;
    }

    ~LRUCache() {
        Node* curr = head;
        while (curr) {
            Node* next = curr->next;
            delete curr;
            curr = next;
        }
    }

    // Part 1 interface
    int get(int key) {
        if (cache.find(key) == cache.end()) return -1;
        Node* node = cache[key];
        moveToFront(node);
        return node->value;
    }

    void put(int key, int value) {
        if (cache.find(key) != cache.end()) {
            Node* node = cache[key];
            node->value = value;
            moveToFront(node);
            return;
        }
        if ((int)cache.size() >= capacity) {
            evictLRU();
        }
        Node* newNode = new Node(key, value);
        addToFront(newNode);
        cache[key] = newNode;
    }

    // Part 2 interface
    int get(int key, long currentTime) {
        if (cache.find(key) == cache.end()) return -1;
        Node* node = cache[key];
        if (isExpired(node, currentTime)) {
            evictNode(node, EvictionReason::TTL_EXPIRED);
            return -1;
        }
        moveToFront(node);
        return node->value;
    }

    void put(int key, int value, long currentTime, int ttl = 0) {
        long expiresAt = (ttl > 0) ? currentTime + ttl : 0;

        if (cache.find(key) != cache.end()) {
            Node* node = cache[key];
            if (isExpired(node, currentTime)) {
                evictNode(node, EvictionReason::TTL_EXPIRED);
            } else {
                node->value = value;
                node->expiresAt = expiresAt;
                moveToFront(node);
                return;
            }
        }
        if ((int)cache.size() >= capacity) {
            evictLRU();
        }
        Node* newNode = new Node(key, value, expiresAt);
        addToFront(newNode);
        cache[key] = newNode;
    }

    bool deleteKey(int key) {
        if (cache.find(key) == cache.end()) return false;
        evictNode(cache[key], EvictionReason::EXPLICIT_DELETE);
        return true;
    }

    int size() {
        return (int)cache.size();
    }

    // Part 3 interface
    void addEvictionListener(EvictionListener* listener) {
        listeners.push_back(listener);
    }

    void removeEvictionListener(EvictionListener* listener) {
        listeners.erase(
            remove(listeners.begin(), listeners.end(), listener),
            listeners.end());
    }
};

// ─── Ops simulator (used by spec-based tests) ───────────────────────────────
//
// LRUCache is a stateful "design" problem. The simulator drives a single
// LRUCache instance through a list of operations and returns one string per
// operation describing its outcome.
//
// Op fields (ints in i1..i4, string in s1):
//   "new"            — create cache(capacity=i1)              -> "ok"
//   "put"            — put(i1, i2)                            -> "ok"
//   "put_t"          — put(i1, i2, i3=now, i4=ttl)            -> "ok"
//   "get"            — get(i1)                                -> int as string
//   "get_t"          — get(i1, i2=now)                        -> int as string
//   "delete"         — deleteKey(i1)                          -> "ok"/"fail"
//   "size"           —                                         -> int as string
//   "add_listener"   — i1 = listener slot index                -> "ok"
//   "remove_listener"— i1 = listener slot index                -> "ok"
//   "events_count"   — i1 = listener slot index                -> int as string
//   "event_key"      — i1 = listener slot, i2 = event index    -> int as string
//   "event_value"    — i1 = listener slot, i2 = event index    -> int as string
//   "event_reason"   — i1 = listener slot, i2 = event index    -> "CAPACITY"|"TTL_EXPIRED"|"EXPLICIT_DELETE"

struct LruOp {
    string kind;
    int i1;
    int i2;
    int i3;
    int i4;
};

struct CapturedEvent { int key; int value; EvictionReason reason; };

class CapturingListener : public EvictionListener {
public:
    vector<CapturedEvent> events;
    void onEviction(int key, int value, EvictionReason r) override {
        events.push_back({key, value, r});
    }
};

static const char* reason_str(EvictionReason r) {
    switch (r) {
        case EvictionReason::CAPACITY:        return "CAPACITY";
        case EvictionReason::TTL_EXPIRED:     return "TTL_EXPIRED";
        case EvictionReason::EXPLICIT_DELETE: return "EXPLICIT_DELETE";
    }
    return "UNKNOWN";
}

vector<string> lru_simulate(vector<LruOp> ops) {
    vector<string> out;
    unique_ptr<LRUCache> cache;
    vector<unique_ptr<CapturingListener>> listeners;
    auto ensure_listener = [&](int idx) {
        while ((int)listeners.size() <= idx)
            listeners.push_back(unique_ptr<CapturingListener>(new CapturingListener()));
    };
    for (const auto& op : ops) {
        const string& k = op.kind;
        if (k == "new") {
            cache.reset(new LRUCache(op.i1));
            listeners.clear();
            out.push_back("ok");
        } else if (k == "put") {
            cache->put(op.i1, op.i2);
            out.push_back("ok");
        } else if (k == "put_t") {
            cache->put(op.i1, op.i2, (long)op.i3, op.i4);
            out.push_back("ok");
        } else if (k == "get") {
            out.push_back(to_string(cache->get(op.i1)));
        } else if (k == "get_t") {
            out.push_back(to_string(cache->get(op.i1, (long)op.i2)));
        } else if (k == "delete") {
            out.push_back(cache->deleteKey(op.i1) ? "ok" : "fail");
        } else if (k == "size") {
            out.push_back(to_string(cache->size()));
        } else if (k == "add_listener") {
            ensure_listener(op.i1);
            cache->addEvictionListener(listeners[op.i1].get());
            out.push_back("ok");
        } else if (k == "remove_listener") {
            if (op.i1 < (int)listeners.size())
                cache->removeEvictionListener(listeners[op.i1].get());
            out.push_back("ok");
        } else if (k == "events_count") {
            out.push_back(op.i1 < (int)listeners.size()
                          ? to_string((int)listeners[op.i1]->events.size())
                          : "0");
        } else if (k == "event_key") {
            if (op.i1 < (int)listeners.size() && op.i2 < (int)listeners[op.i1]->events.size())
                out.push_back(to_string(listeners[op.i1]->events[op.i2].key));
            else out.push_back("");
        } else if (k == "event_value") {
            if (op.i1 < (int)listeners.size() && op.i2 < (int)listeners[op.i1]->events.size())
                out.push_back(to_string(listeners[op.i1]->events[op.i2].value));
            else out.push_back("");
        } else if (k == "event_reason") {
            if (op.i1 < (int)listeners.size() && op.i2 < (int)listeners[op.i1]->events.size())
                out.push_back(reason_str(listeners[op.i1]->events[op.i2].reason));
            else out.push_back("");
        } else {
            out.push_back("unknown:" + k);
        }
    }
    return out;
}

// ─── Main ────────────────────────────────────────────────────────────────────

#ifndef RUNNING_TESTS
int main() {
    LRUCache cache(3);

    cache.put(1, 10);
    cache.put(2, 20);
    cache.put(3, 30);
    cout << "get(1) = " << cache.get(1) << endl;  // 10
    cache.put(4, 40);  // evicts 2
    cout << "get(2) = " << cache.get(2) << endl;  // -1 (evicted)
    cout << "get(3) = " << cache.get(3) << endl;  // 30
    cout << "get(4) = " << cache.get(4) << endl;  // 40

    return 0;
}
#endif
