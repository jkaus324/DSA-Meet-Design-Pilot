// Ride-sharing service — onboarding, selection strategies, end + statistics (Go port).
package main

import "strconv"

type RideOp struct {
	kind string
	s1   string
	s2   string
	s3   string
	s4   string
	i1   int
	i2   int
}

type rsUser struct {
	id           string
	name         string
	ridesOffered int
	ridesTaken   int
}

type vehicle struct {
	id        string
	ownerID   string
	model     string
	regNumber string
}

type ride struct {
	id             string
	driverID       string
	vehicleID      string
	origin         string
	destination    string
	totalSeats     int
	availableSeats int
	active         bool
}

type selectionStrategy interface {
	selectRide(candidates []*ride, seatsNeeded int, preference string) *ride
}

type mostVacantStrategy struct{}

func (mostVacantStrategy) selectRide(candidates []*ride, seatsNeeded int, preference string) *ride {
	var best *ride
	for _, r := range candidates {
		if r.availableSeats >= seatsNeeded {
			if best == nil || r.availableSeats > best.availableSeats {
				best = r
			}
		}
	}
	return best
}

type preferredVehicleStrategy struct {
	vehicleStore map[string]*vehicle
}

func (p preferredVehicleStrategy) selectRide(candidates []*ride, seatsNeeded int, preference string) *ride {
	for _, r := range candidates {
		if r.availableSeats >= seatsNeeded {
			v := p.vehicleStore[r.vehicleID]
			if v != nil && v.model == preference {
				return r
			}
		}
	}
	return nil
}

type rideService struct {
	users          map[string]*rsUser
	vehicles       map[string]*vehicle
	rides          map[string]*ride
	rideOrder      []string
	activeVehicles map[string]string
	rideCounter    int
}

func newRideService() *rideService {
	return &rideService{
		users:          map[string]*rsUser{},
		vehicles:       map[string]*vehicle{},
		rides:          map[string]*ride{},
		rideOrder:      []string{},
		activeVehicles: map[string]string{},
		rideCounter:    0,
	}
}

func (s *rideService) addUser(name string) {
	if _, ok := s.users[name]; ok {
		return
	}
	s.users[name] = &rsUser{id: name, name: name}
}

func (s *rideService) addVehicle(userName, model, regNumber string) {
	if _, ok := s.users[userName]; !ok {
		return
	}
	s.vehicles[regNumber] = &vehicle{id: regNumber, ownerID: userName, model: model, regNumber: regNumber}
}

func (s *rideService) offerRide(userName, origin, dest string, seats int, vehicleRegNumber string) string {
	if _, ok := s.users[userName]; !ok {
		return ""
	}
	v, ok := s.vehicles[vehicleRegNumber]
	if !ok {
		return ""
	}
	if v.ownerID != userName {
		return ""
	}
	if _, ok := s.activeVehicles[vehicleRegNumber]; ok {
		return ""
	}
	s.rideCounter++
	rideID := "RIDE-" + strconv.Itoa(s.rideCounter)
	s.rides[rideID] = &ride{
		id:             rideID,
		driverID:       userName,
		vehicleID:      vehicleRegNumber,
		origin:         origin,
		destination:    dest,
		totalSeats:     seats,
		availableSeats: seats,
		active:         true,
	}
	s.rideOrder = append(s.rideOrder, rideID)
	s.activeVehicles[vehicleRegNumber] = rideID
	s.users[userName].ridesOffered++
	return rideID
}

func (s *rideService) selectRide(passengerName, origin, dest string, seats int, strategy selectionStrategy) string {
	if _, ok := s.users[passengerName]; !ok {
		return ""
	}
	candidates := []*ride{}
	for _, rid := range s.rideOrder {
		r := s.rides[rid]
		if r.active && r.origin == origin && r.destination == dest && r.availableSeats >= seats && r.driverID != passengerName {
			candidates = append(candidates, r)
		}
	}
	// preference is passed via the strategy for preferred-vehicle; the
	// MostVacant strategy ignores it. Mirror python: strategy.select takes it.
	return s.selectWith(passengerName, candidates, seats, strategy, "")
}

func (s *rideService) selectRideWithPref(passengerName, origin, dest string, seats int, strategy selectionStrategy, preference string) string {
	if _, ok := s.users[passengerName]; !ok {
		return ""
	}
	candidates := []*ride{}
	for _, rid := range s.rideOrder {
		r := s.rides[rid]
		if r.active && r.origin == origin && r.destination == dest && r.availableSeats >= seats && r.driverID != passengerName {
			candidates = append(candidates, r)
		}
	}
	return s.selectWith(passengerName, candidates, seats, strategy, preference)
}

func (s *rideService) selectWith(passengerName string, candidates []*ride, seats int, strategy selectionStrategy, preference string) string {
	selected := strategy.selectRide(candidates, seats, preference)
	if selected != nil {
		selected.availableSeats -= seats
		s.users[passengerName].ridesTaken++
		return selected.id
	}
	return ""
}

func (s *rideService) endRide(rideID string) {
	r, ok := s.rides[rideID]
	if !ok {
		return
	}
	if !r.active {
		return
	}
	r.active = false
	if _, ok := s.activeVehicles[r.vehicleID]; ok {
		delete(s.activeVehicles, r.vehicleID)
	}
}

func (s *rideService) hasUser(name string) bool    { _, ok := s.users[name]; return ok }
func (s *rideService) hasVehicle(reg string) bool  { _, ok := s.vehicles[reg]; return ok }
func (s *rideService) hasRide(rid string) bool     { _, ok := s.rides[rid]; return ok }

func ride_simulate(ops []RideOp) []string {
	out := []string{}
	svc := newRideService()
	rideSlots := make([]string, 32)
	for _, op := range ops {
		switch op.kind {
		case "new":
			svc = newRideService()
			rideSlots = make([]string, 32)
			out = append(out, "ok")
		case "add_user":
			svc.addUser(op.s1)
			out = append(out, "ok")
		case "add_veh":
			svc.addVehicle(op.s1, op.s2, op.s3)
			out = append(out, "ok")
		case "offer":
			rid := svc.offerRide(op.s1, op.s2, op.s3, op.i1, op.s4)
			if op.i2 >= 0 && op.i2 < len(rideSlots) {
				rideSlots[op.i2] = rid
			}
			out = append(out, rid)
		case "ride_active":
			rid := rideSlots[op.i2]
			if svc.hasRide(rid) && svc.rides[rid].active {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		case "ride_origin":
			rid := rideSlots[op.i2]
			if svc.hasRide(rid) {
				out = append(out, svc.rides[rid].origin)
			} else {
				out = append(out, "")
			}
		case "ride_dest":
			rid := rideSlots[op.i2]
			if svc.hasRide(rid) {
				out = append(out, svc.rides[rid].destination)
			} else {
				out = append(out, "")
			}
		case "ride_total":
			rid := rideSlots[op.i2]
			if svc.hasRide(rid) {
				out = append(out, strconv.Itoa(svc.rides[rid].totalSeats))
			} else {
				out = append(out, "-1")
			}
		case "ride_avail":
			rid := rideSlots[op.i2]
			if svc.hasRide(rid) {
				out = append(out, strconv.Itoa(svc.rides[rid].availableSeats))
			} else {
				out = append(out, "-1")
			}
		case "ride_driver":
			rid := rideSlots[op.i2]
			if svc.hasRide(rid) {
				out = append(out, svc.rides[rid].driverID)
			} else {
				out = append(out, "")
			}
		case "select_mv":
			rid := svc.selectRide(op.s1, op.s2, op.s3, op.i1, mostVacantStrategy{})
			if op.i2 >= 0 && op.i2 < len(rideSlots) {
				rideSlots[op.i2] = rid
			}
			out = append(out, rid)
		case "select_pv":
			rid := svc.selectRideWithPref(op.s1, op.s2, op.s3, op.i1, preferredVehicleStrategy{vehicleStore: svc.vehicles}, op.s4)
			if op.i2 >= 0 && op.i2 < len(rideSlots) {
				rideSlots[op.i2] = rid
			}
			out = append(out, rid)
		case "end":
			svc.endRide(rideSlots[op.i2])
			out = append(out, "ok")
		case "end_id":
			svc.endRide(op.s1)
			out = append(out, "ok")
		case "user_offered":
			if svc.hasUser(op.s1) {
				out = append(out, strconv.Itoa(svc.users[op.s1].ridesOffered))
			} else {
				out = append(out, "0")
			}
		case "user_taken":
			if svc.hasUser(op.s1) {
				out = append(out, strconv.Itoa(svc.users[op.s1].ridesTaken))
			} else {
				out = append(out, "0")
			}
		case "has_user":
			if svc.hasUser(op.s1) {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		case "has_vehicle":
			if svc.hasVehicle(op.s1) {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		case "has_ride":
			if svc.hasRide(rideSlots[op.i2]) {
				out = append(out, "yes")
			} else {
				out = append(out, "no")
			}
		default:
			out = append(out, "unknown:"+op.kind)
		}
	}
	return out
}
