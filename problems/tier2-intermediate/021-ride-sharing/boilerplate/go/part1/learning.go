package main

import "fmt"

// User holds registration data and ride statistics for a platform participant.
type User struct {
	Id           string
	Name         string
	RidesOffered int
	RidesTaken   int
}

// Vehicle represents a vehicle registered on the platform.
type Vehicle struct {
	Id        string
	OwnerId   string
	Model     string
	RegNumber string
}

// Ride represents a single ride offered on the platform.
type Ride struct {
	Id             string
	DriverId       string
	VehicleId      string
	Origin         string
	Destination    string
	TotalSeats     int
	AvailableSeats int
	Active         bool
}

// RideService manages users, vehicles, and rides.
type RideService struct {
	users          map[string]User
	vehicles       map[string]Vehicle // keyed by regNumber
	rides          map[string]Ride    // keyed by rideId
	activeVehicles map[string]string  // regNumber -> rideId
	rideCounter    int
}

func NewRideService() *RideService {
	return &RideService{
		users:          make(map[string]User),
		vehicles:       make(map[string]Vehicle),
		rides:          make(map[string]Ride),
		activeVehicles: make(map[string]string),
		rideCounter:    0,
	}
}

func (r *RideService) AddUser(name string) {
	// TODO: Check if user already exists (r.users[name] exists) — skip if so
	// TODO: Create User{Id: name, Name: name, RidesOffered: 0, RidesTaken: 0}
	// TODO: Store in r.users keyed by name
}

func (r *RideService) AddVehicle(userName, model, regNumber string) {
	// TODO: Check if user exists in r.users — skip if not
	// TODO: Create Vehicle{Id: regNumber, OwnerId: userName, Model: model, RegNumber: regNumber}
	// TODO: Store in r.vehicles keyed by regNumber
}

func (r *RideService) OfferRide(userName, origin, dest string, seats int, vehicleRegNumber string) string {
	// TODO: Validate user exists — return "" if not
	// TODO: Validate vehicle exists — return "" if not
	// TODO: Validate vehicle.OwnerId == userName — return "" if not
	// TODO: Check activeVehicles — return "" if vehicle already has an active ride
	// TODO: Generate rideId: fmt.Sprintf("RIDE-%d", r.rideCounter+1), then increment
	r.rideCounter++
	rideId := fmt.Sprintf("RIDE-%d", r.rideCounter)
	_ = rideId
	// TODO: Create Ride{Id: rideId, DriverId: userName, VehicleId: vehicleRegNumber,
	//         Origin: origin, Destination: dest, TotalSeats: seats, AvailableSeats: seats, Active: true}
	// TODO: Store ride in r.rides
	// TODO: Mark vehicle as active: r.activeVehicles[vehicleRegNumber] = rideId
	// TODO: Increment r.users[userName].RidesOffered (remember to write back to map)
	// TODO: Return rideId
	return ""
}

// Accessors for testing
func (r *RideService) HasUser(name string) bool {
	_, ok := r.users[name]
	return ok
}

func (r *RideService) HasVehicle(reg string) bool {
	_, ok := r.vehicles[reg]
	return ok
}

func (r *RideService) HasRide(id string) bool {
	_, ok := r.rides[id]
	return ok
}

func (r *RideService) GetUser(name string) User {
	return r.users[name]
}

func (r *RideService) GetRide(id string) Ride {
	return r.rides[id]
}
