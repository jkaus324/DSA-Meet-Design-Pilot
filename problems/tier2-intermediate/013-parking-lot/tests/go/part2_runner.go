package main

import (
	"fmt"
	"math"
)

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

	approxEqual := func(a, b float64) bool {
		return math.Abs(a-b) < 0.01
	}

	// Test 1: FlatRate pricing — same fee regardless of duration
	test("test_flat_rate", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		flat := NewFlatRate(10.0)
		lot.SetPricingStrategy(flat)
		car := Vehicle{LicensePlate: "CAR1", Type: CAR}
		t := lot.ParkVehicle(car, 1000, "G1")
		if t == nil {
			panic("expected ticket")
		}
		tid := t.TicketID
		fee := lot.UnparkVehicle(tid, 8600, "G2") // 7600s ~= 2.1 hours
		if !approxEqual(fee, 10.0) {
			panic(fmt.Sprintf("expected fee 10.0, got %.2f", fee))
		}
	})

	// Test 2: Hourly pricing — rounds up to full hours
	test("test_hourly_rate", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		hourly := NewHourly(5.0)
		lot.SetPricingStrategy(hourly)
		car := Vehicle{LicensePlate: "CAR2", Type: CAR}
		t := lot.ParkVehicle(car, 0, "G1")
		if t == nil {
			panic("expected ticket")
		}
		tid := t.TicketID
		// 9000s = 2.5h, ceil = 3h => $15
		fee := lot.UnparkVehicle(tid, 9000, "G2")
		if !approxEqual(fee, 15.0) {
			panic(fmt.Sprintf("expected fee 15.0, got %.2f", fee))
		}
	})

	// Test 3: Hourly pricing — exactly 1 hour
	test("test_hourly_exact_hour", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		hourly := NewHourly(5.0)
		lot.SetPricingStrategy(hourly)
		car := Vehicle{LicensePlate: "CAR3", Type: CAR}
		t := lot.ParkVehicle(car, 0, "G1")
		if t == nil {
			panic("expected ticket")
		}
		tid := t.TicketID
		fee := lot.UnparkVehicle(tid, 3600, "G2") // 1 hour => $5
		if !approxEqual(fee, 5.0) {
			panic(fmt.Sprintf("expected fee 5.0, got %.2f", fee))
		}
	})

	// Test 4: Tiered pricing — under 1 hour (base rate only)
	test("test_tiered_base_rate", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		tiered := NewTiered(10.0, 8.0, 5.0)
		lot.SetPricingStrategy(tiered)
		car := Vehicle{LicensePlate: "CAR4", Type: CAR}
		t := lot.ParkVehicle(car, 0, "G1")
		if t == nil {
			panic("expected ticket")
		}
		tid := t.TicketID
		fee := lot.UnparkVehicle(tid, 3000, "G2") // 3000s < 1h, ceil=1 => $10
		if !approxEqual(fee, 10.0) {
			panic(fmt.Sprintf("expected fee 10.0, got %.2f", fee))
		}
	})

	// Test 5: Tiered pricing — 2 hours (base + 1 mid)
	test("test_tiered_mid_rate", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		tiered := NewTiered(10.0, 8.0, 5.0)
		lot.SetPricingStrategy(tiered)
		car := Vehicle{LicensePlate: "CAR5", Type: CAR}
		t := lot.ParkVehicle(car, 0, "G1")
		if t == nil {
			panic("expected ticket")
		}
		tid := t.TicketID
		fee := lot.UnparkVehicle(tid, 7200, "G2") // 7200s = 2h => $10 + $8 = $18
		if !approxEqual(fee, 18.0) {
			panic(fmt.Sprintf("expected fee 18.0, got %.2f", fee))
		}
	})

	// Test 6: Tiered pricing — 5 hours (base + 2*mid + 2*high)
	test("test_tiered_high_rate", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		tiered := NewTiered(10.0, 8.0, 5.0)
		lot.SetPricingStrategy(tiered)
		car := Vehicle{LicensePlate: "CAR6", Type: CAR}
		t := lot.ParkVehicle(car, 0, "G1")
		if t == nil {
			panic("expected ticket")
		}
		tid := t.TicketID
		fee := lot.UnparkVehicle(tid, 18000, "G2") // 18000s = 5h => $10 + $16 + $10 = $36
		if !approxEqual(fee, 36.0) {
			panic(fmt.Sprintf("expected fee 36.0, got %.2f", fee))
		}
	})

	// Test 7: Swap strategy at runtime
	test("test_swap_strategy", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		lot.AddSpot(0, MEDIUM)
		flat := NewFlatRate(10.0)
		hourly := NewHourly(5.0)
		lot.SetPricingStrategy(flat)
		car1 := Vehicle{LicensePlate: "SW1", Type: CAR}
		t1 := lot.ParkVehicle(car1, 0, "G1")
		if t1 == nil {
			panic("expected ticket for SW1")
		}
		fee1 := lot.UnparkVehicle(t1.TicketID, 7200, "G2") // flat = $10
		if !approxEqual(fee1, 10.0) {
			panic(fmt.Sprintf("expected fee1 10.0, got %.2f", fee1))
		}
		lot.SetPricingStrategy(hourly)
		car2 := Vehicle{LicensePlate: "SW2", Type: CAR}
		t2 := lot.ParkVehicle(car2, 0, "G1")
		if t2 == nil {
			panic("expected ticket for SW2")
		}
		fee2 := lot.UnparkVehicle(t2.TicketID, 7200, "G2") // hourly: 2h * $5 = $10
		if !approxEqual(fee2, 10.0) {
			panic(fmt.Sprintf("expected fee2 10.0, got %.2f", fee2))
		}
	})

	// Test 8: Add and retrieve gates
	test("test_gate_management", func() {
		lot := NewParkingLot(1)
		lot.AddGate("E1", ENTRY)
		lot.AddGate("E2", ENTRY)
		lot.AddGate("X1", EXIT)
		entryGates := lot.GetGates(ENTRY)
		exitGates := lot.GetGates(EXIT)
		if len(entryGates) != 2 {
			panic(fmt.Sprintf("expected 2 entry gates, got %d", len(entryGates)))
		}
		if len(exitGates) != 1 {
			panic(fmt.Sprintf("expected 1 exit gate, got %d", len(exitGates)))
		}
		if exitGates[0] != "X1" {
			panic("expected X1 as exit gate")
		}
	})

	// Test 9: Gate IDs recorded on ticket
	test("test_gate_on_ticket", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, MEDIUM)
		flat := NewFlatRate(10.0)
		lot.SetPricingStrategy(flat)
		lot.AddGate("ENTRY1", ENTRY)
		lot.AddGate("EXIT1", EXIT)
		car := Vehicle{LicensePlate: "GATE_CAR", Type: CAR}
		t := lot.ParkVehicle(car, 0, "ENTRY1")
		if t == nil {
			panic("expected ticket")
		}
		if t.EntryGateID != "ENTRY1" {
			panic("expected entryGateID == ENTRY1")
		}
	})

	// Test 10: Hourly pricing — very short stay rounds up to 1 hour
	test("test_short_stay_rounds_up", func() {
		lot := NewParkingLot(1)
		lot.AddSpot(0, SMALL)
		hourly := NewHourly(5.0)
		lot.SetPricingStrategy(hourly)
		moto := Vehicle{LicensePlate: "SHORT", Type: MOTORCYCLE}
		t := lot.ParkVehicle(moto, 0, "G1")
		if t == nil {
			panic("expected ticket")
		}
		fee := lot.UnparkVehicle(t.TicketID, 1, "G2") // 1 second => ceil = 1 hour => $5
		if !approxEqual(fee, 5.0) {
			panic(fmt.Sprintf("expected fee 5.0, got %.2f", fee))
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
