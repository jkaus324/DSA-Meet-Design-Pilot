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

	// Test 1: Search movies by city
	test("test_search_movies_by_city", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR Phoenix", "Mumbai")
		AddShow("S1", "T1", "Inception", "18:00", 5, 10)
		movies := SearchMovies("Mumbai")
		if len(movies) != 1 {
			panic("expected 1 movie in Mumbai")
		}
		if movies[0] != "Inception" {
			panic("expected Inception")
		}
	})

	// Test 2: Search movies in unknown city returns empty
	test("test_search_movies_unknown_city", func() {
		ResetBookingSystem()
		movies := SearchMovies("Delhi")
		if len(movies) != 0 {
			panic("expected 0 movies for unknown city")
		}
	})

	// Test 3: Available seats count matches show dimensions
	test("test_available_seats_count", func() {
		ResetBookingSystem()
		AddTheater("T1", "INOX", "Bangalore")
		AddShow("S1", "T1", "Dune", "20:00", 5, 10)
		seats := GetAvailableSeats("S1")
		if len(seats) != 50 {
			panic("expected 50 available seats for 5x10 show")
		}
	})

	// Test 4: Book seats reduces available count
	test("test_book_seats_reduces_available", func() {
		ResetBookingSystem()
		AddTheater("T1", "Cinepolis", "Pune")
		AddShow("S1", "T1", "Oppenheimer", "15:00", 3, 5)
		ok := BookSeats("B1", "S1", [][2]int{{0, 0}, {0, 1}}, "user1")
		if !ok {
			panic("booking should succeed")
		}
		seats := GetAvailableSeats("S1")
		if len(seats) != 13 {
			panic("expected 13 remaining seats after booking 2 of 15")
		}
	})

	// Test 5: Cannot double-book same seat
	test("test_cannot_double_book_same_seat", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Hyderabad")
		AddShow("S1", "T1", "RRR", "10:00", 2, 2)
		ok1 := BookSeats("B1", "S1", [][2]int{{0, 0}}, "user1")
		if !ok1 {
			panic("first booking should succeed")
		}
		ok2 := BookSeats("B2", "S1", [][2]int{{0, 0}}, "user2")
		if ok2 {
			panic("second booking of same seat should fail")
		}
	})

	// Test 6: Atomic booking — if any seat unavailable, none are booked
	test("test_atomic_booking", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Chennai")
		AddShow("S1", "T1", "Vikram", "12:00", 3, 3)
		BookSeats("B1", "S1", [][2]int{{0, 0}}, "user1")
		// Try to book {0,0} (taken) and {0,1} (free) together
		ok := BookSeats("B2", "S1", [][2]int{{0, 0}, {0, 1}}, "user2")
		if ok {
			panic("atomic booking should fail when one seat is taken")
		}
		// {0,1} must still be available
		seats := GetAvailableSeats("S1")
		found := false
		for _, s := range seats {
			if s[0] == 0 && s[1] == 1 {
				found = true
			}
		}
		if !found {
			panic("seat {0,1} should still be available after failed atomic booking")
		}
	})

	// Test 7: Book seat in non-existent show returns false
	test("test_book_nonexistent_show", func() {
		ResetBookingSystem()
		ok := BookSeats("B1", "GHOST", [][2]int{{0, 0}}, "user1")
		if ok {
			panic("booking in non-existent show should return false")
		}
	})

	// Test 8: Same movie in multiple cities
	test("test_movie_in_multiple_cities", func() {
		ResetBookingSystem()
		AddTheater("T1", "PVR", "Mumbai")
		AddTheater("T2", "INOX", "Delhi")
		AddShow("S1", "T1", "Avatar", "18:00", 2, 2)
		AddShow("S2", "T2", "Avatar", "19:00", 2, 2)
		AddShow("S3", "T1", "Inception", "20:00", 2, 2)
		mumbai := SearchMovies("Mumbai")
		if len(mumbai) != 2 {
			panic("expected 2 movies in Mumbai")
		}
		delhi := SearchMovies("Delhi")
		if len(delhi) != 1 {
			panic("expected 1 movie in Delhi")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
