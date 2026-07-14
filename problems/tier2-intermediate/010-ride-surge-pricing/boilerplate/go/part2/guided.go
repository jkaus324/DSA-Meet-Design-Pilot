package main

// Data class (given).
type PricingContext struct {
	baseFare float64
	availableDrivers int
	activeRideRequests int
	timeOfDay string
	weather string
}

type RideRequest struct {
	userId string
	pickup string
	dropoff string
	rideType string
}

// HINT: introduce an abstraction so new rules don't change existing code.

// HINT: pick the field that defines 'better' for this ranking and compare the two.
func calculateSurge(ctx PricingContext) float64 {
	// TODO: write your solution
	return 0.0
}

// HINT: pick the field that defines 'better' for this ranking and compare the two.
func calculateFare(req RideRequest, ctx PricingContext) float64 {
	// TODO: write your solution
	return 0.0
}
