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

	// Test 1: initial state is Idle
	test("test_initial_state_idle", func() {
		Reset()
		if GetState() != "Idle" {
			panic("expected Idle state")
		}
	})

	// Test 2: select item transitions to PaymentPending (or ItemSelected)
	test("test_select_transitions_state", func() {
		Reset()
		SelectItem("Cola")
		s := GetState()
		if s != "PaymentPending" && s != "ItemSelected" {
			panic("expected PaymentPending or ItemSelected state, got: " + s)
		}
	})

	// Test 3: insert enough money and dispense returns to Idle
	test("test_full_purchase_cycle", func() {
		Reset()
		SelectItem("Cola")
		InsertMoney(25.0) // Cola costs 25
		Dispense()
		if GetState() != "Idle" {
			panic("expected Idle state after dispense")
		}
	})

	// Test 4: cancel from PaymentPending returns to Idle
	test("test_cancel_returns_idle", func() {
		Reset()
		SelectItem("Cola")
		InsertMoney(10.0)
		Cancel()
		if GetState() != "Idle" {
			panic("expected Idle state after cancel")
		}
	})

	// Test 5: pay before select does nothing harmful
	test("test_pay_before_select_safe", func() {
		Reset()
		InsertMoney(50.0) // should be ignored or print warning
		if GetState() != "Idle" {
			panic("expected machine to stay in Idle")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
