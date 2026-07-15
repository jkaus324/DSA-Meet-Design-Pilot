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

// TODO: design and implement your solution.
// Required free functions:
//   func calculateSurge(ctx PricingContext) float64
//   func calculateFare(req RideRequest, ctx PricingContext) float64

func calculateSurge(ctx PricingContext) float64 {
	// TODO: write your solution
	return 0.0
}

func calculateFare(req RideRequest, ctx PricingContext) float64 {
	// TODO: write your solution
	return 0.0
}
