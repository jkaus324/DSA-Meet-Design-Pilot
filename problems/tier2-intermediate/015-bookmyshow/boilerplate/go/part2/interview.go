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

// --- Your Design Starts Here (Part 2) ---------------------------------------
//
// Extend Part 1 to support a seat-locking flow:
//   - LockSeats: temporarily reserves seats for a user for ttlMinutes.
//     Returns a lockID, or "" if any seat is unavailable.
//   - ConfirmBooking: converts a valid, non-expired lock to a permanent booking.
//     Returns false if the lock is expired, already confirmed, or released.
//   - ReleaseLock: explicitly releases a lock before expiry, freeing the seats.
//   - GetAvailableSeats now accepts a currentTime parameter to expire stale locks.
//   - BookSeats also accepts currentTime (direct booking still supported).
//
// Think about:
//   - How do you handle TTL — when should a locked seat become available again?
//   - What prevents double-booking when two users lock simultaneously? (single-threaded is fine)
//
// Entry points (must exist for tests — include Part 1 entry points too):
//   func AddTheater(theaterID, name, city string)
//   func AddShow(showID, theaterID, movie, time string, rows, cols int)
//   func SearchMovies(city string) []string
//   func GetAvailableSeats(showID string, currentTime int64) [][2]int
//   func BookSeats(bookingID, showID string, seats [][2]int, userID string, currentTime int64) bool
//   func LockSeats(showID string, seats [][2]int, userID string, ttlMinutes int, currentTime int64) string
//   func ConfirmBooking(lockID string, currentTime int64) bool
//   func ReleaseLock(lockID string, currentTime int64) bool

// -------------------------------------------------------------------------
