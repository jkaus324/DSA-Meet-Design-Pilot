package main

// --- Data Model (given -- do not modify) ------------------------------------

type SeatStatus int

const (
	SeatAvailable SeatStatus = iota
	SeatLocked
	SeatBooked
)

type Seat struct {
	Row      int
	Col      int
	Status   SeatStatus
	BookedBy string
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

// --- Booking System ---------------------------------------------------------
// HINT: Use map[string]Theater, map[string]Show, map[string]Booking.
// HINT: Use map[string]map[string]bool (city → set of movie titles) for search.
// HINT: When AddShow is called, also register the movie in cityMovies[city].

// type BookingSystem struct {
//     theaters    map[string]Theater
//     shows       map[string]Show
//     bookings    map[string]Booking
//     cityMovies  map[string]map[string]bool
// }

// func NewBookingSystem() *BookingSystem

// func (bs *BookingSystem) AddTheater(theaterID, name, city string)

// func (bs *BookingSystem) AddShow(showID, theaterID, movie, time string, rows, cols int)
// HINT: Initialise Seats as a rows×cols slice, all SeatAvailable.

// func (bs *BookingSystem) SearchMovies(city string) []string
// HINT: Return all unique movie titles for the city (order not required).

// func (bs *BookingSystem) GetAvailableSeats(showID string) [][2]int
// HINT: Return [row, col] pairs where Status == SeatAvailable.

// func (bs *BookingSystem) BookSeats(bookingID, showID string, seats [][2]int, userID string) bool
// HINT: First pass — verify ALL seats are SeatAvailable; return false if any is not.
// HINT: Second pass — mark each seat SeatBooked with BookedBy=userID.

// --- Global Entry Points (required by tests) --------------------------------

// var bookingSystem *BookingSystem

// func ResetBookingSystem()
// func AddTheater(theaterID, name, city string)
// func AddShow(showID, theaterID, movie, time string, rows, cols int)
// func SearchMovies(city string) []string
// func GetAvailableSeats(showID string) [][2]int
// func BookSeats(bookingID, showID string, seats [][2]int, userID string) bool
