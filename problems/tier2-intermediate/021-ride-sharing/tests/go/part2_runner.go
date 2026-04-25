package main

import "fmt"

func setupTestService() *RideService {
	service := NewRideService()
	service.AddUser("Rohan")
	service.AddUser("Deepa")
	service.AddUser("Amit")
	service.AddUser("Priya")

	service.AddVehicle("Rohan", "Swift", "KA-01-1234")
	service.AddVehicle("Deepa", "XUV", "KA-02-5678")
	service.AddVehicle("Amit", "Swift", "KA-03-9012")

	// Rohan offers Bangalore->Mysore, 3 seats, Swift
	service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
	// Deepa offers Bangalore->Mysore, 5 seats, XUV
	service.OfferRide("Deepa", "Bangalore", "Mysore", 5, "KA-02-5678")
	// Amit offers Bangalore->Chennai, 2 seats, Swift
	service.OfferRide("Amit", "Bangalore", "Chennai", 2, "KA-03-9012")

	return service
}

func part2Tests() int {
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

	// Test 1: MostVacant selects ride with most available seats
	test("test_most_vacant_strategy", func() {
		service := setupTestService()
		strategy := &MostVacantStrategy{}
		rideId := service.SelectRide("Priya", "Bangalore", "Mysore", 1, strategy, "")
		if rideId == "" {
			panic("expected a ride to be selected")
		}
		ride := service.GetRide(rideId)
		// Deepa's ride has 5 seats (most vacant)
		if ride.DriverId != "Deepa" {
			panic("expected MostVacant to select Deepa's ride (5 seats)")
		}
	})

	// Test 2: PreferredVehicle selects ride with matching model
	test("test_preferred_vehicle_strategy", func() {
		service := setupTestService()
		strategy := NewPreferredVehicleStrategy(service.GetVehicles())
		rideId := service.SelectRide("Priya", "Bangalore", "Mysore", 1, strategy, "Swift")
		if rideId == "" {
			panic("expected a ride to be selected")
		}
		ride := service.GetRide(rideId)
		// Rohan has a Swift on this route
		if ride.DriverId != "Rohan" {
			panic("expected PreferredVehicle to select Rohan's Swift ride")
		}
	})

	// Test 3: No match for route returns empty
	test("test_no_matching_route", func() {
		service := setupTestService()
		strategy := &MostVacantStrategy{}
		rideId := service.SelectRide("Priya", "Delhi", "Mumbai", 1, strategy, "")
		if rideId != "" {
			panic("expected empty rideId for unmatched route")
		}
	})

	// Test 4: Selecting ride decrements available seats
	test("test_seats_decremented", func() {
		service := setupTestService()
		strategy := &MostVacantStrategy{}
		rideId := service.SelectRide("Priya", "Bangalore", "Mysore", 2, strategy, "")
		if rideId == "" {
			panic("expected a ride to be selected")
		}
		ride := service.GetRide(rideId)
		// Deepa's ride: 5 total, now 3 available
		if ride.AvailableSeats != 3 {
			panic("expected 3 available seats after selecting 2")
		}
	})

	// Test 5: Selecting ride increments passenger's ridesTaken
	test("test_rides_taken_incremented", func() {
		service := setupTestService()
		strategy := &MostVacantStrategy{}
		if service.GetUser("Priya").RidesTaken != 0 {
			panic("expected 0 rides taken initially")
		}
		service.SelectRide("Priya", "Bangalore", "Mysore", 1, strategy, "")
		if service.GetUser("Priya").RidesTaken != 1 {
			panic("expected 1 ride taken after selection")
		}
	})

	// Test 6: Cannot select own offered ride
	test("test_cannot_select_own_ride", func() {
		service := setupTestService()
		strategy := &MostVacantStrategy{}
		// Deepa tries to select a Bangalore->Mysore ride, but she offered one
		// Only Rohan's ride should be a candidate for Deepa
		rideId := service.SelectRide("Deepa", "Bangalore", "Mysore", 1, strategy, "")
		if rideId != "" {
			ride := service.GetRide(rideId)
			if ride.DriverId == "Deepa" {
				panic("must not select own ride")
			}
		}
	})

	// Test 7: Not enough seats returns empty
	test("test_not_enough_seats", func() {
		service := setupTestService()
		strategy := &MostVacantStrategy{}
		// Request 10 seats — no ride has that many
		rideId := service.SelectRide("Priya", "Bangalore", "Mysore", 10, strategy, "")
		if rideId != "" {
			panic("expected empty rideId when not enough seats")
		}
	})

	// Test 8: PreferredVehicle with no matching model returns empty
	test("test_no_matching_vehicle_model", func() {
		service := setupTestService()
		strategy := NewPreferredVehicleStrategy(service.GetVehicles())
		rideId := service.SelectRide("Priya", "Bangalore", "Mysore", 1, strategy, "BMW")
		if rideId != "" {
			panic("expected empty rideId for non-matching model")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
