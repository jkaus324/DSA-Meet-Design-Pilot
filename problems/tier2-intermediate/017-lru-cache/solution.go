// LRU Cache with TTL expiry and eviction listeners (Go).
package main

import (
	"container/list"
	"fmt"
)

type LruOp struct {
	kind string
	i1   int
	i2   int
	i3   int
	i4   int
}

// Eviction reasons
const (
	reasonCapacity        = "CAPACITY"
	reasonTTLExpired      = "TTL_EXPIRED"
	reasonExplicitDelete  = "EXPLICIT_DELETE"
)

type evictionEvent struct {
	key    int
	value  int
	reason string
}

type capturingListener struct {
	events []evictionEvent
}

type entry struct {
	key       int
	value     int
	expiresAt int
}

// lruCache preserves ordering: front of list = LRU, back = MRU.
type lruCache struct {
	capacity  int
	ll        *list.List
	items     map[int]*list.Element
	listeners []*capturingListener
}

func newLRUCache(capacity int) *lruCache {
	return &lruCache{
		capacity:  capacity,
		ll:        list.New(),
		items:     map[int]*list.Element{},
		listeners: []*capturingListener{},
	}
}

func (c *lruCache) isExpired(expiresAt, currentTime int) bool {
	return expiresAt > 0 && currentTime >= expiresAt
}

func (c *lruCache) notify(key, value int, reason string) {
	for _, l := range c.listeners {
		l.events = append(l.events, evictionEvent{key: key, value: value, reason: reason})
	}
}

func (c *lruCache) evictKey(key int, reason string) {
	el := c.items[key]
	e := el.Value.(*entry)
	value := e.value
	c.ll.Remove(el)
	delete(c.items, key)
	c.notify(key, value, reason)
}

func (c *lruCache) evictLRU() {
	if c.ll.Len() == 0 {
		return
	}
	front := c.ll.Front()
	key := front.Value.(*entry).key
	c.evictKey(key, reasonCapacity)
}

func (c *lruCache) moveToEnd(key int) {
	c.ll.MoveToBack(c.items[key])
}

// Part 1
func (c *lruCache) get(key int) int {
	el, ok := c.items[key]
	if !ok {
		return -1
	}
	c.moveToEnd(key)
	return el.Value.(*entry).value
}

func (c *lruCache) put(key, value int) {
	if el, ok := c.items[key]; ok {
		e := el.Value.(*entry)
		e.value = value
		c.moveToEnd(key)
		return
	}
	if c.ll.Len() >= c.capacity {
		c.evictLRU()
	}
	el := c.ll.PushBack(&entry{key: key, value: value, expiresAt: 0})
	c.items[key] = el
}

// Part 2 (TTL)
func (c *lruCache) getT(key, currentTime int) int {
	el, ok := c.items[key]
	if !ok {
		return -1
	}
	e := el.Value.(*entry)
	if c.isExpired(e.expiresAt, currentTime) {
		c.evictKey(key, reasonTTLExpired)
		return -1
	}
	c.moveToEnd(key)
	return e.value
}

func (c *lruCache) putT(key, value, currentTime, ttl int) {
	expiresAt := 0
	if ttl > 0 {
		expiresAt = currentTime + ttl
	}
	if el, ok := c.items[key]; ok {
		e := el.Value.(*entry)
		if c.isExpired(e.expiresAt, currentTime) {
			c.evictKey(key, reasonTTLExpired)
		} else {
			e.value = value
			e.expiresAt = expiresAt
			c.moveToEnd(key)
			return
		}
	}
	if c.ll.Len() >= c.capacity {
		c.evictLRU()
	}
	el := c.ll.PushBack(&entry{key: key, value: value, expiresAt: expiresAt})
	c.items[key] = el
}

func (c *lruCache) deleteKey(key int) bool {
	if _, ok := c.items[key]; !ok {
		return false
	}
	c.evictKey(key, reasonExplicitDelete)
	return true
}

func (c *lruCache) size() int {
	return c.ll.Len()
}

// Part 3
func (c *lruCache) addEvictionListener(l *capturingListener) {
	c.listeners = append(c.listeners, l)
}

func (c *lruCache) removeEvictionListener(target *capturingListener) {
	res := []*capturingListener{}
	for _, l := range c.listeners {
		if l != target {
			res = append(res, l)
		}
	}
	c.listeners = res
}

func lru_simulate(ops []LruOp) []string {
	out := []string{}
	var cache *lruCache
	listeners := []*capturingListener{}

	ensureListener := func(idx int) {
		for len(listeners) <= idx {
			listeners = append(listeners, &capturingListener{})
		}
	}

	for _, op := range ops {
		k := op.kind
		switch k {
		case "new":
			cache = newLRUCache(op.i1)
			listeners = []*capturingListener{}
			out = append(out, "ok")
		case "put":
			cache.put(op.i1, op.i2)
			out = append(out, "ok")
		case "put_t":
			cache.putT(op.i1, op.i2, op.i3, op.i4)
			out = append(out, "ok")
		case "get":
			out = append(out, fmt.Sprintf("%d", cache.get(op.i1)))
		case "get_t":
			out = append(out, fmt.Sprintf("%d", cache.getT(op.i1, op.i2)))
		case "delete":
			out = append(out, okFail(cache.deleteKey(op.i1)))
		case "size":
			out = append(out, fmt.Sprintf("%d", cache.size()))
		case "add_listener":
			ensureListener(op.i1)
			cache.addEvictionListener(listeners[op.i1])
			out = append(out, "ok")
		case "remove_listener":
			if op.i1 < len(listeners) {
				cache.removeEvictionListener(listeners[op.i1])
			}
			out = append(out, "ok")
		case "events_count":
			if op.i1 < len(listeners) {
				out = append(out, fmt.Sprintf("%d", len(listeners[op.i1].events)))
			} else {
				out = append(out, "0")
			}
		case "event_key":
			if op.i1 < len(listeners) && op.i2 < len(listeners[op.i1].events) {
				out = append(out, fmt.Sprintf("%d", listeners[op.i1].events[op.i2].key))
			} else {
				out = append(out, "")
			}
		case "event_value":
			if op.i1 < len(listeners) && op.i2 < len(listeners[op.i1].events) {
				out = append(out, fmt.Sprintf("%d", listeners[op.i1].events[op.i2].value))
			} else {
				out = append(out, "")
			}
		case "event_reason":
			if op.i1 < len(listeners) && op.i2 < len(listeners[op.i1].events) {
				out = append(out, listeners[op.i1].events[op.i2].reason)
			} else {
				out = append(out, "")
			}
		default:
			out = append(out, "unknown:"+k)
		}
	}
	return out
}

func okFail(b bool) string {
	if b {
		return "ok"
	}
	return "fail"
}
