package main

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

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Design and implement a ParkingLot that:
//   1. Has multiple floors, each with spots of various sizes
//   2. Parks a vehicle in the nearest compatible spot (lowest floor first,
//      then lowest spot index)
//   3. Unparks a vehicle given a ticket ID and returns the parking fee
//   4. Reports available spot counts
//
// Compatibility:
//   MOTORCYCLE -> SMALL or larger
//   CAR        -> MEDIUM or larger
//   TRUCK      -> LARGE only
//
// Think about:
//   - How do you map vehicle types to minimum spot sizes?
//   - What data structure organizes spots by floor for nearest-first allocation?
//   - Should spot creation logic be in the lot, or in a separate factory?
//
// Entry points (must exist for tests):
//   func (lot *ParkingLot) AddSpot(floor int, size SpotSize)
//   func (lot *ParkingLot) ParkVehicle(vehicle Vehicle, entryTime int64) *Ticket
//   func (lot *ParkingLot) UnparkVehicle(ticketID string, exitTime int64) float64
//   func (lot *ParkingLot) GetAvailableSpots(size SpotSize) int
//   func (lot *ParkingLot) GetAvailableSpotsByFloor(floor int, size SpotSize) int
//
// ─────────────────────────────────────────────────────────────────────────────

type ParkingLot struct {
	// TODO: add your fields here
}

func NewParkingLot(numFloors int) *ParkingLot {
	return &ParkingLot{}
}

func (lot *ParkingLot) AddSpot(floor int, size SpotSize) {}

func (lot *ParkingLot) ParkVehicle(vehicle Vehicle, entryTime int64) *Ticket {
	return nil
}

func (lot *ParkingLot) UnparkVehicle(ticketID string, exitTime int64) float64 {
	return -1.0
}

func (lot *ParkingLot) GetAvailableSpots(size SpotSize) int {
	return 0
}

func (lot *ParkingLot) GetAvailableSpotsByFloor(floor int, size SpotSize) int {
	return 0
}
