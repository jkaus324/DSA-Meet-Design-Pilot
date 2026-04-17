package main

// --- LRU Cache with TTL -------------------------------------------------------
//
// HINT: Extend Node to carry expiresAt int64 (0 means no expiry):
//         type Node struct { key, value int; expiresAt int64; prev, next *Node }
//
// HINT: isExpired helper:
//         return node.expiresAt > 0 && currentTime >= node.expiresAt
//
// HINT: evictNode(node, reason) — generic eviction: removeNode, delete from cache.
//       (reason is used in Part 3 for listener callbacks; ignore in Part 2)
//
// HINT: Get(key, currentTime):
//         if expired: evictNode; return -1
//         else: moveToFront; return value
//
// HINT: Put(key, value, currentTime, ttl):
//         expiresAt = currentTime + int64(ttl) if ttl > 0 else 0
//         if key exists:
//             if expired: evictNode(TTL_EXPIRED)
//             else: update value + expiresAt + moveToFront; return
//         if at capacity: evictLRU (CAPACITY)
//         addToFront new node
//
// HINT: Delete(key) bool — evictNode(EXPLICIT_DELETE); return false if not found.
// HINT: Size() int — return len(c.cache)

// type Node struct { key, value int; expiresAt int64; prev, next *Node }

// type LRUCache struct {
//     capacity int
//     cache    map[int]*Node
//     head     *Node
//     tail     *Node
// }

// func NewLRUCache(capacity int) *LRUCache
// func (c *LRUCache) Get(key int, currentTime int64) int
// func (c *LRUCache) Put(key int, value int, currentTime int64, ttl int)
// func (c *LRUCache) Delete(key int) bool
// func (c *LRUCache) Size() int
