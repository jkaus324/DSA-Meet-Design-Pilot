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
	Id         string
	OwnerId    string
	Model      string
	RegNumber  string
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

// RideService — design and implement this struct so that:
//   1. Users can be registered by name (AddUser)
//   2. Vehicles can be registered to a user (AddVehicle)
//   3. Users can offer rides with a vehicle (OfferRide)
//   4. A vehicle cannot have multiple active rides simultaneously
//
// Think about:
//   - What data structures give O(1) lookup for users, vehicles, rides?
//   - How do you check if a vehicle already has an active ride?
//   - What happens if a user offers a ride with someone else's vehicle?
//
// Entry points (must exist for tests):
//   NewRideService() *RideService
//   (*RideService).AddUser(name string)
//   (*RideService).AddVehicle(userName, model, regNumber string)
//   (*RideService).OfferRide(userName, origin, dest string, seats int, vehicleRegNumber string) string
//   (*RideService).HasUser(name string) bool
//   (*RideService).HasVehicle(regNumber string) bool
//   (*RideService).HasRide(rideId string) bool
//   (*RideService).GetUser(name string) User
//   (*RideService).GetRide(rideId string) Ride

type RideService struct {
}

func NewRideService() *RideService {
	return &RideService{}
}

func (r *RideService) AddUser(name string) {
}

func (r *RideService) AddVehicle(userName, model, regNumber string) {
}

func (r *RideService) OfferRide(userName, origin, dest string, seats int, vehicleRegNumber string) string {
	return ""
}

func (r *RideService) HasUser(name string) bool      { return false }
func (r *RideService) HasVehicle(reg string) bool    { return false }
func (r *RideService) HasRide(rideId string) bool    { return false }
func (r *RideService) GetUser(name string) User      { return User{} }
func (r *RideService) GetRide(rideId string) Ride    { return Ride{} }
