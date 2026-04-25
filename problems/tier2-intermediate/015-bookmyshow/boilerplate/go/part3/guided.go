package main

// --- Data Model (given -- do not modify) ------------------------------------

type SeatStatus int

const (
	SeatAvailable SeatStatus = iota
	SeatLocked
	SeatBooked
)

type Seat struct {
	Row        int
	Col        int
	Status     SeatStatus
	LockedBy   string
	LockExpiry int64
	BookedBy   string
}

type Show struct {
	ID        string
	TheaterID string
	Movie     string
	Time      string
	Rows      int
	Cols      int
	Seats     [][]Seat
}

type Theater struct {
	ID   string
	Name string
	City string
}

type Booking struct {
	ID            string
	ShowID        string
	UserID        string
	SeatPositions [][2]int
}

type SeatLock struct {
	ID            string
	ShowID        string
	UserID        string
	SeatPositions [][2]int
	Expiry        int64
	Confirmed     bool
	Released      bool
}

// --- Booking System (same as Part 2, full integration) -----------------------
// HINT: This part tests the complete state machine — no new types needed.
// HINT: Double-check that expireSeat is called defensively in every read path
//       (GetAvailableSeats, BookSeats, LockSeats, ConfirmBooking).
// HINT: Confirm that ConfirmBooking rejects an already-confirmed lockID
//       (idempotency: calling it twice should return false on the second call).
// HINT: Verify that ReleaseLock returns false on an already-released lock.

// --- Full class definition same as Part 2 ---

// type BookingSystem struct { ... }
// func NewBookingSystem() *BookingSystem
// ... all methods ...

// --- Global Entry Points (required by tests) --------------------------------

// func ResetBookingSystem()
// func AddTheater / AddShow / SearchMovies
// func GetAvailableSeats(showID string, currentTime int64) [][2]int
// func BookSeats(bookingID, showID string, seats [][2]int, userID string, currentTime int64) bool
// func LockSeats(showID string, seats [][2]int, userID string, ttlMinutes int, currentTime int64) string
// func ConfirmBooking(lockID string, currentTime int64) bool
// func ReleaseLock(lockID string, currentTime int64) bool
