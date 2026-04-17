package main

import "fmt"

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type VehicleType int

const (
	MOTORCYCLE VehicleType = iota
	CAR
	TRUCK
)

type SpotSize int

const (
	SMALL  SpotSize = iota
	MEDIUM
	LARGE
)

type Vehicle struct {
	LicensePlate string
	Type         VehicleType
}

type ParkingSpot struct {
	SpotID              string
	Floor               int
	Size                SpotSize
	IsOccupied          bool
	VehicleLicensePlate string
}

type Ticket struct {
	TicketID     string
	LicensePlate string
	SpotID       string
	Floor        int
	EntryTime    int64
}

// ─── Spot Factory ─────────────────────────────────────────────────────────────
// HINT: Encapsulate spot creation. Returns a ParkingSpot with IsOccupied=false.

// func CreateSpot(spotID string, floor int, size SpotSize) ParkingSpot

// ─── Compatibility Helpers ────────────────────────────────────────────────────
// HINT: Map VehicleType to minimum SpotSize. A vehicle can park in any spot
//       whose size >= its minimum. Compare using int(size) >= int(minRequired).

// func GetMinSpotSize(vehicleType VehicleType) SpotSize
// func IsCompatible(spotSize, minRequired SpotSize) bool

// ─── Parking Lot ──────────────────────────────────────────────────────────────
// HINT: Use [][]ParkingSpot for floors. Scan floor 0 first, then floor 1, etc.
//       Within a floor, scan spots in index order.
// HINT: Use map[string]Ticket to track active tickets by ID.
// HINT: Generate ticket IDs incrementally: "T1", "T2", etc.

type ParkingLot struct {
	// HINT: floors        [][]ParkingSpot
	// HINT: activeTickets map[string]Ticket
	// HINT: nextTicketID  int
}

func NewParkingLot(numFloors int) *ParkingLot {
	// HINT: initialise floors with numFloors empty slices
	return &ParkingLot{}
}

func (lot *ParkingLot) AddSpot(floor int, size SpotSize) {
	// HINT: spotID = fmt.Sprintf("F%dS%d", floor, len(lot.floors[floor]))
	// HINT: create spot via SpotFactory and append to lot.floors[floor]
	_ = fmt.Sprintf
}

func (lot *ParkingLot) ParkVehicle(vehicle Vehicle, entryTime int64) *Ticket {
	// HINT: get minSize from GetMinSpotSize(vehicle.Type)
	// HINT: iterate floors in order; within each floor iterate spots
	// HINT: find first unoccupied spot where IsCompatible(spot.Size, minSize)
	// HINT: mark occupied, create Ticket with ID "T" + nextTicketID++
	// HINT: store in activeTickets, return pointer
	// HINT: return nil if no compatible spot found
	return nil
}

func (lot *ParkingLot) UnparkVehicle(ticketID string, exitTime int64) float64 {
	// HINT: find ticket in activeTickets; return -1 if not found
	// HINT: free the spot (IsOccupied=false, clear plate)
	// HINT: duration = exitTime - ticket.EntryTime
	// HINT: for Part 1, fee = float64(duration)
	// HINT: delete ticket, return fee
	return -1.0
}

func (lot *ParkingLot) GetAvailableSpots(size SpotSize) int {
	// HINT: count all unoccupied spots across all floors matching exact size
	return 0
}

func (lot *ParkingLot) GetAvailableSpotsByFloor(floor int, size SpotSize) int {
	// HINT: count unoccupied spots on given floor matching exact size
	return 0
}
