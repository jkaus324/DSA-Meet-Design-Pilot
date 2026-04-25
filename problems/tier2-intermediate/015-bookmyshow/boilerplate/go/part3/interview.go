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

// --- Your Design Starts Here (Part 3) ---------------------------------------
//
// This part is an integration challenge — no new features are added.
// Ensure the entire flow works end-to-end:
//   - Register theaters and shows
//   - Search movies by city
//   - Lock seats → Confirm booking (happy path)
//   - Lock seats → Expiry → seats become available again
//   - Lock seats → Release → seats become available again
//   - Concurrent-style: two users try to lock the same seat; only one succeeds
//
// All entry points are the same as Part 2.
// Focus on correctness of the state machine:
//   SeatAvailable → SeatLocked (LockSeats)
//   SeatLocked    → SeatBooked (ConfirmBooking, within TTL)
//   SeatLocked    → SeatAvailable (expiry or ReleaseLock)
//
// Entry points (same as Part 2):
//   func AddTheater / AddShow / SearchMovies
//   func GetAvailableSeats(showID string, currentTime int64) [][2]int
//   func BookSeats(bookingID, showID string, seats [][2]int, userID string, currentTime int64) bool
//   func LockSeats(showID string, seats [][2]int, userID string, ttlMinutes int, currentTime int64) string
//   func ConfirmBooking(lockID string, currentTime int64) bool
//   func ReleaseLock(lockID string, currentTime int64) bool

// -------------------------------------------------------------------------
