package main

// --- Your Design Starts Here (Part 2) ---------------------------------------
//
// Extend the LRU Cache to support TTL (time-to-live) per key:
//
//   Get(key int, currentTime int64) int
//     — If key is expired (expiresAt > 0 && currentTime >= expiresAt), evict it
//       and return -1. Otherwise move to front and return value.
//
//   Put(key int, value int, currentTime int64, ttl int)
//     — ttl is in seconds; expiresAt = currentTime + ttl if ttl > 0, else 0 (no expiry).
//     — If the key already exists and is expired, evict it first (TTL_EXPIRED).
//     — If at capacity, evict LRU.
//     — Insert the new node at front.
//
//   Delete(key int) bool
//     — Remove the key explicitly; return false if not present.
//
//   Size() int
//     — Return the current number of entries.
//
// Think about:
//   - How do you extend Node to carry an expiresAt timestamp?
//   - When should expiry be checked — only on access, or proactively?
//
// Entry points (must exist for tests — include Part 1 entry points too):
//   func NewLRUCache(capacity int) *LRUCache
//   func (c *LRUCache) Get(key int, currentTime int64) int
//   func (c *LRUCache) Put(key int, value int, currentTime int64, ttl int)
//   func (c *LRUCache) Delete(key int) bool
//   func (c *LRUCache) Size() int

// -------------------------------------------------------------------------
