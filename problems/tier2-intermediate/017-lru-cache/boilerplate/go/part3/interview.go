package main

// --- Your Design Starts Here (Part 3) ---------------------------------------
//
// Extend Part 2 with eviction listeners (Observer pattern):
//
//   EvictionListener interface:
//     OnEviction(key, value int, reason EvictionReason)
//
//   AddEvictionListener(l EvictionListener)
//     — Register a listener; all future evictions notify it.
//
//   RemoveEvictionListener(l EvictionListener)
//     — Deregister a previously added listener.
//
// All three eviction paths must notify listeners:
//   - Capacity eviction (Put when full)
//   - TTL expiry (Get or Put on an expired key)
//   - Explicit delete (Delete)
//
// Think about:
//   - Should eviction listeners be called while holding any lock? (single-threaded is fine)
//   - How do you safely remove a listener from a slice by reference?
//
// Entry points (must exist for tests — include all prior entry points):
//   func NewLRUCache(capacity int) *LRUCache
//   func (c *LRUCache) Get(key int, currentTime int64) int
//   func (c *LRUCache) Put(key int, value int, currentTime int64, ttl int)
//   func (c *LRUCache) Delete(key int) bool
//   func (c *LRUCache) Size() int
//   func (c *LRUCache) AddEvictionListener(l EvictionListener)
//   func (c *LRUCache) RemoveEvictionListener(l EvictionListener)

// -------------------------------------------------------------------------
