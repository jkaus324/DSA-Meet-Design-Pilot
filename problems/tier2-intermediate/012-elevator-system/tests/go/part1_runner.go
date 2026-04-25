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

	// Test 1: Initial state — elevator starts at floor 0, IDLE
	test("test_initial_state", func() {
		e := NewElevator()
		if e.GetCurrentFloor() != 0 {
			panic("expected floor 0")
		}
		if e.GetState() != IDLE {
			panic("expected IDLE state")
		}
	})

	// Test 2: Add upward request — elevator starts moving up
	test("test_add_upward_request", func() {
		e := NewElevator()
		e.AddRequest(3, UP)
		if e.GetState() != MOVING_UP {
			panic("expected MOVING_UP")
		}
	})

	// Test 3: Step moves one floor at a time
	test("test_step_moves_one_floor", func() {
		e := NewElevator()
		e.AddRequest(3, UP)
		e.Step() // 0 -> 1
		if e.GetCurrentFloor() != 1 {
			panic("expected floor 1")
		}
		e.Step() // 1 -> 2
		if e.GetCurrentFloor() != 2 {
			panic("expected floor 2")
		}
	})

	// Test 4: Elevator opens doors at requested floor
	test("test_door_opens_at_target", func() {
		e := NewElevator()
		e.AddRequest(2, UP)
		e.Step() // 0 -> 1
		e.Step() // 1 -> 2, door opens
		if e.GetCurrentFloor() != 2 {
			panic("expected floor 2")
		}
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN")
		}
	})

	// Test 5: After door open with no more requests, goes IDLE
	test("test_idle_after_last_request", func() {
		e := NewElevator()
		e.AddRequest(1, UP)
		e.Step() // 0 -> 1, door opens
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN")
		}
		e.Step() // close doors, go idle
		if e.GetState() != IDLE {
			panic("expected IDLE")
		}
		if e.GetCurrentFloor() != 1 {
			panic("expected floor 1")
		}
	})

	// Test 6: SCAN order — serve all upward requests before reversing
	test("test_scan_upward_order", func() {
		e := NewElevator()
		e.AddRequest(5, UP)
		e.AddRequest(2, UP)
		e.Step() // 1
		e.Step() // 2, DOOR_OPEN
		if e.GetCurrentFloor() != 2 {
			panic("expected floor 2")
		}
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN")
		}
		e.Step() // close doors, resume MOVING_UP
		if e.GetState() != MOVING_UP {
			panic("expected MOVING_UP")
		}
		e.Step() // 3
		e.Step() // 4
		e.Step() // 5, DOOR_OPEN
		if e.GetCurrentFloor() != 5 {
			panic("expected floor 5")
		}
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN at floor 5")
		}
	})

	// Test 7: SCAN reversal — after upward done, serve downward requests
	test("test_scan_reversal", func() {
		e := NewElevator()
		e.AddRequest(3, UP)
		e.AddRequest(1, DOWN) // at floor 0, 1 > 0 so goes to upRequests
		e.AddRequest(3, UP)
		// Step to floor 3
		e.Step() // 1
		e.Step() // 2
		e.Step() // 3, DOOR_OPEN
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN at 3")
		}
		// Add a downward request
		e.AddRequest(1, DOWN) // 1 < 3, goes to downRequests
		e.Step()               // close doors, should switch to MOVING_DOWN
		if e.GetState() != MOVING_DOWN {
			panic("expected MOVING_DOWN")
		}
		e.Step() // 2
		e.Step() // 1, DOOR_OPEN
		if e.GetCurrentFloor() != 1 {
			panic("expected floor 1")
		}
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN at 1")
		}
	})

	// Test 8: Request at current floor while IDLE opens doors immediately
	test("test_request_at_current_floor", func() {
		e := NewElevator()
		e.AddRequest(0, UP)
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN")
		}
		if e.GetCurrentFloor() != 0 {
			panic("expected floor 0")
		}
	})

	// Test 9: Step on IDLE elevator does nothing
	test("test_idle_step_noop", func() {
		e := NewElevator()
		e.Step()
		if e.GetCurrentFloor() != 0 {
			panic("expected floor 0")
		}
		if e.GetState() != IDLE {
			panic("expected IDLE")
		}
		e.Step()
		if e.GetCurrentFloor() != 0 {
			panic("expected floor 0")
		}
	})

	// Test 10: Multiple requests in same direction served in order
	test("test_multiple_stops_in_order", func() {
		e := NewElevator()
		e.AddRequest(5, UP)
		e.AddRequest(3, UP)
		e.AddRequest(1, UP)
		e.Step() // floor 1, DOOR_OPEN
		if e.GetCurrentFloor() != 1 {
			panic("expected floor 1")
		}
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN at 1")
		}
		e.Step() // close, MOVING_UP
		e.Step() // floor 2
		e.Step() // floor 3, DOOR_OPEN
		if e.GetCurrentFloor() != 3 {
			panic("expected floor 3")
		}
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN at 3")
		}
		e.Step() // close, MOVING_UP
		e.Step() // floor 4
		e.Step() // floor 5, DOOR_OPEN
		if e.GetCurrentFloor() != 5 {
			panic("expected floor 5")
		}
		if e.GetState() != DOOR_OPEN {
			panic("expected DOOR_OPEN at 5")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
