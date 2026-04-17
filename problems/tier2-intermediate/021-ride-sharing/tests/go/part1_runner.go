package main

import "fmt"

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

	// Test 1: Add user and verify exists
	test("test_add_user", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		if !service.HasUser("Rohan") {
			panic("expected Rohan to exist")
		}
		if service.HasUser("Unknown") {
			panic("expected Unknown to not exist")
		}
	})

	// Test 2: Add vehicle and verify exists
	test("test_add_vehicle", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		if !service.HasVehicle("KA-01-1234") {
			panic("expected vehicle to exist")
		}
	})

	// Test 3: Offer ride returns valid rideId
	test("test_offer_ride", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		rideId := service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
		if rideId == "" {
			panic("expected non-empty rideId")
		}
		if !service.HasRide(rideId) {
			panic("expected ride to exist")
		}
		ride := service.GetRide(rideId)
		if ride.Origin != "Bangalore" {
			panic("expected origin Bangalore")
		}
		if ride.Destination != "Mysore" {
			panic("expected destination Mysore")
		}
		if ride.TotalSeats != 3 {
			panic("expected 3 total seats")
		}
		if ride.AvailableSeats != 3 {
			panic("expected 3 available seats")
		}
		if !ride.Active {
			panic("expected ride to be active")
		}
	})

	// Test 4: Cannot offer ride with vehicle already in active ride
	test("test_no_duplicate_active_ride_per_vehicle", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		ride1 := service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
		if ride1 == "" {
			panic("expected first ride to succeed")
		}
		ride2 := service.OfferRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234")
		if ride2 != "" {
			panic("expected second ride with same vehicle to fail")
		}
	})

	// Test 5: Cannot offer ride with someone else's vehicle
	test("test_cannot_use_others_vehicle", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddUser("Deepa")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		rideId := service.OfferRide("Deepa", "Bangalore", "Mysore", 2, "KA-01-1234")
		if rideId != "" {
			panic("expected Deepa to be unable to use Rohan's vehicle")
		}
	})

	// Test 6: Offering ride increments ridesOffered
	test("test_rides_offered_counter", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		service.AddVehicle("Rohan", "XUV", "KA-01-5678")
		if service.GetUser("Rohan").RidesOffered != 0 {
			panic("expected 0 rides offered initially")
		}
		service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
		if service.GetUser("Rohan").RidesOffered != 1 {
			panic("expected 1 ride offered after first offer")
		}
		service.OfferRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-5678")
		if service.GetUser("Rohan").RidesOffered != 2 {
			panic("expected 2 rides offered after second offer")
		}
	})

	// Test 7: Offer ride with nonexistent user returns empty
	test("test_offer_ride_invalid_user", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		rideId := service.OfferRide("Ghost", "A", "B", 2, "KA-01-1234")
		if rideId != "" {
			panic("expected empty rideId for nonexistent user")
		}
	})

	// Test 8: Offer ride with nonexistent vehicle returns empty
	test("test_offer_ride_invalid_vehicle", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		rideId := service.OfferRide("Rohan", "A", "B", 2, "INVALID-REG")
		if rideId != "" {
			panic("expected empty rideId for nonexistent vehicle")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
