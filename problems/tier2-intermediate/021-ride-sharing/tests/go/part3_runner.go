package main

import "fmt"

func part3Tests() int {
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

	// Test 1: endRide marks ride as inactive
	test("test_end_ride_marks_inactive", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		rideId := service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
		if !service.GetRide(rideId).Active {
			panic("expected ride to be active initially")
		}
		service.EndRide(rideId)
		if service.GetRide(rideId).Active {
			panic("expected ride to be inactive after EndRide")
		}
	})

	// Test 2: After endRide, vehicle can be used for new ride
	test("test_vehicle_freed_after_end", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		ride1 := service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
		if ride1 == "" {
			panic("expected first ride to succeed")
		}

		// Cannot offer again while active
		ride2 := service.OfferRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234")
		if ride2 != "" {
			panic("expected second ride with active vehicle to fail")
		}

		// End first ride
		service.EndRide(ride1)

		// Now can offer again
		ride3 := service.OfferRide("Rohan", "Bangalore", "Chennai", 2, "KA-01-1234")
		if ride3 == "" {
			panic("expected third ride to succeed after ending first")
		}
	})

	// Test 3: Ending already-ended ride is a no-op
	test("test_end_ride_idempotent", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		rideId := service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
		service.EndRide(rideId)
		service.EndRide(rideId) // should not panic
		if service.GetRide(rideId).Active {
			panic("expected ride to remain inactive")
		}
	})

	// Test 4: Ending nonexistent ride is a no-op
	test("test_end_nonexistent_ride", func() {
		service := NewRideService()
		service.EndRide("RIDE-999") // should not panic
	})

	// Test 5: GetRideStats returns correct counts
	test("test_ride_stats", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddUser("Deepa")
		service.AddUser("Priya")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		service.AddVehicle("Deepa", "XUV", "KA-02-5678")

		service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
		service.OfferRide("Deepa", "Bangalore", "Mysore", 5, "KA-02-5678")

		strategy := &MostVacantStrategy{}
		service.SelectRide("Priya", "Bangalore", "Mysore", 1, strategy, "")

		stats := service.GetRideStats()
		if len(stats) == 0 {
			panic("expected non-empty stats")
		}

		foundRohan, foundDeepa, foundPriya := false, false, false
		for _, s := range stats {
			if s.Name == "Rohan" {
				if s.RidesOffered != 1 {
					panic("expected Rohan to have offered 1 ride")
				}
				if s.RidesTaken != 0 {
					panic("expected Rohan to have taken 0 rides")
				}
				foundRohan = true
			}
			if s.Name == "Deepa" {
				if s.RidesOffered != 1 {
					panic("expected Deepa to have offered 1 ride")
				}
				if s.RidesTaken != 0 {
					panic("expected Deepa to have taken 0 rides")
				}
				foundDeepa = true
			}
			if s.Name == "Priya" {
				if s.RidesOffered != 0 {
					panic("expected Priya to have offered 0 rides")
				}
				if s.RidesTaken != 1 {
					panic("expected Priya to have taken 1 ride")
				}
				foundPriya = true
			}
		}
		if !foundRohan || !foundDeepa || !foundPriya {
			panic("expected stats for Rohan, Deepa, and Priya")
		}
	})

	// Test 6: Ended ride not selectable by future passengers
	test("test_ended_ride_not_selectable", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddUser("Priya")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")
		rideId := service.OfferRide("Rohan", "Bangalore", "Mysore", 3, "KA-01-1234")
		service.EndRide(rideId)

		strategy := &MostVacantStrategy{}
		selected := service.SelectRide("Priya", "Bangalore", "Mysore", 1, strategy, "")
		if selected != "" {
			panic("expected ended ride to not be selectable")
		}
	})

	// Test 7: Full workflow — offer, select, end, re-offer
	test("test_full_workflow", func() {
		service := NewRideService()
		service.AddUser("Rohan")
		service.AddUser("Deepa")
		service.AddVehicle("Rohan", "Swift", "KA-01-1234")

		// Rohan offers ride
		ride1 := service.OfferRide("Rohan", "Bangalore", "Mysore", 2, "KA-01-1234")
		if ride1 == "" {
			panic("expected first ride to succeed")
		}

		// Deepa takes ride
		strategy := &MostVacantStrategy{}
		selected := service.SelectRide("Deepa", "Bangalore", "Mysore", 1, strategy, "")
		if selected != ride1 {
			panic("expected Deepa to select ride1")
		}
		if service.GetRide(ride1).AvailableSeats != 1 {
			panic("expected 1 available seat after selection")
		}

		// Rohan ends ride
		service.EndRide(ride1)
		if service.GetRide(ride1).Active {
			panic("expected ride to be inactive after end")
		}

		// Rohan offers new ride with same vehicle
		ride2 := service.OfferRide("Rohan", "Mysore", "Bangalore", 2, "KA-01-1234")
		if ride2 == "" {
			panic("expected second ride to succeed")
		}
		if ride2 == ride1 {
			panic("expected different ride ID")
		}

		// Check stats
		if service.GetUser("Rohan").RidesOffered != 2 {
			panic("expected Rohan to have offered 2 rides")
		}
		if service.GetUser("Deepa").RidesTaken != 1 {
			panic("expected Deepa to have taken 1 ride")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
