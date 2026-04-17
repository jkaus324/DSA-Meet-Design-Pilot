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

type PricingContext struct {
	BaseFare           float64
	AvailableDrivers   int
	ActiveRideRequests int
	TimeOfDay          string // "morning", "evening", "night"
	Weather            string // "clear", "rain", "storm"
}

// ─── Surge Strategy Interface ─────────────────────────────────────────────────

type SurgeStrategy interface {
	Multiplier(ctx PricingContext) float64
}

// ─── Concrete Surge Strategies ────────────────────────────────────────────────

type DemandSurge struct{}

func (d *DemandSurge) Multiplier(ctx PricingContext) float64 {
	// TODO: Calculate demand ratio (requests / drivers)
	//       Return 1.0 if balanced, up to 2.0x if demand >> supply
	//       If drivers == 0, return 2.5
	//       ratio > 3.0 → 2.0, > 2.0 → 1.5, > 1.5 → 1.25, else 1.0
	return 1.0
}

type WeatherSurge struct{}

func (w *WeatherSurge) Multiplier(ctx PricingContext) float64 {
	// TODO: Return multiplier based on weather
	//       "clear" → 1.0, "rain" → 1.3, "storm" → 1.8
	return 1.0
}

type TimeSurge struct{}

func (t *TimeSurge) Multiplier(ctx PricingContext) float64 {
	// TODO: Return multiplier based on time of day
	//       "morning" → 1.2, "evening" → 1.5, "night" → 1.0
	return 1.0
}

// ─── Surge Observer Interface ─────────────────────────────────────────────────

type SurgeObserver interface {
	OnSurgeChange(oldMultiplier, newMultiplier float64)
}

type DriverNotifier struct{}

func (d *DriverNotifier) OnSurgeChange(old, new_ float64) {
	// TODO: Print driver-facing alert if new surge > 1.5x
	_ = fmt.Sprintf("driver alert: %.2f -> %.2f", old, new_)
}

type RiderNotifier struct{}

func (r *RiderNotifier) OnSurgeChange(old, new_ float64) {
	// TODO: Print rider-facing warning if new surge > 1.5x
	_ = fmt.Sprintf("rider warning: %.2f -> %.2f", old, new_)
}

// ─── Pricing Engine ───────────────────────────────────────────────────────────

type PricingEngine struct {
	strategies []SurgeStrategy
	observers  []SurgeObserver
	lastSurge  float64
}

func NewPricingEngine() *PricingEngine {
	return &PricingEngine{lastSurge: 1.0}
}

func (e *PricingEngine) AddStrategy(s SurgeStrategy) { e.strategies = append(e.strategies, s) }
func (e *PricingEngine) AddObserver(o SurgeObserver)  { e.observers = append(e.observers, o) }

func (e *PricingEngine) CalculateSurge(ctx PricingContext) float64 {
	mult := 1.0
	for _, s := range e.strategies {
		// TODO: multiply all strategy multipliers together
		_ = s
	}
	// TODO: cap at 3.0
	mult = math.Min(mult, 3.0)
	// TODO: if changed by > 0.5 from lastSurge, notify observers
	// TODO: update e.lastSurge and return mult
	return mult
}

func (e *PricingEngine) CalculateFare(ctx PricingContext) float64 {
	// TODO: return ctx.BaseFare * e.CalculateSurge(ctx)
	return ctx.BaseFare
}

// ─── Test Entry Points ────────────────────────────────────────────────────────

func CalculateSurge(ctx PricingContext) float64 {
	engine := NewPricingEngine()
	engine.AddStrategy(&DemandSurge{})
	engine.AddStrategy(&WeatherSurge{})
	engine.AddStrategy(&TimeSurge{})
	return engine.CalculateSurge(ctx)
}

func CalculateFare(req RideRequest, ctx PricingContext) float64 {
	engine := NewPricingEngine()
	engine.AddStrategy(&DemandSurge{})
	engine.AddStrategy(&WeatherSurge{})
	engine.AddStrategy(&TimeSurge{})
	return engine.CalculateFare(ctx)
}
