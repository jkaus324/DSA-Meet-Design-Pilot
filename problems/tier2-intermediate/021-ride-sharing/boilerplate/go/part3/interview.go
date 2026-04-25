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

// RideService — extend with ride lifecycle management and statistics:
//
//   EndRide(rideId):
//     - Mark the ride as inactive
//     - Free the vehicle so it can be used again
//     - No-op if ride is already ended or doesn't exist
//
//   GetRideStats():
//     - Return per-user statistics: name, ridesOffered, ridesTaken
//
//   PrintRideStats():
//     - Print: "User: NAME — Offered: X, Taken: Y" for each user
//
// Think about:
//   - How does ending a ride free the vehicle for future rides?
//   - What if someone tries to end a ride that's already ended?
//   - How do you track per-user statistics efficiently?
//
// Entry points (all Part 1, Part 2 plus):
//   (*RideService).EndRide(rideId string)
//   (*RideService).GetRideStats() []RideStat
//   (*RideService).PrintRideStats()

// RideStat holds statistics for one user.
type RideStat struct {
	Name         string
	RidesOffered int
	RidesTaken   int
}

// RideSelectionStrategy (from Part 2)
type RideSelectionStrategy interface {
	Select(candidates []*Ride, seatsNeeded int, preference string) *Ride
}

type MostVacantStrategy struct{}

func (s *MostVacantStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	return nil
}

type PreferredVehicleStrategy struct{}

func (s *PreferredVehicleStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	return nil
}

type RideService struct {
}

func NewRideService() *RideService {
	return &RideService{}
}

func (r *RideService) AddUser(name string)                                     {}
func (r *RideService) AddVehicle(userName, model, regNumber string)            {}
func (r *RideService) OfferRide(userName, origin, dest string, seats int, vehicleRegNumber string) string {
	return ""
}
func (r *RideService) SelectRide(passengerName, origin, dest string, seats int,
	strategy RideSelectionStrategy, preference string) string {
	return ""
}
func (r *RideService) EndRide(rideId string) {}
func (r *RideService) GetRideStats() []RideStat {
	return nil
}
func (r *RideService) PrintRideStats() {}

func (r *RideService) GetVehicles() map[string]Vehicle { return nil }
func (r *RideService) HasUser(name string) bool        { return false }
func (r *RideService) HasVehicle(reg string) bool      { return false }
func (r *RideService) HasRide(id string) bool          { return false }
func (r *RideService) GetUser(name string) User        { return User{} }
func (r *RideService) GetRide(id string) Ride          { return Ride{} }
