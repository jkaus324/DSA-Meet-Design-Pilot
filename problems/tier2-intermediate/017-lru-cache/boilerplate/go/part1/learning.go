package main

// --- Node (doubly-linked list) -----------------------------------------------

type Node struct {
	key   int
	value int
	prev  *Node
	next  *Node
}

// --- LRU Cache ---------------------------------------------------------------

type LRUCache struct {
	capacity int
	cache    map[int]*Node
	head     *Node // sentinel: most-recently-used end
	tail     *Node // sentinel: least-recently-used end
}

func NewLRUCache(capacity int) *LRUCache {
	// TODO: Create sentinel head and tail nodes linked to each other:
	//         head.next = tail; tail.prev = head
	// TODO: Initialise cache map; set capacity
	// TODO: Return the LRUCache
	return nil
}

func (c *LRUCache) addToFront(node *Node) {
	// TODO: Insert node between head and head.next
	//         node.next = head.next
	//         node.prev = head
	//         head.next.prev = node
	//         head.next = node
}

func (c *LRUCache) removeNode(node *Node) {
	// TODO: Unlink node: node.prev.next = node.next; node.next.prev = node.prev
}

func (c *LRUCache) moveToFront(node *Node) {
	// TODO: removeNode(node); addToFront(node)
}

func (c *LRUCache) evictLRU() {
	// TODO: If head.next == tail (empty), return
	// TODO: lru = tail.prev; removeNode(lru); delete(c.cache, lru.key)
}

func (c *LRUCache) Get(key int) int {
	// TODO: If key not in c.cache, return -1
	// TODO: moveToFront(c.cache[key]); return c.cache[key].value
	return -1
}

func (c *LRUCache) Put(key int, value int) {
	// TODO: If key exists: update value, moveToFront; return
	// TODO: If len(c.cache) >= c.capacity: evictLRU
	// TODO: Create new Node{key, value}; addToFront; store in c.cache
}
