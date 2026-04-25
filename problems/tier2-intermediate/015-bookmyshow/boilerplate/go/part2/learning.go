package main

import "fmt"

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

// --- Booking System ---------------------------------------------------------

type BookingSystem struct {
	theaters    map[string]Theater
	shows       map[string]Show
	bookings    map[string]Booking
	locks       map[string]SeatLock
	cityMovies  map[string]map[string]bool
	lockCounter int
}

func NewBookingSystem() *BookingSystem {
	// TODO: Initialise and return a BookingSystem with empty maps
	return nil
}

func (bs *BookingSystem) expireSeat(seat *Seat, currentTime int64) {
	// TODO: If seat.Status==SeatLocked && currentTime >= seat.LockExpiry,
	//       reset to SeatAvailable, clear LockedBy and LockExpiry
}

func (bs *BookingSystem) isSeatAvailable(seat *Seat, currentTime int64) bool {
	// TODO: Return true if Status==SeatAvailable
	// TODO: Also return true if Status==SeatLocked && currentTime >= seat.LockExpiry
	return false
}

func (bs *BookingSystem) AddTheater(theaterID, name, city string) {
	// TODO: Store Theater in bs.theaters
}

func (bs *BookingSystem) AddShow(showID, theaterID, movie, timeStr string, rows, cols int) {
	// TODO: Build rows×cols Seat grid (all SeatAvailable)
	// TODO: Store Show; register movie in cityMovies[city]
}

func (bs *BookingSystem) SearchMovies(city string) []string {
	// TODO: Return movie titles for the city
	return nil
}

func (bs *BookingSystem) GetAvailableSeats(showID string, currentTime int64) [][2]int {
	// TODO: For each seat, call expireSeat first, then collect if SeatAvailable
	return nil
}

func (bs *BookingSystem) BookSeats(bookingID, showID string, seats [][2]int, userID string, currentTime int64) bool {
	// TODO: Atomic check (expireSeat + SeatAvailable) for all seats; then book all
	return false
}

func (bs *BookingSystem) LockSeats(showID string, seats [][2]int, userID string, ttlMinutes int, currentTime int64) string {
	// TODO: Check all seats available (call expireSeat first)
	// TODO: bs.lockCounter++; lockID = fmt.Sprintf("LOCK_%d", bs.lockCounter)
	// TODO: expiry = currentTime + int64(ttlMinutes)*60
	// TODO: Mark each seat SeatLocked; store SeatLock; return lockID
	_ = fmt.Sprintf // use fmt for ID generation
	return ""
}

func (bs *BookingSystem) ConfirmBooking(lockID string, currentTime int64) bool {
	// TODO: Return false if lock not found, already confirmed, or already released
	// TODO: If currentTime >= lock.Expiry: call expireSeat for each seat,
	//       set lock.Released=true, return false
	// TODO: Mark each seat SeatBooked; create Booking with ID="BK_"+lockID
	// TODO: Set lock.Confirmed=true; return true
	return false
}

func (bs *BookingSystem) ReleaseLock(lockID string, currentTime int64) bool {
	// TODO: Return false if lock not found, confirmed, or released
	// TODO: Reset each seat to SeatAvailable; clear LockedBy and LockExpiry
	// TODO: Set lock.Released=true; return true
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

func GetAvailableSeats(showID string, currentTime int64) [][2]int {
	// TODO: return bookingSystem.GetAvailableSeats(showID, currentTime)
	return nil
}

func BookSeats(bookingID, showID string, seats [][2]int, userID string, currentTime int64) bool {
	// TODO: return bookingSystem.BookSeats(...)
	return false
}

func LockSeats(showID string, seats [][2]int, userID string, ttlMinutes int, currentTime int64) string {
	// TODO: return bookingSystem.LockSeats(...)
	return ""
}

func ConfirmBooking(lockID string, currentTime int64) bool {
	// TODO: return bookingSystem.ConfirmBooking(lockID, currentTime)
	return false
}

func ReleaseLock(lockID string, currentTime int64) bool {
	// TODO: return bookingSystem.ReleaseLock(lockID, currentTime)
	return false
}
