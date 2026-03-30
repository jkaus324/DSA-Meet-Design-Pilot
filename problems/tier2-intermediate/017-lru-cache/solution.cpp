#include <iostream>
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
