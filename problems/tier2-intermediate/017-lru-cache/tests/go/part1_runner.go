package main

import "fmt"

func part1Tests() int {
	passed := 0
	failed := 0

	test := func(name string, fn func()) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("FAIL", name)
					failed++
				}
			}()
			fn()
			fmt.Println("PASS", name)
			passed++
		}()
	}

	// Test 1: Get from empty cache returns -1
	test("test_get_empty_cache", func() {
		c := NewLRUCache(3)
		if c.Get(1) != -1 {
			panic("empty cache should return -1")
		}
	})

	// Test 2: Put and Get basic
	test("test_put_and_get", func() {
		c := NewLRUCache(3)
		c.Put(1, 10)
		if c.Get(1) != 10 {
			panic("expected value 10")
		}
	})

	// Test 3: Get missing key returns -1
	test("test_get_missing_key", func() {
		c := NewLRUCache(3)
		c.Put(1, 10)
		if c.Get(2) != -1 {
			panic("missing key should return -1")
		}
	})

	// Test 4: Update existing key
	test("test_update_existing_key", func() {
		c := NewLRUCache(3)
		c.Put(1, 10)
		c.Put(1, 20)
		if c.Get(1) != 20 {
			panic("expected updated value 20")
		}
	})

	// Test 5: Evict LRU on capacity overflow
	test("test_evict_lru_on_overflow", func() {
		c := NewLRUCache(3)
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)
		c.Put(4, 40) // should evict key 1 (LRU)
		if c.Get(1) != -1 {
			panic("key 1 should have been evicted")
		}
		if c.Get(4) != 40 {
			panic("key 4 should be present")
		}
	})

	// Test 6: Get updates recency — accessed key is not evicted
	test("test_get_updates_recency", func() {
		c := NewLRUCache(3)
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(3, 30)
		c.Get(1) // access key 1, making key 2 the LRU
		c.Put(4, 40) // should evict key 2
		if c.Get(2) != -1 {
			panic("key 2 should have been evicted")
		}
		if c.Get(1) != 10 {
			panic("key 1 should still be present")
		}
	})

	// Test 7: Put updates recency — updated key is not evicted
	test("test_put_updates_recency", func() {
		c := NewLRUCache(2)
		c.Put(1, 10)
		c.Put(2, 20)
		c.Put(1, 15) // update key 1, making key 2 the LRU
		c.Put(3, 30) // should evict key 2
		if c.Get(2) != -1 {
			panic("key 2 should have been evicted")
		}
		if c.Get(1) != 15 {
			panic("key 1 should be present with updated value 15")
		}
	})

	// Test 8: Capacity-1 cache evicts immediately on second insert
	test("test_capacity_one_evicts_on_second", func() {
		c := NewLRUCache(1)
		c.Put(1, 10)
		c.Put(2, 20)
		if c.Get(1) != -1 {
			panic("key 1 should be evicted from capacity-1 cache")
		}
		if c.Get(2) != 20 {
			panic("key 2 should be present")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
