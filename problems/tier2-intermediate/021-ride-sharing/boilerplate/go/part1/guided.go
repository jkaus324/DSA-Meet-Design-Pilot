package main

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
// HINT: Use maps for O(1) lookups:
//   - users keyed by name
//   - vehicles keyed by regNumber
//   - rides keyed by rideId
// HINT: Track active vehicles with a separate map: regNumber -> rideId
//   This makes "is this vehicle in an active ride?" an O(1) check
type RideService struct {
	// HINT: Declare your maps here
	// HINT: Use an int counter for generating ride IDs like "RIDE-1", "RIDE-2"
}

func NewRideService() *RideService {
	// HINT: Initialize all maps and the counter
	return &RideService{}
}

func (r *RideService) AddUser(name string) {
	// HINT: Check if user already exists before adding
	// HINT: Create User with RidesOffered=0, RidesTaken=0
}

func (r *RideService) AddVehicle(userName, model, regNumber string) {
	// HINT: Validate that the user exists before adding the vehicle
}

func (r *RideService) OfferRide(userName, origin, dest string, seats int, vehicleRegNumber string) string {
	// HINT: 1. Validate user exists
	// HINT: 2. Validate vehicle exists and belongs to this user
	// HINT: 3. Check vehicle doesn't have an active ride (use activeVehicles map)
	// HINT: 4. Create ride, mark vehicle as active, increment ridesOffered
	// HINT: 5. Return rideId or "" on failure
	return ""
}

// Accessors for testing
func (r *RideService) HasUser(name string) bool   { return false }
func (r *RideService) HasVehicle(reg string) bool { return false }
func (r *RideService) HasRide(id string) bool     { return false }
func (r *RideService) GetUser(name string) User   { return User{} }
func (r *RideService) GetRide(id string) Ride     { return Ride{} }
