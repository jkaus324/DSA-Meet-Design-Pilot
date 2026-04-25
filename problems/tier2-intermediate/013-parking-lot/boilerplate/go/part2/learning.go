package main

import (
	"fmt"
	"math"
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

// ─── Spot Factory (from Part 1) ───────────────────────────────────────────────

func CreateSpot(spotID string, floor int, size SpotSize) ParkingSpot {
	return ParkingSpot{SpotID: spotID, Floor: floor, Size: size, IsOccupied: false}
}

// ─── Helpers (from Part 1) ────────────────────────────────────────────────────

func GetMinSpotSize(vehicleType VehicleType) SpotSize {
	switch vehicleType {
	case MOTORCYCLE:
		return SMALL
	case CAR:
		return MEDIUM
	case TRUCK:
		return LARGE
	}
	return LARGE
}

func IsCompatible(spotSize, minRequired SpotSize) bool {
	return int(spotSize) >= int(minRequired)
}

// ─── Pricing Strategy Interface ───────────────────────────────────────────────

type PricingStrategy interface {
	CalculateFee(durationSeconds int64) float64
}

// ─── FlatRate Strategy ────────────────────────────────────────────────────────

type FlatRate struct {
	fee float64
}

func NewFlatRate(fee float64) *FlatRate { return &FlatRate{fee: fee} }

func (f *FlatRate) CalculateFee(durationSeconds int64) float64 {
	// TODO: return f.fee regardless of duration
	return 0.0
}

// ─── Hourly Strategy ──────────────────────────────────────────────────────────

type Hourly struct {
	ratePerHour float64
}

func NewHourly(rate float64) *Hourly { return &Hourly{ratePerHour: rate} }

func (h *Hourly) CalculateFee(durationSeconds int64) float64 {
	// TODO: hours := math.Ceil(float64(durationSeconds) / 3600.0)
	// TODO: return h.ratePerHour * hours
	_ = math.Ceil
	return 0.0
}

// ─── Tiered Strategy ──────────────────────────────────────────────────────────

type Tiered struct {
	baseRate float64
	midRate  float64
	highRate float64
}

func NewTiered(base, mid, high float64) *Tiered {
	return &Tiered{baseRate: base, midRate: mid, highRate: high}
}

func (t *Tiered) CalculateFee(durationSeconds int64) float64 {
	hours := math.Ceil(float64(durationSeconds) / 3600.0)
	// TODO: if hours <= 1, return t.baseRate
	// TODO: if hours <= 3, return t.baseRate + t.midRate * (hours - 1)
	// TODO: if hours > 3,  return t.baseRate + t.midRate * 2 + t.highRate * (hours - 3)
	_ = hours
	return 0.0
}

// ─── Parking Lot ──────────────────────────────────────────────────────────────

type ParkingLot struct {
	floors        [][]ParkingSpot
	activeTickets map[string]*Ticket
	gates         []Gate
	strategy      PricingStrategy
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
	lot.floors[floor] = append(lot.floors[floor], CreateSpot(spotID, floor, size))
}

func (lot *ParkingLot) SetPricingStrategy(strategy PricingStrategy) {
	// TODO: lot.strategy = strategy
}

func (lot *ParkingLot) AddGate(gateID string, gateType GateType) {
	// TODO: lot.gates = append(lot.gates, Gate{GateID: gateID, GateType: gateType})
}

func (lot *ParkingLot) GetGates(gateType GateType) []string {
	var result []string
	for _, g := range lot.gates {
		// TODO: if g.GateType == gateType, append g.GateID to result
		_ = g
	}
	return result
}

func (lot *ParkingLot) ParkVehicle(vehicle Vehicle, entryTime int64, gateID string) *Ticket {
	minSize := GetMinSpotSize(vehicle.Type)
	for fi := range lot.floors {
		for si := range lot.floors[fi] {
			spot := &lot.floors[fi][si]
			if !spot.IsOccupied && IsCompatible(spot.Size, minSize) {
				spot.IsOccupied = true
				spot.VehicleLicensePlate = vehicle.LicensePlate
				ticketID := "T" + strconv.Itoa(lot.nextTicketID)
				lot.nextTicketID++
				t := &Ticket{
					TicketID:     ticketID,
					LicensePlate: vehicle.LicensePlate,
					SpotID:       spot.SpotID,
					Floor:        spot.Floor,
					EntryTime:    entryTime,
					EntryGateID:  gateID,
				}
				// TODO: store t in lot.activeTickets[ticketID]
				return t
			}
		}
	}
	return nil
}

func (lot *ParkingLot) UnparkVehicle(ticketID string, exitTime int64, gateID string) float64 {
	ticket, ok := lot.activeTickets[ticketID]
	if !ok {
		return -1.0
	}
	// TODO: set ticket.ExitGateID = gateID
	// TODO: find and free the spot
	duration := exitTime - ticket.EntryTime
	var fee float64
	if lot.strategy != nil {
		// TODO: fee = lot.strategy.CalculateFee(duration)
		_ = duration
	} else {
		fee = float64(duration)
	}
	// TODO: delete lot.activeTickets[ticketID]
	return fee
}

func (lot *ParkingLot) GetAvailableSpots(size SpotSize) int {
	count := 0
	for _, floor := range lot.floors {
		for _, spot := range floor {
			if !spot.IsOccupied && spot.Size == size {
				count++
			}
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
		if !spot.IsOccupied && spot.Size == size {
			count++
		}
	}
	return count
}
