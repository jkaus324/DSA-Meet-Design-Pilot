package main

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
// HINT: Each surge factor (demand, weather, time) is an independent strategy.
// They each contribute a multiplier that combines into the final surge.

type SurgeStrategy interface {
	// HINT: name your method Multiplier(ctx PricingContext) float64
	Multiplier(ctx PricingContext) float64
}

// TODO: Implement concrete surge strategies:
//   - DemandSurge  (based on AvailableDrivers vs ActiveRideRequests ratio)
//   - WeatherSurge (based on Weather condition)
//   - TimeSurge    (based on TimeOfDay)

// ─── Surge Observer Interface ─────────────────────────────────────────────────
// HINT: These are notified when the surge multiplier changes significantly.

type SurgeObserver interface {
	OnSurgeChange(oldMultiplier, newMultiplier float64)
}

// TODO: Implement observers:
//   - DriverNotifier (tells drivers surge is high)
//   - RiderNotifier  (warns riders about high surge)

// ─── Pricing Engine ───────────────────────────────────────────────────────────
// TODO: Implement PricingEngine that:
//   - Holds a list of surge strategies
//   - Combines their multipliers (multiply them together, cap at 3.0x)
//   - Notifies observers when surge changes by > 0.5x
//   - Has CalculateSurge() and CalculateFare() methods

// ─── Test Entry Points ────────────────────────────────────────────────────────

func CalculateSurge(ctx PricingContext) float64 {
	return 1.0
}

func CalculateFare(req RideRequest, ctx PricingContext) float64 {
	return ctx.BaseFare
}
