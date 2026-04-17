package main

// --- Node (doubly-linked list with TTL) --------------------------------------

type Node struct {
	key       int
	value     int
	expiresAt int64 // 0 = no expiry
	prev      *Node
	next      *Node
}

// --- Eviction Reason (used in Part 3; define here for forward compatibility) -

type EvictionReason int

const (
	EvictCapacity       EvictionReason = iota
	EvictTTLExpired     EvictionReason = iota
	EvictExplicitDelete EvictionReason = iota
)

// --- LRU Cache ---------------------------------------------------------------

type LRUCache struct {
	capacity int
	cache    map[int]*Node
	head     *Node // sentinel MRU
	tail     *Node // sentinel LRU
}

func NewLRUCache(capacity int) *LRUCache {
	// TODO: Create sentinel head and tail; link them; initialise cache map
	return nil
}

func (c *LRUCache) addToFront(node *Node) {
	// TODO: Insert node right after head sentinel
}

func (c *LRUCache) removeNode(node *Node) {
	// TODO: Unlink node from its neighbours
}

func (c *LRUCache) moveToFront(node *Node) {
	// TODO: removeNode + addToFront
}

func (c *LRUCache) isExpired(node *Node, currentTime int64) bool {
	// TODO: return node.expiresAt > 0 && currentTime >= node.expiresAt
	return false
}

func (c *LRUCache) evictNode(node *Node, reason EvictionReason) {
	// TODO: removeNode(node); delete(c.cache, node.key)
	// (reason ignored until Part 3)
}

func (c *LRUCache) evictLRU() {
	// TODO: If head.next == tail (empty), return
	// TODO: evictNode(tail.prev, EvictCapacity)
}

func (c *LRUCache) Get(key int, currentTime int64) int {
	// TODO: If key not in cache, return -1
	// TODO: If isExpired: evictNode(EvictTTLExpired); return -1
	// TODO: moveToFront; return value
	return -1
}

func (c *LRUCache) Put(key int, value int, currentTime int64, ttl int) {
	// TODO: expiresAt = currentTime + int64(ttl) if ttl > 0 else 0
	// TODO: If key exists:
	//         if isExpired: evictNode(EvictTTLExpired) (fall through to insert)
	//         else: update value + expiresAt + moveToFront; return
	// TODO: If len(c.cache) >= capacity: evictLRU
	// TODO: Create Node{key, value, expiresAt}; addToFront; c.cache[key] = node
}

func (c *LRUCache) Delete(key int) bool {
	// TODO: If key not in cache, return false
	// TODO: evictNode(c.cache[key], EvictExplicitDelete); return true
	return false
}

func (c *LRUCache) Size() int {
	// TODO: return len(c.cache)
	return 0
}
