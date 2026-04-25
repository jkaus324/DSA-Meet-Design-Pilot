package main

import "fmt"

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

	// Test 1: FREE tier allows only 10 requests per minute
	test("test_free_tier_limit", func() {
		for i := 0; i < 10; i++ {
			if !AllowRequestForTier(FREE, Request{ClientId: "free_user", Timestamp: int64(5000 + i), Endpoint: "/api/data"}) {
				panic(fmt.Sprintf("request %d should be allowed", i+1))
			}
		}
		// 11th should be rejected
		if AllowRequestForTier(FREE, Request{ClientId: "free_user", Timestamp: 5010, Endpoint: "/api/data"}) {
			panic("11th request should be rejected for FREE tier")
		}
	})

	// Test 2: PRO tier allows 100 requests per minute
	test("test_pro_tier_limit", func() {
		for i := 0; i < 100; i++ {
			if !AllowRequestForTier(PRO, Request{ClientId: "pro_user", Timestamp: int64(6000 + i), Endpoint: "/api/data"}) {
				panic(fmt.Sprintf("request %d should be allowed", i+1))
			}
		}
		// 101st should be rejected
		if AllowRequestForTier(PRO, Request{ClientId: "pro_user", Timestamp: 6100, Endpoint: "/api/data"}) {
			panic("101st request should be rejected for PRO tier")
		}
	})

	// Test 3: ENTERPRISE tier allows 1000 requests per minute
	test("test_enterprise_tier_limit", func() {
		for i := 0; i < 1000; i++ {
			if !AllowRequestForTier(ENTERPRISE, Request{ClientId: "enterprise_user", Timestamp: int64(7000 + i), Endpoint: "/api/data"}) {
				panic(fmt.Sprintf("request %d should be allowed", i+1))
			}
		}
		// 1001st should be rejected
		if AllowRequestForTier(ENTERPRISE, Request{ClientId: "enterprise_user", Timestamp: 8000, Endpoint: "/api/data"}) {
			panic("1001st request should be rejected for ENTERPRISE tier")
		}
	})

	// Test 4: different tiers are independent
	test("test_tier_independence", func() {
		// Free user hits limit at 10
		for i := 0; i < 10; i++ {
			AllowRequestForTier(FREE, Request{ClientId: "free_user_2", Timestamp: int64(9000 + i), Endpoint: "/api/x"})
		}
		if AllowRequestForTier(FREE, Request{ClientId: "free_user_2", Timestamp: 9010, Endpoint: "/api/x"}) {
			panic("free_user_2 should be rejected")
		}
		// Pro user still has quota
		if !AllowRequestForTier(PRO, Request{ClientId: "pro_user_2", Timestamp: 9010, Endpoint: "/api/x"}) {
			panic("pro_user_2 should still pass")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
