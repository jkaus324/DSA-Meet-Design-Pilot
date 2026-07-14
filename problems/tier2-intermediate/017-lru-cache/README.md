# Problem 017 — LRU Cache

**Tier:** 2 (Intermediate) | **Pattern:** Strategy + Observer | **DSA:** HashMap + Doubly Linked List
**Companies:** Amazon, Google, Microsoft, Flipkart, Kutumb | **Time:** 60 minutes

---

## Problem Statement

Build a Least Recently Used (LRU) Cache with O(1) average time complexity for all operations. The cache has a fixed capacity; when it is full, the least recently used entry is evicted to make room. Entries support optional TTL (time-to-live). Registered listeners are notified whenever an entry is evicted — whether due to capacity overflow, TTL expiry, or explicit deletion.

**Constraints:**
- Capacity: 1 <= capacity <= 10^4
- Keys and values are integers
- All get/put/delete operations must run in O(1) average time
- TTL of 0 means no expiry; TTL > 0 means entry expires at `insertTime + ttl` seconds
- Expired entries are lazily evicted (removed on access, not proactively scanned)

---

## Base Requirement — LRU Cache with O(1) Operations

Implement an `LRUCache` supporting `get` and `put` in O(1) time. On capacity overflow, evict the least recently used entry. Any access (get or put) counts as a "use" and moves the entry to the most-recently-used position.

**Example:**
```
LRUCache cache(3)
cache.put(1, 10)
cache.put(2, 20)
cache.put(3, 30)
cache.get(1)      →  10    // 1 is now most recently used
cache.put(4, 40)           // evicts 2 (least recently used)
cache.get(2)      →  -1    // evicted
cache.get(3)      →  30
cache.get(4)      →  40
```

**Public methods:**
- `LRUCache(int capacity)`
- `int get(int key)`
- `void put(int key, int value)`

---

## Extension 1 — Per-Entry TTL

Each entry can now have an optional TTL in seconds. `get` and `put` accept a `currentTime` parameter. An entry whose `expiresAt <= currentTime` is treated as non-existent and evicted lazily on access.

**Example:**
```
cache.put(1, 100, currentTime=0, ttl=60)   // expires at t=60
cache.put(2, 200, currentTime=0, ttl=0)    // no expiry

cache.get(1, currentTime=59)  →  100   // still valid
cache.get(1, currentTime=61)  →  -1    // expired and lazily removed
cache.get(2, currentTime=999) →  200   // no TTL, still alive

cache.deleteKey(2)   →  true
cache.deleteKey(2)   →  false   // already removed
cache.size()         →  0       // all entries gone
```

**Public methods:**
- `void put(int key, int value, long currentTime, int ttl = 0)`
- `int get(int key, long currentTime)`
- `bool deleteKey(int key)`
- `int size()`

---

## Extension 2 — Eviction Listeners

Register listeners that are notified synchronously whenever an entry is evicted, along with the reason.

| Eviction Reason | When |
|---|---|
| CAPACITY | Cache is full; LRU entry bumped to make room |
| TTL_EXPIRED | Entry's TTL elapsed; evicted on lazy access |
| EXPLICIT_DELETE | Caller invoked `deleteKey` |

**Example:**
```
cache.addEvictionListener(&logger)
cache.put(5, 50, currentTime=0, ttl=10)
cache.get(5, currentTime=20)  →  -1
// logger.onEviction(5, 50, TTL_EXPIRED) was called during the get
```

**Public methods:**
- `void addEvictionListener(EvictionListener* listener)`
- `void removeEvictionListener(EvictionListener* listener)`

---

## Running Tests

```bash
./run-tests.sh 017-lru-cache cpp
```
