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

type GateType int

const (
	ENTRY GateType = iota
	EXIT
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
	EntryGateID  string
	ExitGateID   string
}

type Gate struct {
	GateID   string
	GateType GateType
}

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Extend your Part 1 ParkingLot to support:
//   1. Pluggable pricing strategies: FlatRate, Hourly, Tiered
//   2. Entry/exit gate registration and tracking
//   3. Gate IDs recorded on tickets
//
// Think about:
//   - How do you define a pricing interface so new strategies can be
//     added without modifying the lot?
//   - How does the Tiered strategy calculate fees across brackets?
//   - Should gates validate their type?
//
// Entry points (must exist for tests):
//   func (lot *ParkingLot) SetPricingStrategy(strategy PricingStrategy)
//   func (lot *ParkingLot) AddGate(gateID string, gateType GateType)
//   func (lot *ParkingLot) ParkVehicle(vehicle Vehicle, entryTime int64, gateID string) *Ticket
//   func (lot *ParkingLot) UnparkVehicle(ticketID string, exitTime int64, gateID string) float64
//   func (lot *ParkingLot) GetGates(gateType GateType) []string
//
// ─────────────────────────────────────────────────────────────────────────────

type PricingStrategy interface {
	CalculateFee(durationSeconds int64) float64
}

type ParkingLot struct {
	// TODO: add your fields here
}

func NewParkingLot(numFloors int) *ParkingLot {
	return &ParkingLot{}
}

func (lot *ParkingLot) AddSpot(floor int, size SpotSize) {}

func (lot *ParkingLot) SetPricingStrategy(strategy PricingStrategy) {}

func (lot *ParkingLot) AddGate(gateID string, gateType GateType) {}

func (lot *ParkingLot) GetGates(gateType GateType) []string { return nil }

func (lot *ParkingLot) ParkVehicle(vehicle Vehicle, entryTime int64, gateID string) *Ticket {
	return nil
}

func (lot *ParkingLot) UnparkVehicle(ticketID string, exitTime int64, gateID string) float64 {
	return -1.0
}

func (lot *ParkingLot) GetAvailableSpots(size SpotSize) int { return 0 }

func (lot *ParkingLot) GetAvailableSpotsByFloor(floor int, size SpotSize) int { return 0 }

// ─── Concrete Strategies ──────────────────────────────────────────────────────

type FlatRate struct {
	fee float64
}

func NewFlatRate(fee float64) *FlatRate { return &FlatRate{fee: fee} }

func (f *FlatRate) CalculateFee(durationSeconds int64) float64 {
	return 0.0
}

type Hourly struct {
	ratePerHour float64
}

func NewHourly(rate float64) *Hourly { return &Hourly{ratePerHour: rate} }

func (h *Hourly) CalculateFee(durationSeconds int64) float64 {
	return 0.0
}

type Tiered struct {
	baseRate float64 // first hour
	midRate  float64 // hours 1-3
	highRate float64 // hours 3+
}

func NewTiered(base, mid, high float64) *Tiered {
	return &Tiered{baseRate: base, midRate: mid, highRate: high}
}

func (t *Tiered) CalculateFee(durationSeconds int64) float64 {
	return 0.0
}
