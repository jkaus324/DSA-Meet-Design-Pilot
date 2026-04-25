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

// RideSelectionStrategy picks one ride from a list of candidates.
type RideSelectionStrategy interface {
	Select(candidates []*Ride, seatsNeeded int, preference string) *Ride
}

// MostVacantStrategy selects the ride with the most available seats.
type MostVacantStrategy struct{}

func (s *MostVacantStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	// TODO: Find the ride with the most AvailableSeats (>= seatsNeeded)
	// TODO: Iterate through candidates, track the best (most seats)
	// TODO: Return nil if no candidate has enough seats
	return nil
}

// PreferredVehicleStrategy selects the first ride whose vehicle model matches preference.
type PreferredVehicleStrategy struct {
	vehicleStore map[string]Vehicle
}

func NewPreferredVehicleStrategy(vehicles map[string]Vehicle) *PreferredVehicleStrategy {
	return &PreferredVehicleStrategy{vehicleStore: vehicles}
}

func (s *PreferredVehicleStrategy) Select(candidates []*Ride, seatsNeeded int, preference string) *Ride {
	// TODO: Find the first ride whose vehicle model matches preference
	// TODO: Use s.vehicleStore[ride.VehicleId].Model to get the model
	// TODO: Only consider rides with AvailableSeats >= seatsNeeded
	// TODO: Return nil if no match found
	return nil
}

// RideService manages users, vehicles, and rides.
type RideService struct {
	users          map[string]User
	vehicles       map[string]Vehicle
	rides          map[string]Ride
	activeVehicles map[string]string
	rideCounter    int
}

func NewRideService() *RideService {
	return &RideService{
		users:          make(map[string]User),
		vehicles:       make(map[string]Vehicle),
		rides:          make(map[string]Ride),
		activeVehicles: make(map[string]string),
	}
}

func (r *RideService) AddUser(name string) {
	// TODO: Same as Part 1
	if _, exists := r.users[name]; exists {
		return
	}
	r.users[name] = User{Id: name, Name: name}
}

func (r *RideService) AddVehicle(userName, model, regNumber string) {
	// TODO: Same as Part 1
	if _, exists := r.users[userName]; !exists {
		return
	}
	r.vehicles[regNumber] = Vehicle{Id: regNumber, OwnerId: userName, Model: model, RegNumber: regNumber}
}

func (r *RideService) OfferRide(userName, origin, dest string, seats int, vehicleRegNumber string) string {
	// TODO: Same as Part 1
	if _, exists := r.users[userName]; !exists {
		return ""
	}
	vehicle, exists := r.vehicles[vehicleRegNumber]
	if !exists || vehicle.OwnerId != userName {
		return ""
	}
	if _, active := r.activeVehicles[vehicleRegNumber]; active {
		return ""
	}
	r.rideCounter++
	rideId := fmt.Sprintf("RIDE-%d", r.rideCounter)
	r.rides[rideId] = Ride{
		Id: rideId, DriverId: userName, VehicleId: vehicleRegNumber,
		Origin: origin, Destination: dest,
		TotalSeats: seats, AvailableSeats: seats, Active: true,
	}
	r.activeVehicles[vehicleRegNumber] = rideId
	u := r.users[userName]
	u.RidesOffered++
	r.users[userName] = u
	return rideId
}

func (r *RideService) SelectRide(passengerName, origin, dest string, seats int,
	strategy RideSelectionStrategy, preference string) string {
	// TODO: Validate passenger exists — return "" if not
	if _, exists := r.users[passengerName]; !exists {
		return ""
	}

	// TODO: Build candidate list
	var candidates []*Ride
	for id := range r.rides {
		ride := r.rides[id]
		if ride.Active && ride.Origin == origin && ride.Destination == dest &&
			ride.AvailableSeats >= seats && ride.DriverId != passengerName {
			rideCopy := ride
			candidates = append(candidates, &rideCopy)
		}
	}

	// TODO: Call strategy.Select(candidates, seats, preference)
	selected := strategy.Select(candidates, seats, preference)
	if selected == nil {
		return ""
	}

	// TODO: Decrement AvailableSeats, increment RidesTaken, return rideId
	ride := r.rides[selected.Id]
	ride.AvailableSeats -= seats
	r.rides[selected.Id] = ride

	u := r.users[passengerName]
	u.RidesTaken++
	r.users[passengerName] = u

	return selected.Id
}

func (r *RideService) GetVehicles() map[string]Vehicle { return r.vehicles }

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
