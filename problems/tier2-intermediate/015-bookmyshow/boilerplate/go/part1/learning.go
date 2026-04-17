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

type BookingSystem struct {
	theaters   map[string]Theater
	shows      map[string]Show
	bookings   map[string]Booking
	cityMovies map[string]map[string]bool
}

func NewBookingSystem() *BookingSystem {
	// TODO: Initialise and return a BookingSystem with empty maps
	return nil
}

func (bs *BookingSystem) AddTheater(theaterID, name, city string) {
	// TODO: Store Theater{ID: theaterID, Name: name, City: city} in bs.theaters
}

func (bs *BookingSystem) AddShow(showID, theaterID, movie, timeStr string, rows, cols int) {
	// TODO: Build a rows×cols Seat grid (all SeatAvailable)
	// TODO: Store Show in bs.shows
	// TODO: Look up the theater's city and add movie to bs.cityMovies[city]
}

func (bs *BookingSystem) SearchMovies(city string) []string {
	// TODO: Return all movie titles registered under city
	// TODO: Return empty slice if city not found
	return nil
}

func (bs *BookingSystem) GetAvailableSeats(showID string) [][2]int {
	// TODO: Return [row, col] pairs where seat.Status == SeatAvailable
	return nil
}

func (bs *BookingSystem) BookSeats(bookingID, showID string, seats [][2]int, userID string) bool {
	// TODO: Return false if show not found
	// TODO: First pass: check every seat position is in bounds and SeatAvailable
	// TODO: If any check fails, return false (atomicity — do not book any)
	// TODO: Second pass: mark each seat SeatBooked, set BookedBy=userID
	// TODO: Store Booking in bs.bookings; return true
	return false
}

// --- Global Entry Points (required by tests) --------------------------------

var bookingSystem *BookingSystem

func ResetBookingSystem() {
	bookingSystem = NewBookingSystem()
}

func AddTheater(theaterID, name, city string) {
	// TODO: bookingSystem.AddTheater(...)
}

func AddShow(showID, theaterID, movie, timeStr string, rows, cols int) {
	// TODO: bookingSystem.AddShow(...)
}

func SearchMovies(city string) []string {
	// TODO: return bookingSystem.SearchMovies(city)
	return nil
}

func GetAvailableSeats(showID string) [][2]int {
	// TODO: return bookingSystem.GetAvailableSeats(showID)
	return nil
}

func BookSeats(bookingID, showID string, seats [][2]int, userID string) bool {
	// TODO: return bookingSystem.BookSeats(...)
	return false
}
