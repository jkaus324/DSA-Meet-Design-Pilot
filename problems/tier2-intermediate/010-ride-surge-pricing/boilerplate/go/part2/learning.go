package main

import (
	"fmt"
	"math"
)

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type RideRequest struct {
	UserID   string
	Pickup   string
	Dropoff  string
	RideType string
}

type Driver struct {
	ID        string
	Rating    float64
	RideType  string
	Available bool
}

type PricingContext struct {
	BaseFare           float64
	AvailableDrivers   int
	ActiveRideRequests int
	TimeOfDay          string
	Weather            string
}

// ─── Surge Strategy Interface ─────────────────────────────────────────────────

type SurgeStrategy interface {
	Multiplier(ctx PricingContext) float64
}

// ─── Concrete Surge Strategies (from Part 1) ──────────────────────────────────

type DemandSurge struct{}

func (d *DemandSurge) Multiplier(ctx PricingContext) float64 {
	if ctx.AvailableDrivers == 0 {
		return 2.5
	}
	ratio := float64(ctx.ActiveRideRequests) / float64(ctx.AvailableDrivers)
	if ratio > 3.0 {
		return 2.0
	}
	if ratio > 2.0 {
		return 1.5
	}
	if ratio > 1.5 {
		return 1.25
	}
	return 1.0
}

type WeatherSurge struct{}

func (w *WeatherSurge) Multiplier(ctx PricingContext) float64 {
	switch ctx.Weather {
	case "storm":
		return 2.0
	case "rain":
		return 1.3
	}
	return 1.0
}

type TimeSurge struct{}

func (t *TimeSurge) Multiplier(ctx PricingContext) float64 {
	switch ctx.TimeOfDay {
	case "evening":
		return 1.4
	case "morning":
		return 1.2
	}
	return 1.0
}

// ─── Surge Observer Interface ─────────────────────────────────────────────────

type SurgeObserver interface {
	OnSurgeChange(oldMult, newMult float64, rideType string)
}

type DriverObserver struct {
	Driver Driver
}

func (d *DriverObserver) OnSurgeChange(oldMult, newMult float64, rideType string) {
	if d.Driver.RideType == rideType || rideType == "all" {
		fmt.Printf("[DRIVER %s] Surge changed: %.2fx -> %.2fx (%s)\n", d.Driver.ID, oldMult, newMult, rideType)
	}
}

type OpsDashboardObserver struct{}

func (o *OpsDashboardObserver) OnSurgeChange(oldMult, newMult float64, rideType string) {
	fmt.Printf("[OPS] Surge alert for %s: %.2fx -> %.2fx\n", rideType, oldMult, newMult)
}

// ─── SurgePricingEngine ───────────────────────────────────────────────────────

type SurgePricingEngine struct {
	strategies      []SurgeStrategy
	observers       []SurgeObserver
	lastMultiplier  float64
	changeThreshold float64
}

func NewSurgePricingEngine() *SurgePricingEngine {
	return &SurgePricingEngine{lastMultiplier: 1.0, changeThreshold: 0.5}
}

func (e *SurgePricingEngine) AddStrategy(s SurgeStrategy) { e.strategies = append(e.strategies, s) }
func (e *SurgePricingEngine) AddObserver(o SurgeObserver)  { e.observers = append(e.observers, o) }

func (e *SurgePricingEngine) CalculateSurge(ctx PricingContext, rideType string) float64 {
	mult := 1.0
	for _, s := range e.strategies {
		m := s.Multiplier(ctx)
		if m > mult {
			mult = m
		}
	}
	// TODO: if math.Abs(mult - e.lastMultiplier) > e.changeThreshold, notify all observers
	// HINT: for _, o := range e.observers { o.OnSurgeChange(e.lastMultiplier, mult, rideType) }
	_ = math.Abs(mult - e.lastMultiplier)
	e.lastMultiplier = mult
	return mult
}

// ─── Global engine and entry points ──────────────────────────────────────────

var globalEngine *SurgePricingEngine

func init() {
	globalEngine = NewSurgePricingEngine()
	globalEngine.AddStrategy(&DemandSurge{})
	globalEngine.AddStrategy(&WeatherSurge{})
	globalEngine.AddStrategy(&TimeSurge{})
}

func CalculateSurge(ctx PricingContext) float64 {
	return globalEngine.CalculateSurge(ctx, "all")
}

func CalculateFare(req RideRequest, ctx PricingContext) float64 {
	return ctx.BaseFare * globalEngine.CalculateSurge(ctx, req.RideType)
}

func RegisterSurgeObserver(obs SurgeObserver) {
	// TODO: globalEngine.AddObserver(obs)
}
