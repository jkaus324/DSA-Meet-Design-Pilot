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

// ─── SurgeStrategy Interface ──────────────────────────────────────────────────

type SurgeStrategy interface {
	Multiplier(ctx PricingContext) float64
}

// Copy your Part 1 strategies here (DemandSurge, WeatherSurge, TimeSurge)

// ─── NEW: SurgeObserver Interface ─────────────────────────────────────────────
// HINT: Observers receive the old multiplier, new multiplier, and ride type.

type SurgeObserver interface {
	OnSurgeChange(oldMult, newMult float64, rideType string)
}

// TODO: Implement DriverObserver (filters by driver's RideType)
// TODO: Implement OpsDashboardObserver (receives all surge changes)

// ─── SurgePricingEngine ───────────────────────────────────────────────────────

type SurgePricingEngine struct {
	strategies      []SurgeStrategy
	observers       []SurgeObserver
	lastMultiplier  float64
	changeThreshold float64
}

func NewSurgePricingEngine() *SurgePricingEngine {
	// HINT: initialise lastMultiplier=1.0, changeThreshold=0.5
	return &SurgePricingEngine{}
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
	// TODO: if |mult - e.lastMultiplier| > e.changeThreshold, notify observers
	e.lastMultiplier = mult
	return mult
}

// ─── Global engine and entry points ──────────────────────────────────────────

var globalEngine = NewSurgePricingEngine()

func CalculateSurge(ctx PricingContext) float64 {
	// TODO: implement using globalEngine
	return 1.0
}

func CalculateFare(req RideRequest, ctx PricingContext) float64 {
	return ctx.BaseFare * globalEngine.CalculateSurge(ctx, req.RideType)
}

func RegisterSurgeObserver(obs SurgeObserver) {
	// TODO: register with globalEngine
}
