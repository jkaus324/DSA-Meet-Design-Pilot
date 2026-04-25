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

	// Test 1: notify sends to subscribed channels (no panic)
	test("test_notify_no_throw", func() {
		u1 := User{"user1", "u1@test.com", "+91-9000000001", []string{"email"}}
		u2 := User{"user2", "u2@test.com", "+91-9000000002", []string{"sms", "push"}}
		u3 := User{"user3", "u3@test.com", "+91-9000000003", []string{"email", "sms"}}
		Notify("Order shipped", []User{u1, u2, u3})
	})

	// Test 2: empty user list causes no crash
	test("test_empty_user_list", func() {
		Notify("Event", []User{})
	})

	// Test 3: user subscribed to multiple channels receives on each (no panic)
	test("test_multi_channel_user", func() {
		u := User{"u1", "u1@test.com", "+1-555-0001", []string{"email", "sms", "push"}}
		Notify("Flash sale", []User{u})
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
