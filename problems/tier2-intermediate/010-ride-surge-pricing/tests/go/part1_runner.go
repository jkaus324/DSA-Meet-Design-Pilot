package main

import (
	"fmt"
	"math"
)

func part1Tests() int {
	passed := 0
	failed := 0

	test := func(name string, fn func()) {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("FAIL", name)
					failed++
				}
			}()
			fn()
			fmt.Println("PASS", name)
			passed++
		}()
	}

	approxEqual := func(a, b, eps float64) bool {
		return math.Abs(a-b) < eps
	}

	// Test 1: no surge when supply > demand
	test("test_no_surge_normal_conditions", func() {
		ctx := PricingContext{BaseFare: 10.0, AvailableDrivers: 20, ActiveRideRequests: 5, TimeOfDay: "morning", Weather: "clear"}
		surge := CalculateSurge(ctx)
		if surge < 1.0 {
			panic("surge must be >= 1.0")
		}
		if surge > 1.5 {
			panic("surge too high for low demand")
		}
	})

	// Test 2: high demand ratio causes surge
	test("test_high_demand_causes_surge", func() {
		ctx := PricingContext{BaseFare: 10.0, AvailableDrivers: 2, ActiveRideRequests: 10, TimeOfDay: "evening", Weather: "clear"}
		surge := CalculateSurge(ctx)
		if surge <= 1.0 {
			panic("expected surge > 1.0 for high demand")
		}
	})

	// Test 3: storm weather increases surge
	test("test_storm_increases_surge", func() {
		ctx := PricingContext{BaseFare: 10.0, AvailableDrivers: 10, ActiveRideRequests: 10, TimeOfDay: "morning", Weather: "storm"}
		surge := CalculateSurge(ctx)
		if surge <= 1.0 {
			panic("expected surge > 1.0 during storm")
		}
	})

	// Test 4: fare = baseFare * surgeMultiplier
	test("test_fare_calculation", func() {
		ctx := PricingContext{BaseFare: 100.0, AvailableDrivers: 5, ActiveRideRequests: 5, TimeOfDay: "morning", Weather: "clear"}
		req := RideRequest{UserID: "u1", Pickup: "A", Dropoff: "B", RideType: "economy"}
		fare := CalculateFare(req, ctx)
		surge := CalculateSurge(ctx)
		if !approxEqual(fare, 100.0*surge, 0.01) {
			panic(fmt.Sprintf("fare mismatch: got %.2f, expected %.2f", fare, 100.0*surge))
		}
	})

	// Test 5: surge is never below 1.0
	test("test_surge_minimum_one", func() {
		ctx := PricingContext{BaseFare: 50.0, AvailableDrivers: 100, ActiveRideRequests: 1, TimeOfDay: "morning", Weather: "clear"}
		surge := CalculateSurge(ctx)
		if surge < 1.0 {
			panic("surge must always be >= 1.0")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
