package main

// --- LRU Cache (doubly-linked list + hash map) --------------------------------
//
// HINT: Use a doubly-linked list where the head-side is most-recently-used
//       and the tail-side is least-recently-used.
//       Maintain sentinel head and tail nodes to avoid nil checks.
//
// HINT: Node struct:
//         type Node struct { key, value int; prev, next *Node }
//
// HINT: LRUCache struct:
//         type LRUCache struct {
//             capacity int
//             cache    map[int]*Node
//             head     *Node  // sentinel MRU end
//             tail     *Node  // sentinel LRU end
//         }
//
// HINT: addToFront(node) — insert node right after head sentinel.
// HINT: removeNode(node) — unlink node from its current position.
// HINT: moveToFront(node) = removeNode + addToFront.
// HINT: evictLRU() — remove tail.prev (the actual LRU node), delete from cache map.
//
// HINT: Get — if key missing return -1; else moveToFront and return value.
// HINT: Put — if key exists: update value + moveToFront.
//             if at capacity: evictLRU first.
//             Create new node, addToFront, store in cache.

// type Node struct { key, value int; prev, next *Node }

// type LRUCache struct {
//     capacity int
//     cache    map[int]*Node
//     head     *Node
//     tail     *Node
// }

// func NewLRUCache(capacity int) *LRUCache
// func (c *LRUCache) Get(key int) int
// func (c *LRUCache) Put(key int, value int)
