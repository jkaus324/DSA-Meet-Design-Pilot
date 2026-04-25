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

	// Test 1: Lock seats reduces available count
	test("test_lock_reduces_available", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 3, 3)
		lockID := LockSeats("S1", [][2]int{{0, 0}, {0, 1}}, "user1", 10, 1000)
		if lockID == "" {
			panic("lock should succeed")
		}
		seats := GetAvailableSeats("S1", 1000)
		if len(seats) != 7 {
			panic("expected 7 available seats after locking 2 of 9")
		}
	})

	// Test 2: Confirm booking within TTL succeeds
	test("test_confirm_booking_within_ttl", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 3, 3)
		lockID := LockSeats("S1", [][2]int{{0, 0}}, "user1", 10, 1000)
		ok := ConfirmBooking(lockID, 1100) // well within 10-min TTL
		if !ok {
			panic("confirm within TTL should succeed")
		}
		// Seat should now be booked (not available)
		seats := GetAvailableSeats("S1", 1200)
		for _, s := range seats {
			if s[0] == 0 && s[1] == 0 {
				panic("confirmed seat should not appear as available")
			}
		}
	})

	// Test 3: Confirm booking after TTL fails and frees seat
	test("test_confirm_after_expiry_fails", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 2, 2)
		// TTL = 5 minutes = 300 seconds; lock at t=0, try confirm at t=500
		lockID := LockSeats("S1", [][2]int{{0, 0}}, "user1", 5, 0)
		ok := ConfirmBooking(lockID, 500)
		if ok {
			panic("confirm after expiry should fail")
		}
		// Seat must be available again
		seats := GetAvailableSeats("S1", 600)
		found := false
		for _, s := range seats {
			if s[0] == 0 && s[1] == 0 {
				found = true
			}
		}
		if !found {
			panic("expired locked seat should become available")
		}
	})

	// Test 4: Release lock frees seat
	test("test_release_lock_frees_seat", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 2, 2)
		lockID := LockSeats("S1", [][2]int{{0, 0}}, "user1", 10, 1000)
		ok := ReleaseLock(lockID, 1050)
		if !ok {
			panic("release should succeed")
		}
		seats := GetAvailableSeats("S1", 1100)
		found := false
		for _, s := range seats {
			if s[0] == 0 && s[1] == 0 {
				found = true
			}
		}
		if !found {
			panic("released seat should be available again")
		}
	})

	// Test 5: Two users cannot lock the same seat
	test("test_two_users_cannot_lock_same_seat", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 2, 2)
		lock1 := LockSeats("S1", [][2]int{{0, 0}}, "user1", 10, 1000)
		if lock1 == "" {
			panic("first lock should succeed")
		}
		lock2 := LockSeats("S1", [][2]int{{0, 0}}, "user2", 10, 1001)
		if lock2 != "" {
			panic("second lock on same seat should fail")
		}
	})

	// Test 6: Expired lock allows another user to lock
	test("test_expired_lock_allows_relock", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 2, 2)
		// user1 locks for 5 min (300s) at t=0
		LockSeats("S1", [][2]int{{0, 0}}, "user1", 5, 0)
		// user2 tries at t=400 (after expiry)
		lock2 := LockSeats("S1", [][2]int{{0, 0}}, "user2", 10, 400)
		if lock2 == "" {
			panic("lock after expiry should succeed for new user")
		}
	})

	// Test 7: Confirm already-confirmed lock returns false
	test("test_confirm_already_confirmed", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 2, 2)
		lockID := LockSeats("S1", [][2]int{{0, 0}}, "user1", 10, 1000)
		ConfirmBooking(lockID, 1100)
		ok := ConfirmBooking(lockID, 1200)
		if ok {
			panic("confirming already-confirmed lock should return false")
		}
	})

	// Test 8: Release already-released lock returns false
	test("test_release_already_released", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 2, 2)
		lockID := LockSeats("S1", [][2]int{{0, 0}}, "user1", 10, 1000)
		ReleaseLock(lockID, 1050)
		ok := ReleaseLock(lockID, 1100)
		if ok {
			panic("releasing already-released lock should return false")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
