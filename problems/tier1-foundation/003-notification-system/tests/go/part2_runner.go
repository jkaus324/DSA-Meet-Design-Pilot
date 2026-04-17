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

	// Test 1: critical event reaches all users regardless of min priority
	test("test_critical_reaches_all", func() {
		u1 := User{"u1", "u1@test.com", "+1-555-0001", []string{"email"}}
		prefs := map[string]string{"*": "critical"}
		// Should not panic
		Notify("CRITICAL: System down", "critical", []User{u1}, prefs)
	})

	// Test 2: promotional event blocked when user wants info+ only
	test("test_promotional_filtered", func() {
		u1 := User{"u1", "u1@test.com", "+1-555-0001", []string{"email"}}
		prefs := map[string]string{"*": "info"}
		// Promotional should be filtered out — no panic expected
		Notify("50% off sale!", "promotional", []User{u1}, prefs)
	})

	// Test 3: empty priority prefs defaults to sending all events
	test("test_empty_prefs_allow_all", func() {
		u1 := User{"u1", "u1@test.com", "+1-555-0001", []string{"email"}}
		emptyPrefs := map[string]string{}
		Notify("Informational update", "info", []User{u1}, emptyPrefs)
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
