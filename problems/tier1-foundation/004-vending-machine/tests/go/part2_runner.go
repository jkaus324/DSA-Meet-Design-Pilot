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

	// Test 1: valid PIN enters maintenance mode
	test("test_enter_maintenance_valid_pin", func() {
		Reset()
		EnterMaintenance("1234")
		if GetState() != "Maintenance" {
			panic("expected Maintenance state")
		}
	})

	// Test 2: invalid PIN does not enter maintenance
	test("test_enter_maintenance_invalid_pin", func() {
		Reset()
		EnterMaintenance("wrong")
		if GetState() != "Idle" {
			panic("expected machine to stay in Idle with invalid PIN")
		}
	})

	// Test 3: in maintenance, user actions are blocked
	test("test_user_blocked_in_maintenance", func() {
		Reset()
		EnterMaintenance("1234")
		SelectItem("Cola") // should print warning, not panic
		if GetState() != "Maintenance" {
			panic("expected machine to stay in Maintenance")
		}
	})

	// Test 4: exit maintenance returns to Idle
	test("test_exit_maintenance", func() {
		Reset()
		EnterMaintenance("1234")
		ExitMaintenance("1234")
		if GetState() != "Idle" {
			panic("expected Idle state after exiting maintenance")
		}
	})

	// Test 5: restock only works in maintenance
	test("test_restock_in_maintenance", func() {
		Reset()
		Restock("Cola", 10) // should be blocked outside maintenance
		EnterMaintenance("1234")
		Restock("Cola", 5) // should work in maintenance
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
