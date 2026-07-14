package main

// Data class (given — do not modify).
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

func calculateSurge(ctx PricingContext) float64 {
	// TODO: implement this
	return 0.0
}

func calculateFare(req RideRequest, ctx PricingContext) float64 {
	// TODO: implement this
	return 0.0
}
