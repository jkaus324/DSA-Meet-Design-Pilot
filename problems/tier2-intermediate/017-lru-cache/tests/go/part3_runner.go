package main

import "fmt"

// MockEvictionListener records evictions for test assertions.
type MockEvictionListener struct {
	Evictions []struct {
		Key    int
		Value  int
		Reason EvictionReason
	}
}

func (m *MockEvictionListener) OnEviction(key, value int, reason EvictionReason) {
	m.Evictions = append(m.Evictions, struct {
		Key    int
		Value  int
		Reason EvictionReason
	}{key, value, reason})
}

func part3Tests() int {
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

	// Test 1: Capacity eviction notifies listener
	test("test_capacity_eviction_notifies_listener", func() {
		c := NewLRUCache(2)
		mock := &MockEvictionListener{}
		c.AddEvictionListener(mock)
		c.Put(1, 10, 0, 0)
		c.Put(2, 20, 0, 0)
		c.Put(3, 30, 0, 0) // evicts key 1
		if len(mock.Evictions) != 1 {
			panic("expected 1 eviction notification")
		}
		ev := mock.Evictions[0]
		if ev.Key != 1 || ev.Value != 10 || ev.Reason != EvictCapacity {
			panic("wrong eviction details for capacity eviction")
		}
	})

	// Test 2: TTL eviction notifies listener
	test("test_ttl_eviction_notifies_listener", func() {
		c := NewLRUCache(3)
		mock := &MockEvictionListener{}
		c.AddEvictionListener(mock)
		c.Put(1, 10, 0, 50) // expires at t=50
		c.Get(1, 100)        // triggers TTL eviction
		if len(mock.Evictions) != 1 {
			panic("expected 1 TTL eviction notification")
		}
		if mock.Evictions[0].Reason != EvictTTLExpired {
			panic("reason should be EvictTTLExpired")
		}
	})

	// Test 3: Explicit delete notifies listener
	test("test_explicit_delete_notifies_listener", func() {
		c := NewLRUCache(3)
		mock := &MockEvictionListener{}
		c.AddEvictionListener(mock)
		c.Put(1, 10, 0, 0)
		c.Delete(1)
		if len(mock.Evictions) != 1 {
			panic("expected 1 explicit-delete eviction notification")
		}
		if mock.Evictions[0].Reason != EvictExplicitDelete {
			panic("reason should be EvictExplicitDelete")
		}
	})

	// Test 4: Remove listener stops notifications
	test("test_remove_listener_stops_notifications", func() {
		c := NewLRUCache(2)
		mock := &MockEvictionListener{}
		c.AddEvictionListener(mock)
		c.Put(1, 10, 0, 0)
		c.Put(2, 20, 0, 0)
		c.RemoveEvictionListener(mock)
		c.Put(3, 30, 0, 0) // evicts key 1, but listener was removed
		if len(mock.Evictions) != 0 {
			panic("removed listener should not receive notifications")
		}
	})

	// Test 5: Multiple listeners all notified
	test("test_multiple_listeners_all_notified", func() {
		c := NewLRUCache(1)
		mock1 := &MockEvictionListener{}
		mock2 := &MockEvictionListener{}
		c.AddEvictionListener(mock1)
		c.AddEvictionListener(mock2)
		c.Put(1, 10, 0, 0)
		c.Put(2, 20, 0, 0) // evicts key 1
		if len(mock1.Evictions) != 1 || len(mock2.Evictions) != 1 {
			panic("both listeners should receive the eviction notification")
		}
	})

	// Test 6: Eviction reason is correct per eviction type
	test("test_eviction_reasons_are_correct", func() {
		c := NewLRUCache(2)
		mock := &MockEvictionListener{}
		c.AddEvictionListener(mock)

		// Put with TTL → TTL eviction on Get
		c.Put(1, 10, 0, 100)
		c.Get(1, 200) // TTL eviction

		// Put until capacity → capacity eviction
		c.Put(2, 20, 300, 0)
		c.Put(3, 30, 300, 0)
		c.Put(4, 40, 300, 0) // evicts key 2 (capacity)

		reasons := map[EvictionReason]int{}
		for _, e := range mock.Evictions {
			reasons[e.Reason]++
		}
		if reasons[EvictTTLExpired] != 1 {
			panic("expected 1 TTL eviction")
		}
		if reasons[EvictCapacity] != 1 {
			panic("expected 1 capacity eviction")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
