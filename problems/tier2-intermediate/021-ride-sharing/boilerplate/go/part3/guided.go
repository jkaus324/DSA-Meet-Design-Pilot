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

// TODO: Implement MostVacantStrategy and PreferredVehicleStrategy (same as Part 2)
type MostVacantStrategy struct{}

func (s *MostVacantStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	return nil
}

type PreferredVehicleStrategy struct {
	// HINT: Store reference to vehicles map
}

func NewPreferredVehicleStrategy(vehicles map[string]Vehicle) *PreferredVehicleStrategy {
	return &PreferredVehicleStrategy{}
}

func (s *PreferredVehicleStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	return nil
}

// RideService — adds ride lifecycle and statistics to Part 2.
// HINT: EndRide should:
//   1. Validate ride exists
//   2. Check if ride is still active (no-op if already ended)
//   3. Set ride.Active = false
//   4. Remove vehicle from activeVehicles map (frees it for future rides)
//
// HINT: GetRideStats iterates the users map and returns RideStat entries
// HINT: PrintRideStats formats as: "User: NAME — Offered: X, Taken: Y"
type RideService struct {
	// HINT: Same fields as Part 2
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
	// TODO: Same as Part 2
	return ""
}

func (r *RideService) EndRide(rideId string) {
	// HINT: No-op if ride doesn't exist or already inactive
	// HINT: Set ride.Active = false, write back to map
	// HINT: Delete vehicle from activeVehicles
}

func (r *RideService) GetRideStats() []RideStat {
	// HINT: Iterate r.users and build a RideStat for each
	return nil
}

func (r *RideService) PrintRideStats() {
	// HINT: For each user, print: "User: NAME — Offered: X, Taken: Y"
}

func (r *RideService) GetVehicles() map[string]Vehicle { return nil }
func (r *RideService) HasUser(name string) bool        { return false }
func (r *RideService) HasVehicle(reg string) bool      { return false }
func (r *RideService) HasRide(id string) bool          { return false }
func (r *RideService) GetUser(name string) User        { return User{} }
func (r *RideService) GetRide(id string) Ride          { return Ride{} }
