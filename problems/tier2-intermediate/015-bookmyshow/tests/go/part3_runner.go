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

	// Test 1: Full happy path — lock → confirm → seat is booked
	test("test_full_lock_confirm_flow", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 4, 4)
		lockID := LockSeats("S1", [][2]int{{1, 1}, {1, 2}}, "user1", 10, 0)
		if lockID == "" {
			panic("lock should succeed")
		}
		ok := ConfirmBooking(lockID, 100)
		if !ok {
			panic("confirm should succeed within TTL")
		}
		seats := GetAvailableSeats("S1", 200)
		for _, s := range seats {
			if (s[0] == 1 && s[1] == 1) || (s[0] == 1 && s[1] == 2) {
				panic("confirmed seats should not appear available")
			}
		}
	})

	// Test 2: Lock expiry restores seats; another user can then book directly
	test("test_lock_expire_then_direct_book", func() {
		ResetBookingSystem()
		AddTheater("T1", "INOX", "Delhi")
		AddShow("S1", "T1", "Dune", "20:00", 3, 3)
		// user1 locks {0,0} for 5 min at t=0; does NOT confirm
		LockSeats("S1", [][2]int{{0, 0}}, "user1", 5, 0)
		// user2 books at t=400 (after 300s expiry)
		ok := BookSeats("B1", "S1", [][2]int{{0, 0}}, "user2", 400)
		if !ok {
			panic("direct booking after lock expiry should succeed")
		}
	})

	// Test 3: Release → another user can immediately lock
	test("test_release_then_relock", func() {
		ResetBookingSystem()
		AddTheater("T1", "Cinepolis", "Bangalore")
		AddShow("S1", "T1", "RRR", "15:00", 3, 3)
		lock1 := LockSeats("S1", [][2]int{{2, 2}}, "user1", 10, 0)
		ReleaseLock(lock1, 50)
		lock2 := LockSeats("S1", [][2]int{{2, 2}}, "user2", 10, 100)
		if lock2 == "" {
			panic("relock after release should succeed")
		}
	})

	// Test 4: Multi-show independence — booking one show does not affect another
	test("test_shows_are_independent", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "14:00", 2, 2)
		AddShow("S2", "T1", "Inception", "18:00", 2, 2)
		BookSeats("B1", "S1", [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}, "user1", 0)
		seats := GetAvailableSeats("S2", 0)
		if len(seats) != 4 {
			panic("second show should still have all seats available")
		}
	})

	// Test 5: GetAvailableSeats reflects real-time expiry without explicit call
	test("test_available_seats_reflects_ttl_expiry", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Hyderabad")
		AddShow("S1", "T1", "Bahubali", "10:00", 2, 2)
		LockSeats("S1", [][2]int{{0, 0}}, "user1", 5, 0) // 5 min = 300s
		// Check before expiry
		before := GetAvailableSeats("S1", 100)
		// Check after expiry
		after := GetAvailableSeats("S1", 400)
		if len(before) != 3 {
			panic("seat should be locked (3 available) before expiry")
		}
		if len(after) != 4 {
			panic("seat should be released (4 available) after expiry")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
