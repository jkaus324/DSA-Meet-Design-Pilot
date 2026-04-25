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

	// Test 1: allow requests within limit
	test("test_allow_within_limit", func() {
		InitLimiter(3, 60) // 3 requests per 60 seconds
		r1 := Request{ClientId: "client_A", Timestamp: 1000, Endpoint: "/api/search"}
		r2 := Request{ClientId: "client_A", Timestamp: 1001, Endpoint: "/api/search"}
		r3 := Request{ClientId: "client_A", Timestamp: 1002, Endpoint: "/api/search"}
		if !AllowRequest(r1) {
			panic("r1 should be allowed")
		}
		if !AllowRequest(r2) {
			panic("r2 should be allowed")
		}
		if !AllowRequest(r3) {
			panic("r3 should be allowed")
		}
	})

	// Test 2: reject requests exceeding limit
	test("test_reject_over_limit", func() {
		InitLimiter(3, 60)
		if !AllowRequest(Request{ClientId: "client_B", Timestamp: 2000, Endpoint: "/api/pay"}) {
			panic("1st should be allowed")
		}
		if !AllowRequest(Request{ClientId: "client_B", Timestamp: 2001, Endpoint: "/api/pay"}) {
			panic("2nd should be allowed")
		}
		if !AllowRequest(Request{ClientId: "client_B", Timestamp: 2002, Endpoint: "/api/pay"}) {
			panic("3rd should be allowed")
		}
		if AllowRequest(Request{ClientId: "client_B", Timestamp: 2003, Endpoint: "/api/pay"}) {
			panic("4th should be rejected")
		}
		if AllowRequest(Request{ClientId: "client_B", Timestamp: 2004, Endpoint: "/api/pay"}) {
			panic("5th should be rejected")
		}
	})

	// Test 3: different clients have independent limits
	test("test_independent_client_limits", func() {
		InitLimiter(2, 60)
		if !AllowRequest(Request{ClientId: "alice", Timestamp: 3000, Endpoint: "/api/x"}) {
			panic("alice 1st should pass")
		}
		if !AllowRequest(Request{ClientId: "alice", Timestamp: 3001, Endpoint: "/api/x"}) {
			panic("alice 2nd should pass")
		}
		if AllowRequest(Request{ClientId: "alice", Timestamp: 3002, Endpoint: "/api/x"}) {
			panic("alice 3rd should be rejected")
		}
		if !AllowRequest(Request{ClientId: "bob", Timestamp: 3003, Endpoint: "/api/x"}) {
			panic("bob 1st should pass")
		}
		if !AllowRequest(Request{ClientId: "bob", Timestamp: 3004, Endpoint: "/api/x"}) {
			panic("bob 2nd should pass")
		}
		if AllowRequest(Request{ClientId: "bob", Timestamp: 3005, Endpoint: "/api/x"}) {
			panic("bob 3rd should be rejected")
		}
	})

	// Test 4: window reset allows new requests
	test("test_window_reset", func() {
		InitLimiter(2, 60)
		if !AllowRequest(Request{ClientId: "client_C", Timestamp: 1000, Endpoint: "/api/y"}) {
			panic("1st should pass")
		}
		if !AllowRequest(Request{ClientId: "client_C", Timestamp: 1020, Endpoint: "/api/y"}) {
			panic("2nd should pass")
		}
		if AllowRequest(Request{ClientId: "client_C", Timestamp: 1040, Endpoint: "/api/y"}) {
			panic("3rd should be rejected (limit hit)")
		}
		// New window starts at 1060
		if !AllowRequest(Request{ClientId: "client_C", Timestamp: 1060, Endpoint: "/api/y"}) {
			panic("should pass in new window")
		}
		if !AllowRequest(Request{ClientId: "client_C", Timestamp: 1080, Endpoint: "/api/y"}) {
			panic("should pass in new window")
		}
	})

	// Test 5: get_request_count tracks correctly
	test("test_get_request_count", func() {
		InitLimiter(5, 60)
		if GetRequestCount("new_client") != 0 {
			panic("initial count should be 0")
		}
		AllowRequest(Request{ClientId: "new_client", Timestamp: 5000, Endpoint: "/api/z"})
		if GetRequestCount("new_client") != 1 {
			panic("count should be 1")
		}
		AllowRequest(Request{ClientId: "new_client", Timestamp: 5001, Endpoint: "/api/z"})
		if GetRequestCount("new_client") != 2 {
			panic("count should be 2")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
