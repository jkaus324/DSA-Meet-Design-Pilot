package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type RideRequest struct {
	UserID   string
	Pickup   string
	Dropoff  string
	RideType string // "economy", "premium", "pool"
}

type Driver struct {
	ID        string
	Rating    float64
	RideType  string
	Available bool
}

type PricingContext struct {
	BaseFare             float64
	AvailableDrivers     int
	ActiveRideRequests   int
	TimeOfDay            string // "morning", "evening", "night"
	Weather              string // "clear", "rain", "storm"
}

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Design and implement a Surge Pricing Engine that:
//   1. Calculates surge multiplier based on supply/demand and conditions
//   2. Lets multiple surge factors (time, weather, demand) combine independently
//   3. Notifies relevant parties when surge changes significantly
//
// Think about:
//   - How do you combine multiple surge factors without one giant if-else chain?
//   - How would you add a "special event" surge factor with zero changes to existing code?
//   - Who needs to be notified when surge changes? How do they subscribe?
//
// Entry points:
//   func CalculateSurge(ctx PricingContext) float64
//   func CalculateFare(req RideRequest, ctx PricingContext) float64
//
// ─────────────────────────────────────────────────────────────────────────────

func CalculateSurge(ctx PricingContext) float64 {
	return 1.0
}

// CalculateFare accepts a RideRequest (for ride-type context) and pricing context.
func CalculateFare(req RideRequest, ctx PricingContext) float64 {
	return ctx.BaseFare
}
