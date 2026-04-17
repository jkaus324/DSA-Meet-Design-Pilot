package main

import "fmt"

func part2Tests() int {
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

	// Test 1: Get on expired key returns -1
	test("test_get_expired_key_returns_minus_one", func() {
		c := NewLRUCache(3)
		c.Put(1, 10, 0, 100) // expires at t=100
		if c.Get(1, 50) != 10 {
			panic("key should be alive at t=50")
		}
		if c.Get(1, 100) != -1 {
			panic("key should be expired at t=100")
		}
	})

	// Test 2: Put with no TTL (ttl=0) never expires
	test("test_put_no_ttl_never_expires", func() {
		c := NewLRUCache(3)
		c.Put(1, 10, 0, 0)
		if c.Get(1, 999999) != 10 {
			panic("key with no TTL should never expire")
		}
	})

	// Test 3: Expired key frees capacity
	test("test_expired_key_frees_capacity", func() {
		c := NewLRUCache(2)
		c.Put(1, 10, 0, 100) // expires at t=100
		c.Put(2, 20, 0, 0)
		// At t=200, key 1 is expired; new put should not evict key 2
		c.Put(3, 30, 200, 0)
		if c.Get(2, 200) != 20 {
			panic("key 2 should not be evicted if key 1 was expired")
		}
	})

	// Test 4: Delete removes a key
	test("test_delete_removes_key", func() {
		c := NewLRUCache(3)
		c.Put(1, 10, 0, 0)
		ok := c.Delete(1)
		if !ok {
			panic("delete should return true for existing key")
		}
		if c.Get(1, 0) != -1 {
			panic("deleted key should not be accessible")
		}
	})

	// Test 5: Delete returns false for missing key
	test("test_delete_missing_key", func() {
		c := NewLRUCache(3)
		ok := c.Delete(99)
		if ok {
			panic("delete on missing key should return false")
		}
	})

	// Test 6: Size reflects current count
	test("test_size_reflects_count", func() {
		c := NewLRUCache(5)
		if c.Size() != 0 {
			panic("empty cache size should be 0")
		}
		c.Put(1, 10, 0, 0)
		c.Put(2, 20, 0, 0)
		if c.Size() != 2 {
			panic("size should be 2 after two puts")
		}
		c.Delete(1)
		if c.Size() != 1 {
			panic("size should be 1 after delete")
		}
	})

	// Test 7: Size decreases when TTL-expired key is accessed
	test("test_size_decreases_on_ttl_eviction", func() {
		c := NewLRUCache(3)
		c.Put(1, 10, 0, 50)
		c.Put(2, 20, 0, 0)
		if c.Size() != 2 {
			panic("size should be 2 before expiry")
		}
		c.Get(1, 100) // triggers eviction of key 1
		if c.Size() != 1 {
			panic("size should be 1 after TTL eviction")
		}
	})

	// Test 8: Updating an expired key inserts fresh entry
	test("test_put_on_expired_key_reinserts", func() {
		c := NewLRUCache(3)
		c.Put(1, 10, 0, 50) // expires at t=50
		c.Put(1, 99, 100, 0) // re-insert at t=100 with no expiry
		if c.Get(1, 200) != 99 {
			panic("re-inserted key should have new value 99")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
