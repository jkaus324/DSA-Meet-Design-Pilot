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

// ─── Include your Part 1 ParkingLot logic here ───────────────────────────────

// ─── Pricing Strategy Interface ───────────────────────────────────────────────
// HINT: The strategy receives duration in seconds and returns a fee (float64).
//       The parking lot calls strategy.CalculateFee(duration) during unpark.

type PricingStrategy interface {
	CalculateFee(durationSeconds int64) float64
}

// ─── FlatRate ─────────────────────────────────────────────────────────────────
// HINT: Constructor takes a fixed fee. CalculateFee ignores duration.

type FlatRate struct {
	// HINT: fee float64
}

func NewFlatRate(fee float64) *FlatRate { return &FlatRate{} }

func (f *FlatRate) CalculateFee(durationSeconds int64) float64 {
	// HINT: return f.fee
	return 0.0
}

// ─── Hourly ───────────────────────────────────────────────────────────────────
// HINT: Constructor takes ratePerHour. Use math.Ceil(seconds / 3600.0) for hours.

type Hourly struct {
	// HINT: ratePerHour float64
}

func NewHourly(rate float64) *Hourly { return &Hourly{} }

func (h *Hourly) CalculateFee(durationSeconds int64) float64 {
	// HINT: hours := math.Ceil(float64(durationSeconds) / 3600.0)
	// HINT: return h.ratePerHour * hours
	return 0.0
}

// ─── Tiered ───────────────────────────────────────────────────────────────────
// HINT: base (first hour), mid (hours 1-3), high (hours 3+).
//       hours <= 1 → baseRate
//       hours <= 3 → baseRate + midRate * (hours - 1)
//       hours > 3  → baseRate + midRate * 2 + highRate * (hours - 3)

type Tiered struct {
	// HINT: baseRate, midRate, highRate float64
}

func NewTiered(base, mid, high float64) *Tiered { return &Tiered{} }

func (t *Tiered) CalculateFee(durationSeconds int64) float64 {
	// HINT: hours := math.Ceil(float64(durationSeconds) / 3600.0)
	return 0.0
}

// ─── Parking Lot (extended from Part 1) ──────────────────────────────────────

type ParkingLot struct {
	// HINT: floors        [][]ParkingSpot
	// HINT: activeTickets map[string]*Ticket
	// HINT: gates         []Gate
	// HINT: strategy      PricingStrategy
	// HINT: nextTicketID  int
}

func NewParkingLot(numFloors int) *ParkingLot { return &ParkingLot{} }

func (lot *ParkingLot) AddSpot(floor int, size SpotSize) {}

func (lot *ParkingLot) SetPricingStrategy(strategy PricingStrategy) {
	// HINT: lot.strategy = strategy
}

func (lot *ParkingLot) AddGate(gateID string, gateType GateType) {
	// HINT: append Gate{gateID, gateType} to lot.gates
}

func (lot *ParkingLot) GetGates(gateType GateType) []string {
	// HINT: collect gateID from lot.gates where gate.GateType == gateType
	return nil
}

func (lot *ParkingLot) ParkVehicle(vehicle Vehicle, entryTime int64, gateID string) *Ticket {
	// HINT: same nearest-compatible-spot logic as Part 1
	// HINT: set ticket.EntryGateID = gateID
	return nil
}

func (lot *ParkingLot) UnparkVehicle(ticketID string, exitTime int64, gateID string) float64 {
	// HINT: same as Part 1 but set ticket.ExitGateID = gateID
	// HINT: use lot.strategy.CalculateFee(duration) if strategy != nil
	//       else default to float64(duration)
	return -1.0
}

func (lot *ParkingLot) GetAvailableSpots(size SpotSize) int          { return 0 }
func (lot *ParkingLot) GetAvailableSpotsByFloor(floor int, size SpotSize) int { return 0 }
