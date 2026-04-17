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

// RideSelectionStrategy — extend RideService to support pluggable ride selection:
//   - MostVacant: selects the ride with the most available seats
//   - PreferredVehicle: selects the ride whose vehicle model matches a preference
//
// Think about:
//   - What abstraction lets you swap selection logic at runtime?
//   - How do you filter candidates (origin, dest, active, enough seats)?
//   - How does PreferredVehicleStrategy access vehicle model information?
//   - What happens when no ride matches the criteria?
//
// Entry points (all Part 1 plus):
//   (*RideService).SelectRide(passengerName, origin, dest string, seats int,
//                             strategy RideSelectionStrategy, preference string) string
//
// You also need:
//   RideSelectionStrategy interface
//   MostVacantStrategy
//   PreferredVehicleStrategy

type RideSelectionStrategy interface {
	Select(candidates []*Ride, seatsNeeded int, preference string) *Ride
}

type MostVacantStrategy struct{}
type PreferredVehicleStrategy struct{}

func (s *MostVacantStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	return nil
}

func (s *PreferredVehicleStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	return nil
}

type RideService struct {
}

func NewRideService() *RideService {
	return &RideService{}
}

func (r *RideService) AddUser(name string) {}

func (r *RideService) AddVehicle(userName, model, regNumber string) {}

func (r *RideService) OfferRide(userName, origin, dest string, seats int, vehicleRegNumber string) string {
	return ""
}

func (r *RideService) SelectRide(passengerName, origin, dest string, seats int,
	strategy RideSelectionStrategy, preference string) string {
	return ""
}

func (r *RideService) GetVehicles() map[string]Vehicle { return nil }
func (r *RideService) HasUser(name string) bool        { return false }
func (r *RideService) HasVehicle(reg string) bool      { return false }
func (r *RideService) HasRide(id string) bool          { return false }
func (r *RideService) GetUser(name string) User        { return User{} }
func (r *RideService) GetRide(id string) Ride          { return Ride{} }
