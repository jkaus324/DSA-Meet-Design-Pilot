package main

// --- Eviction Reason ---------------------------------------------------------

type EvictionReason int

const (
	EvictCapacity       EvictionReason = iota
	EvictTTLExpired     EvictionReason = iota
	EvictExplicitDelete EvictionReason = iota
)

// --- Eviction Listener Interface ---------------------------------------------

type EvictionListener interface {
	OnEviction(key, value int, reason EvictionReason)
}

// --- Node --------------------------------------------------------------------

type Node struct {
	key       int
	value     int
	expiresAt int64
	prev      *Node
	next      *Node
}

// --- LRU Cache ---------------------------------------------------------------

type LRUCache struct {
	capacity  int
	cache     map[int]*Node
	head      *Node
	tail      *Node
	listeners []EvictionListener
}

func NewLRUCache(capacity int) *LRUCache {
	// TODO: Create sentinel head/tail; link them; init cache map; empty listeners
	return nil
}

func (c *LRUCache) addToFront(node *Node) {
	// TODO: Insert node after head sentinel
}

func (c *LRUCache) removeNode(node *Node) {
	// TODO: Unlink node from doubly-linked list
}

func (c *LRUCache) moveToFront(node *Node) {
	// TODO: removeNode + addToFront
}

func (c *LRUCache) isExpired(node *Node, currentTime int64) bool {
	// TODO: return node.expiresAt > 0 && currentTime >= node.expiresAt
	return false
}

func (c *LRUCache) notifyListeners(key, value int, reason EvictionReason) {
	// TODO: Call l.OnEviction(key, value, reason) for each listener
}

func (c *LRUCache) evictNode(node *Node, reason EvictionReason) {
	// TODO: removeNode(node); delete(c.cache, node.key); notifyListeners(...)
}

func (c *LRUCache) evictLRU() {
	// TODO: If empty (head.next == tail), return; else evictNode(tail.prev, EvictCapacity)
}

func (c *LRUCache) Get(key int, currentTime int64) int {
	// TODO: Key missing → return -1
	// TODO: Expired → evictNode(EvictTTLExpired); return -1
	// TODO: Move to front; return value
	return -1
}

func (c *LRUCache) Put(key int, value int, currentTime int64, ttl int) {
	// TODO: Compute expiresAt (0 if ttl <= 0)
	// TODO: Key exists + expired → evictNode(EvictTTLExpired) then fall through
	// TODO: Key exists + not expired → update value/expiresAt + moveToFront; return
	// TODO: At capacity → evictLRU
	// TODO: Create new node; addToFront; c.cache[key] = node
}

func (c *LRUCache) Delete(key int) bool {
	// TODO: Key missing → return false
	// TODO: evictNode(c.cache[key], EvictExplicitDelete); return true
	return false
}

func (c *LRUCache) Size() int {
	// TODO: return len(c.cache)
	return 0
}

func (c *LRUCache) AddEvictionListener(l EvictionListener) {
	// TODO: append l to c.listeners
}

func (c *LRUCache) RemoveEvictionListener(l EvictionListener) {
	// TODO: Find l in c.listeners by interface equality and remove it
}
