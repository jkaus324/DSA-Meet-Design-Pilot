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

// --- Booking System (extends Part 1) -----------------------------------------
// HINT: Add a map[string]SeatLock to BookingSystem.
// HINT: Add a lockCounter int to generate lock IDs like "LOCK_1", "LOCK_2", ...

// HINT: expireSeat helper — if seat.Status==SeatLocked && currentTime >= seat.LockExpiry,
//       reset the seat to SeatAvailable, clear LockedBy and LockExpiry.

// HINT: isSeatAvailable helper — returns true if Status==SeatAvailable,
//       OR if Status==SeatLocked and currentTime >= LockExpiry.

// func (bs *BookingSystem) LockSeats(showID string, seats [][2]int,
//     userID string, ttlMinutes int, currentTime int64) string
// HINT: Check all seats available (call expireSeat first), then lock them.
//       expiry = currentTime + int64(ttlMinutes)*60
//       Store SeatLock; return the lockID.

// func (bs *BookingSystem) ConfirmBooking(lockID string, currentTime int64) bool
// HINT: Reject if lock.Confirmed || lock.Released || currentTime >= lock.Expiry.
//       If expired, call expireSeat on each seat and mark lock.Released=true.
//       On success: mark seats SeatBooked, create a Booking with ID="BK_"+lockID.

// func (bs *BookingSystem) ReleaseLock(lockID string, currentTime int64) bool
// HINT: Reject if not found, already confirmed, or already released.
//       Reset each seat to SeatAvailable, clear LockedBy/LockExpiry.

// --- Global Entry Points (required by tests) --------------------------------

// func ResetBookingSystem()
// func AddTheater(theaterID, name, city string)
// func AddShow(showID, theaterID, movie, time string, rows, cols int)
// func SearchMovies(city string) []string
// func GetAvailableSeats(showID string, currentTime int64) [][2]int
// func BookSeats(bookingID, showID string, seats [][2]int, userID string, currentTime int64) bool
// func LockSeats(showID string, seats [][2]int, userID string, ttlMinutes int, currentTime int64) string
// func ConfirmBooking(lockID string, currentTime int64) bool
// func ReleaseLock(lockID string, currentTime int64) bool
