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

	// Test 1: Add elevators and verify count
	test("test_add_elevators", func() {
		sys := NewElevatorSystem()
		sys.AddElevator(1)
		sys.AddElevator(2)
		if sys.GetElevatorCount() != 2 {
			panic("expected elevator count == 2")
		}
		if sys.GetElevator(0) == nil {
			panic("expected elevator 0 to be non-nil")
		}
		if sys.GetElevator(1) == nil {
			panic("expected elevator 1 to be non-nil")
		}
	})

	// Test 2: Without strategy, request goes to first elevator
	test("test_default_dispatch", func() {
		sys := NewElevatorSystem()
		sys.AddElevator(1)
		sys.AddElevator(2)
		sys.AddRequest(5, UP)
		if sys.GetElevator(0).GetState() == IDLE {
			panic("expected elevator 0 to not be IDLE")
		}
		if sys.GetElevator(1).GetState() != IDLE {
			panic("expected elevator 1 to be IDLE")
		}
	})

	// Test 3: NearestFirst dispatches to closer elevator
	test("test_nearest_first_dispatch", func() {
		sys := NewElevatorSystem()
		nf := &NearestFirst{}
		sys.AddElevator(1)
		sys.AddElevator(2)
		sys.SetDispatchStrategy(nf)
		// Move elevator 0 to floor 5
		sys.GetElevator(0).AddRequest(5, UP)
		for i := 0; i < 5; i++ {
			sys.GetElevator(0).Step()
		}
		// Elevator 0 is at floor 5 (DOOR_OPEN), elevator 1 is at floor 0
		sys.GetElevator(0).Step() // close doors, go idle at 5
		// Request floor 2 — elevator 1 (at 0) is closer than elevator 0 (at 5)
		sys.AddRequest(2, UP)
		if sys.GetElevator(1).GetState() == IDLE {
			panic("expected elevator 1 to not be IDLE")
		}
	})

	// Test 4: LeastLoaded dispatches to elevator with fewer requests
	test("test_least_loaded_dispatch", func() {
		sys := NewElevatorSystem()
		ll := &LeastLoaded{}
		sys.AddElevator(1)
		sys.AddElevator(2)
		sys.SetDispatchStrategy(ll)
		// Give elevator 0 three requests
		sys.GetElevator(0).AddRequest(3, UP)
		sys.GetElevator(0).AddRequest(5, UP)
		sys.GetElevator(0).AddRequest(7, UP)
		// Elevator 0 has 3 pending, elevator 1 has 0
		sys.AddRequest(4, UP)
		// Should go to elevator 1 (least loaded)
		if sys.GetElevator(1).GetPendingCount() == 0 {
			panic("expected elevator 1 to have pending requests")
		}
	})

	// Test 5: Step advances all elevators
	test("test_step_all_elevators", func() {
		sys := NewElevatorSystem()
		sys.AddElevator(1)
		sys.AddElevator(2)
		sys.GetElevator(0).AddRequest(3, UP)
		sys.GetElevator(1).AddRequest(2, UP)
		sys.Step() // both move one floor
		if sys.GetElevator(0).GetCurrentFloor() != 1 {
			panic("expected elevator 0 at floor 1")
		}
		if sys.GetElevator(1).GetCurrentFloor() != 1 {
			panic("expected elevator 1 at floor 1")
		}
	})

	// Test 6: NearestFirst prefers same-direction elevator
	test("test_nearest_prefers_same_direction", func() {
		sys := NewElevatorSystem()
		nf := &NearestFirst{}
		sys.AddElevator(1)
		sys.AddElevator(2)
		sys.SetDispatchStrategy(nf)
		// Move elevator 0 to floor 3 going up (give it request for floor 8)
		sys.GetElevator(0).AddRequest(3, UP)
		sys.GetElevator(0).AddRequest(8, UP)
		for i := 0; i < 3; i++ {
			sys.GetElevator(0).Step()
		}
		// Elevator 0 at floor 3, DOOR_OPEN, still has request for 8 (moving up)
		sys.GetElevator(0).Step() // close doors, MOVING_UP toward 8
		// Elevator 1 at floor 0
		// Request floor 5 UP — elevator 0 is moving up past it, should be preferred
		sys.AddRequest(5, UP)
		// Elevator 0 should get it (moving up, will pass floor 5)
		if sys.GetElevator(0).GetPendingCount() < 2 {
			panic("expected elevator 0 to have at least 2 pending requests")
		}
	})

	// Test 7: Swapping strategy at runtime
	test("test_swap_strategy_runtime", func() {
		sys := NewElevatorSystem()
		nf := &NearestFirst{}
		ll := &LeastLoaded{}
		sys.AddElevator(1)
		sys.AddElevator(2)
		sys.SetDispatchStrategy(nf)
		sys.AddRequest(3, UP) // dispatched via NearestFirst
		sys.SetDispatchStrategy(ll)
		// Now give elevator 0 more requests so elevator 1 is least loaded
		sys.GetElevator(0).AddRequest(5, UP)
		sys.GetElevator(0).AddRequest(7, UP)
		sys.AddRequest(4, UP) // dispatched via LeastLoaded to elevator 1
		if sys.GetElevator(1).GetPendingCount() == 0 {
			panic("expected elevator 1 to have pending requests")
		}
	})

	// Test 8: GetElevator with invalid index returns nil
	test("test_invalid_elevator_index", func() {
		sys := NewElevatorSystem()
		sys.AddElevator(1)
		if sys.GetElevator(-1) != nil {
			panic("expected nil for index -1")
		}
		if sys.GetElevator(5) != nil {
			panic("expected nil for index 5")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
