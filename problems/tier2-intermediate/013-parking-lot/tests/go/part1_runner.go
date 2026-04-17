package main

import (
	"fmt"
	"strings"
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

	// Test 1: Park a car in a medium spot
	test("test_park_car", func() {
		lot := NewParkingLot(2)
		lot.AddSpot(0, MEDIUM)
		car := Vehicle{LicensePlate: "ABC123", Type: CAR}
		t := lot.ParkVehicle(car, 1000)
		if t == nil {
			panic("expected ticket, got nil")
		}
		if t.LicensePlate != "ABC123" {
			panic("expected license ABC123")
		}
		if t.Floor != 0 {
			panic("expected floor 0")
		}
	})

	// Test 2: Park a motorcycle in a small spot
	test("test_park_motorcycle", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, SMALL)
		moto := Vehicle{LicensePlate: "MOTO1", Type: MOTORCYCLE}
		t := lot.ParkVehicle(moto, 1000)
		if t == nil {
			panic("expected ticket")
		}
		if t.LicensePlate != "MOTO1" {
			panic("expected MOTO1")
		}
	})

	// Test 3: Car cannot park in small spot
	test("test_car_no_small_spot", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, SMALL)
		car := Vehicle{LicensePlate: "CAR1", Type: CAR}
		t := lot.ParkVehicle(car, 1000)
		if t != nil {
			panic("expected nil — car should not fit in small spot")
		}
	})

	// Test 4: Motorcycle can park in a larger spot when small is unavailable
	test("test_motorcycle_in_medium_spot", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		moto := Vehicle{LicensePlate: "MOTO2", Type: MOTORCYCLE}
		t := lot.ParkVehicle(moto, 1000)
		if t == nil {
			panic("expected ticket — motorcycle should fit in medium spot")
		}
	})

	// Test 5: Truck requires large spot
	test("test_truck_large_spot", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		lot.AddSpot(0, LARGE)
		truck := Vehicle{LicensePlate: "TRUCK1", Type: TRUCK}
		t := lot.ParkVehicle(truck, 1000)
		if t == nil {
			panic("expected ticket")
		}
		if !strings.Contains(t.SpotID, "S1") {
			panic("expected second spot (large), spotID=" + t.SpotID)
		}
	})

	// Test 6: Unpark returns a fee and frees the spot
	test("test_unpark_frees_spot", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		car := Vehicle{LicensePlate: "CAR2", Type: CAR}
		t := lot.ParkVehicle(car, 1000)
		if t == nil {
			panic("expected ticket")
		}
		tid := t.TicketID
		if lot.GetAvailableSpots(MEDIUM) != 0 {
			panic("expected 0 available spots after parking")
		}
		fee := lot.UnparkVehicle(tid, 4600) // 3600 seconds
		if fee < 0 {
			panic("expected non-negative fee")
		}
		if lot.GetAvailableSpots(MEDIUM) != 1 {
			panic("expected 1 available spot after unparking")
		}
	})

	// Test 7: Invalid ticket returns -1
	test("test_invalid_ticket", func() {
		lot := NewParkingLot(1)
		fee := lot.UnparkVehicle("INVALID", 5000)
		if fee >= 0 {
			panic("expected negative fee for invalid ticket")
		}
	})

	// Test 8: Available spots count is correct
	test("test_available_spots_count", func() {
		lot := NewParkingLot(2)
		lot.AddSpot(0, SMALL)
		lot.AddSpot(0, MEDIUM)
		lot.AddSpot(0, MEDIUM)
		lot.AddSpot(1, LARGE)
		if lot.GetAvailableSpots(SMALL) != 1 {
			panic("expected 1 small spot")
		}
		if lot.GetAvailableSpots(MEDIUM) != 2 {
			panic("expected 2 medium spots")
		}
		if lot.GetAvailableSpots(LARGE) != 1 {
			panic("expected 1 large spot")
		}
	})

	// Test 9: Available spots by floor
	test("test_available_spots_by_floor", func() {
		lot := NewParkingLot(2)
		lot.AddSpot(0, MEDIUM)
		lot.AddSpot(0, MEDIUM)
		lot.AddSpot(1, MEDIUM)
		if lot.GetAvailableSpotsByFloor(0, MEDIUM) != 2 {
			panic("expected 2 medium spots on floor 0")
		}
		if lot.GetAvailableSpotsByFloor(1, MEDIUM) != 1 {
			panic("expected 1 medium spot on floor 1")
		}
	})

	// Test 10: Nearest spot allocation — prefers lower floor
	test("test_nearest_floor_allocation", func() {
		lot := NewParkingLot(3)
		lot.AddSpot(0, SMALL)   // floor 0 has only small
		lot.AddSpot(1, MEDIUM)  // floor 1 has medium
		lot.AddSpot(2, MEDIUM)  // floor 2 has medium
		car := Vehicle{LicensePlate: "CAR3", Type: CAR}
		t := lot.ParkVehicle(car, 1000)
		if t == nil {
			panic("expected ticket")
		}
		if t.Floor != 1 {
			panic(fmt.Sprintf("expected floor 1 (nearest), got %d", t.Floor))
		}
	})

	// Test 11: Full lot returns nil
	test("test_full_lot", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		car1 := Vehicle{LicensePlate: "C1", Type: CAR}
		car2 := Vehicle{LicensePlate: "C2", Type: CAR}
		if lot.ParkVehicle(car1, 1000) == nil {
			panic("expected first car to park")
		}
		if lot.ParkVehicle(car2, 1000) != nil {
			panic("expected nil — lot is full")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
