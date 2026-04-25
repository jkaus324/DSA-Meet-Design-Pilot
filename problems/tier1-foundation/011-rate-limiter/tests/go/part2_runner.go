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

	// Test 1: Factory creates fixed-window limiter
	test("test_factory_fixed_window", func() {
		limiter := CreateLimiter("fixed-window", 3, 60)
		if limiter == nil {
			panic("expected non-nil limiter")
		}
		if !limiter.AllowRequest(Request{ClientId: "user_1", Timestamp: 1000, Endpoint: "/api/a"}) {
			panic("1st should pass")
		}
		if !limiter.AllowRequest(Request{ClientId: "user_1", Timestamp: 1001, Endpoint: "/api/a"}) {
			panic("2nd should pass")
		}
		if !limiter.AllowRequest(Request{ClientId: "user_1", Timestamp: 1002, Endpoint: "/api/a"}) {
			panic("3rd should pass")
		}
		if limiter.AllowRequest(Request{ClientId: "user_1", Timestamp: 1003, Endpoint: "/api/a"}) {
			panic("4th should be rejected")
		}
	})

	// Test 2: Sliding-window counts in rolling window
	test("test_sliding_window", func() {
		limiter := CreateLimiter("sliding-window", 3, 60)
		if limiter == nil {
			panic("expected non-nil limiter")
		}
		if !limiter.AllowRequest(Request{ClientId: "user_2", Timestamp: 1000, Endpoint: "/api/b"}) {
			panic("1st should pass")
		}
		if !limiter.AllowRequest(Request{ClientId: "user_2", Timestamp: 1030, Endpoint: "/api/b"}) {
			panic("2nd should pass")
		}
		if !limiter.AllowRequest(Request{ClientId: "user_2", Timestamp: 1050, Endpoint: "/api/b"}) {
			panic("3rd should pass")
		}
		if limiter.AllowRequest(Request{ClientId: "user_2", Timestamp: 1055, Endpoint: "/api/b"}) {
			panic("4th should be rejected (all 3 still in window)")
		}
		// After first request expires from window (timestamp 1000 is outside [1001, 1061])
		if !limiter.AllowRequest(Request{ClientId: "user_2", Timestamp: 1061, Endpoint: "/api/b"}) {
			panic("should pass after expiry")
		}
	})

	// Test 3: Token-bucket allows bursts then throttles
	test("test_token_bucket", func() {
		limiter := CreateLimiter("token-bucket", 3, 60)
		if limiter == nil {
			panic("expected non-nil limiter")
		}
		// Burst: use all 3 tokens immediately
		if !limiter.AllowRequest(Request{ClientId: "user_3", Timestamp: 1000, Endpoint: "/api/c"}) {
			panic("1st should pass")
		}
		if !limiter.AllowRequest(Request{ClientId: "user_3", Timestamp: 1000, Endpoint: "/api/c"}) {
			panic("2nd should pass")
		}
		if !limiter.AllowRequest(Request{ClientId: "user_3", Timestamp: 1000, Endpoint: "/api/c"}) {
			panic("3rd should pass")
		}
		if limiter.AllowRequest(Request{ClientId: "user_3", Timestamp: 1000, Endpoint: "/api/c"}) {
			panic("4th should be rejected (empty bucket)")
		}
		// Wait for tokens to refill (rate = 3/60 = 0.05/sec, need 20sec for 1 token)
		if !limiter.AllowRequest(Request{ClientId: "user_3", Timestamp: 1020, Endpoint: "/api/c"}) {
			panic("should pass after 20s refill")
		}
		if limiter.AllowRequest(Request{ClientId: "user_3", Timestamp: 1020, Endpoint: "/api/c"}) {
			panic("should be rejected (empty again)")
		}
	})

	// Test 4: Factory returns nil for unknown algorithm
	test("test_factory_unknown_algorithm", func() {
		limiter := CreateLimiter("unknown-algo", 10, 60)
		if limiter != nil {
			panic("expected nil for unknown algorithm")
		}
	})

	// Test 5: allow_request_with_strategy uses correct algorithm
	test("test_allow_request_with_strategy", func() {
		if !AllowRequestWithStrategy("fixed-window", Request{ClientId: "user_4", Timestamp: 2000, Endpoint: "/api/d"}) {
			panic("1st should pass")
		}
		if !AllowRequestWithStrategy("fixed-window", Request{ClientId: "user_4", Timestamp: 2001, Endpoint: "/api/d"}) {
			panic("2nd should pass")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
