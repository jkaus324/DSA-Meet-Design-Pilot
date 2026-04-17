package main

import (
	"fmt"
	"strconv"
)

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

func CreateSpot(spotID string, floor int, size SpotSize) ParkingSpot {
	// TODO: return ParkingSpot{SpotID: spotID, Floor: floor, Size: size, IsOccupied: false, VehicleLicensePlate: ""}
	return ParkingSpot{}
}

// ─── Compatibility Helpers ────────────────────────────────────────────────────

func GetMinSpotSize(vehicleType VehicleType) SpotSize {
	// TODO: return SMALL for MOTORCYCLE, MEDIUM for CAR, LARGE for TRUCK
	return LARGE
}

func IsCompatible(spotSize, minRequired SpotSize) bool {
	// TODO: return int(spotSize) >= int(minRequired)
	return false
}

// ─── Parking Lot ──────────────────────────────────────────────────────────────

type ParkingLot struct {
	floors        [][]ParkingSpot
	activeTickets map[string]*Ticket
	nextTicketID  int
}

func NewParkingLot(numFloors int) *ParkingLot {
	floors := make([][]ParkingSpot, numFloors)
	return &ParkingLot{
		floors:        floors,
		activeTickets: make(map[string]*Ticket),
		nextTicketID:  1,
	}
}

func (lot *ParkingLot) AddSpot(floor int, size SpotSize) {
	if floor < 0 || floor >= len(lot.floors) {
		return
	}
	spotID := fmt.Sprintf("F%dS%d", floor, len(lot.floors[floor]))
	// TODO: use CreateSpot and append to lot.floors[floor]
	_ = spotID
}

func (lot *ParkingLot) ParkVehicle(vehicle Vehicle, entryTime int64) *Ticket {
	minSize := GetMinSpotSize(vehicle.Type)
	// TODO: iterate lot.floors in order (0, 1, 2, ...)
	//   For each floor, iterate spots in index order
	//   Find the first unoccupied spot where IsCompatible(spot.Size, minSize)
	//   Mark it occupied: spot.IsOccupied = true, spot.VehicleLicensePlate = vehicle.LicensePlate
	//   Create Ticket: TicketID = "T" + strconv.Itoa(lot.nextTicketID), lot.nextTicketID++
	//   Store in lot.activeTickets, return pointer
	// TODO: return nil if no compatible spot found
	_ = minSize
	_ = strconv.Itoa
	return nil
}

func (lot *ParkingLot) UnparkVehicle(ticketID string, exitTime int64) float64 {
	ticket, ok := lot.activeTickets[ticketID]
	if !ok {
		return -1.0
	}
	// TODO: find spot in lot.floors[ticket.Floor] where spot.SpotID == ticket.SpotID
	// TODO: set spot.IsOccupied = false, spot.VehicleLicensePlate = ""
	// TODO: duration = exitTime - ticket.EntryTime
	// TODO: fee = float64(duration)  (1.0 per second for Part 1)
	// TODO: delete lot.activeTickets[ticketID]
	// TODO: return fee
	_ = ticket
	return -1.0
}

func (lot *ParkingLot) GetAvailableSpots(size SpotSize) int {
	count := 0
	for _, floor := range lot.floors {
		for _, spot := range floor {
			// TODO: if !spot.IsOccupied && spot.Size == size, count++
			_ = spot
		}
	}
	return count
}

func (lot *ParkingLot) GetAvailableSpotsByFloor(floor int, size SpotSize) int {
	if floor < 0 || floor >= len(lot.floors) {
		return 0
	}
	count := 0
	for _, spot := range lot.floors[floor] {
		// TODO: if !spot.IsOccupied && spot.Size == size, count++
		_ = spot
	}
	return count
}
