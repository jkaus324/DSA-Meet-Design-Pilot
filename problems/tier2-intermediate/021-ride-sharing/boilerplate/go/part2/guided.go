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

// RideSelectionStrategy picks one ride from a list of candidates.
// HINT: Each strategy picks one ride from a list of candidates.
// HINT: The preference string is used by PreferredVehicleStrategy to match model name.
type RideSelectionStrategy interface {
	Select(candidates []*Ride, seatsNeeded int, preference string) *Ride
}

// MostVacantStrategy selects the ride with the most available seats.
// TODO: Implement Select
//   HINT: Iterate candidates, find the one with max AvailableSeats >= seatsNeeded
type MostVacantStrategy struct{}

func (s *MostVacantStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	return nil
}

// PreferredVehicleStrategy selects the first ride whose vehicle model matches preference.
// HINT: Needs access to vehicle data to resolve VehicleId -> Model
// HINT: Pass a reference to the vehicles map in the constructor
type PreferredVehicleStrategy struct {
	// HINT: Store a reference to the vehicles map
}

func NewPreferredVehicleStrategy(vehicles map[string]Vehicle) *PreferredVehicleStrategy {
	return &PreferredVehicleStrategy{}
}

func (s *PreferredVehicleStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	// HINT: Find the first ride whose vehicle model == preference
	// HINT: Only consider rides with AvailableSeats >= seatsNeeded
	return nil
}

// RideService manages users, vehicles, and rides.
// HINT: SelectRide should:
//   1. Filter all active rides matching origin + destination
//   2. Exclude rides where DriverId == passengerName
//   3. Pass candidates to strategy.Select()
//   4. If selected, decrement AvailableSeats and increment passenger's RidesTaken
//   5. Return rideId or "" if no match
type RideService struct {
	// HINT: Same fields as Part 1
}

func NewRideService() *RideService {
	return &RideService{}
}

func (r *RideService) AddUser(name string) {
	// TODO: Same as Part 1
}

func (r *RideService) AddVehicle(userName, model, regNumber string) {
	// TODO: Same as Part 1
}

func (r *RideService) OfferRide(userName, origin, dest string, seats int, vehicleRegNumber string) string {
	// TODO: Same as Part 1
	return ""
}

func (r *RideService) SelectRide(passengerName, origin, dest string, seats int,
	strategy RideSelectionStrategy, preference string) string {
	// HINT: Validate passenger exists
	// HINT: Build candidate list (active, origin match, dest match, seats enough, not own ride)
	// HINT: Call strategy.Select(candidates, seats, preference)
	// HINT: If selected: decrement AvailableSeats, increment RidesTaken, return rideId
	return ""
}

func (r *RideService) GetVehicles() map[string]Vehicle { return nil }
func (r *RideService) HasUser(name string) bool        { return false }
func (r *RideService) HasVehicle(reg string) bool      { return false }
func (r *RideService) HasRide(id string) bool          { return false }
func (r *RideService) GetUser(name string) User        { return User{} }
func (r *RideService) GetRide(id string) Ride          { return Ride{} }
