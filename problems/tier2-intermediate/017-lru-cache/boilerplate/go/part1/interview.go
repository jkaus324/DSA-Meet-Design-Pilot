package main

// --- Your Design Starts Here ------------------------------------------------
//
// Design and implement an LRU (Least Recently Used) Cache with capacity cap.
//
// Rules:
//   - Get(key int) int  — returns the value if present, -1 otherwise.
//     Accessing a key makes it the most recently used.
//   - Put(key int, value int)  — inserts or updates the key.
//     If the cache is at capacity, evict the least recently used key first.
//
// Think about:
//   - What data structure gives O(1) get and O(1) eviction?
//   - How do you track recency order efficiently?
//
// Entry points (must exist for tests):
//   func NewLRUCache(capacity int) *LRUCache
//   func (c *LRUCache) Get(key int) int
//   func (c *LRUCache) Put(key int, value int)

// -------------------------------------------------------------------------
