package main

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

// ─── NEW in Extension 1 ───────────────────────────────────────────────────────
//
// The ops team wants DYNAMIC SURGE NOTIFICATIONS:
//   - When the surge multiplier changes by more than 0.5x, notify stakeholders
//   - Stakeholders: drivers (to encourage them to go online), ops dashboard
//   - Drivers only receive notifications relevant to their ride type
//
// Think about:
//   - How do you combine the Strategy pattern (surge calculation) with
//     the Observer pattern (surge notifications)?
//   - Does the surge engine become the subject? Or is there a separate notifier?
//   - How do you filter driver notifications by ride type?
//
// Entry points:
//   func CalculateSurge(ctx PricingContext) float64
//   func CalculateFare(req RideRequest, ctx PricingContext) float64
//   func RegisterSurgeObserver(observer SurgeObserver)
//
// ─────────────────────────────────────────────────────────────────────────────

type SurgeObserver interface {
	OnSurgeChange(oldMult, newMult float64, rideType string)
}

func CalculateSurge(ctx PricingContext) float64 {
	return 1.0
}

func CalculateFare(req RideRequest, ctx PricingContext) float64 {
	return ctx.BaseFare
}

func RegisterSurgeObserver(observer SurgeObserver) {}
