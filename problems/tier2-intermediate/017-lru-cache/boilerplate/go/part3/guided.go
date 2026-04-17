package main

// --- LRU Cache with Eviction Listeners ----------------------------------------
//
// HINT: Define the EvictionListener interface:
//         type EvictionListener interface {
//             OnEviction(key, value int, reason EvictionReason)
//         }
//
// HINT: Add []EvictionListener to LRUCache.
// HINT: notifyListeners(key, value int, reason EvictionReason) iterates all listeners
//       and calls OnEviction.
// HINT: evictNode now calls notifyListeners before (or after) removing the node.
//
// HINT: AddEvictionListener appends to the slice.
// HINT: RemoveEvictionListener scans the slice for the same pointer and removes it:
//         for i, l := range c.listeners {
//             if l == target { c.listeners = append(c.listeners[:i], c.listeners[i+1:]...); break }
//         }
//
// HINT: The rest of the implementation is identical to Part 2 — only evictNode
//       gains the notification side-effect.

// type EvictionListener interface {
//     OnEviction(key, value int, reason EvictionReason)
// }

// type LRUCache struct {
//     capacity  int
//     cache     map[int]*Node
//     head      *Node
//     tail      *Node
//     listeners []EvictionListener
// }

// func (c *LRUCache) AddEvictionListener(l EvictionListener)
// func (c *LRUCache) RemoveEvictionListener(l EvictionListener)
