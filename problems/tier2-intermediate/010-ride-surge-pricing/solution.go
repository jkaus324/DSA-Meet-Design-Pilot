// Ride surge pricing — Strategy + Observer reference solution (Go).
package main

type PricingContext struct {
	baseFare           float64
	availableDrivers   int
	activeRideRequests int
	timeOfDay          string
	weather            string
}

type RideRequest struct {
	userId   string
	pickup   string
	dropoff  string
	rideType string
}

type SurgeStrategy interface {
	multiplier(ctx PricingContext) float64
}

type DemandSurge struct{}

func (DemandSurge) multiplier(ctx PricingContext) float64 {
	if ctx.availableDrivers == 0 {
		return 2.5
	}
	ratio := float64(ctx.activeRideRequests) / float64(ctx.availableDrivers)
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

func (WeatherSurge) multiplier(ctx PricingContext) float64 {
	if ctx.weather == "storm" {
		return 1.8
	}
	if ctx.weather == "rain" {
		return 1.3
	}
	return 1.0
}

type TimeSurge struct{}

func (TimeSurge) multiplier(ctx PricingContext) float64 {
	if ctx.timeOfDay == "evening" {
		return 1.5
	}
	if ctx.timeOfDay == "morning" {
		return 1.2
	}
	return 1.0
}

type SurgeObserver interface {
	onSurgeChange(oldSurge, newSurge float64, rideType string)
}

type PricingEngine struct {
	strategies []SurgeStrategy
	observers  []SurgeObserver
	lastSurge  float64
}

const changeThreshold = 0.5

func (e *PricingEngine) addStrategy(s SurgeStrategy) { e.strategies = append(e.strategies, s) }
func (e *PricingEngine) addObserver(o SurgeObserver) { e.observers = append(e.observers, o) }
func (e *PricingEngine) clearObservers()             { e.observers = nil }

func (e *PricingEngine) calculateSurge(ctx PricingContext, rideType string) float64 {
	mult := 1.0
	for _, s := range e.strategies {
		if m := s.multiplier(ctx); m > mult {
			mult = m
		}
	}
	if mult > 3.0 {
		mult = 3.0
	}
	if abs(mult-e.lastSurge) > changeThreshold {
		for _, o := range e.observers {
			o.onSurgeChange(e.lastSurge, mult, rideType)
		}
	}
	e.lastSurge = mult
	return mult
}

func (e *PricingEngine) calculateFare(ctx PricingContext, rideType string) float64 {
	return ctx.baseFare * e.calculateSurge(ctx, rideType)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

var engine = func() *PricingEngine {
	e := &PricingEngine{lastSurge: 1.0}
	e.addStrategy(DemandSurge{})
	e.addStrategy(WeatherSurge{})
	e.addStrategy(TimeSurge{})
	return e
}()

func calculateSurge(ctx PricingContext) float64 {
	return engine.calculateSurge(ctx, "all")
}

func calculateFare(req RideRequest, ctx PricingContext) float64 {
	return engine.calculateFare(ctx, req.rideType)
}
