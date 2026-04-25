package main

// --- Data Model (given -- do not modify) ------------------------------------

type SeatStatus int

const (
	SeatAvailable SeatStatus = iota
	SeatLocked
	SeatBooked
)

type Seat struct {
	Row       int
	Col       int
	Status    SeatStatus
	BookedBy  string
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

// --- Your Design Starts Here ------------------------------------------------
//
// Design and implement a BookingSystem that:
//   1. Registers theaters and shows
//   2. Searches movies by city
//   3. Returns available seat positions for a show
//   4. Books a list of seats atomically (all-or-nothing)
//
// Think about:
//   - How do you index shows by city for fast search?
//   - What does "atomic booking" mean — what must you check before marking seats?
//
// Entry points (must exist for tests):
//   func AddTheater(theaterID, name, city string)
//   func AddShow(showID, theaterID, movie, time string, rows, cols int)
//   func SearchMovies(city string) []string
//   func GetAvailableSeats(showID string) [][2]int
//   func BookSeats(bookingID, showID string, seats [][2]int, userID string) bool

// -------------------------------------------------------------------------
