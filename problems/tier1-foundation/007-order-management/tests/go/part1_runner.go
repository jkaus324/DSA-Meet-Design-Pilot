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

	// Test 1: create order returns valid ID and state is Created
	test("test_create_order", func() {
		ResetManager()
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 2}}, 500.0)
		if id == "" {
			panic("expected non-empty ID")
		}
		if GetOrderState(id) != Created {
			panic("expected Created state")
		}
	})

	// Test 2: full valid lifecycle Created -> Confirmed -> Shipped -> Delivered
	test("test_full_lifecycle", func() {
		ResetManager()
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		if !ConfirmOrder(id) {
			panic("confirm should return true")
		}
		if GetOrderState(id) != Confirmed {
			panic("expected Confirmed")
		}
		if !ShipOrder(id) {
			panic("ship should return true")
		}
		if GetOrderState(id) != Shipped {
			panic("expected Shipped")
		}
		if !DeliverOrder(id) {
			panic("deliver should return true")
		}
		if GetOrderState(id) != Delivered {
			panic("expected Delivered")
		}
	})

	// Test 3: invalid transition — skip from Created to Shipped
	test("test_invalid_skip_to_shipped", func() {
		ResetManager()
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		if ShipOrder(id) {
			panic("ship from Created should return false")
		}
		if GetOrderState(id) != Created {
			panic("state should remain Created")
		}
	})

	// Test 4: invalid transition — skip from Created to Delivered
	test("test_invalid_skip_to_delivered", func() {
		ResetManager()
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		if DeliverOrder(id) {
			panic("deliver from Created should return false")
		}
		if GetOrderState(id) != Created {
			panic("state should remain Created")
		}
	})

	// Test 5: invalid transition — backward from Shipped to Confirmed
	test("test_invalid_backward_transition", func() {
		ResetManager()
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		ConfirmOrder(id)
		ShipOrder(id)
		if ConfirmOrder(id) {
			panic("confirm from Shipped should return false")
		}
		if GetOrderState(id) != Shipped {
			panic("state should remain Shipped")
		}
	})

	// Test 6: multiple orders are independent
	test("test_multiple_orders_independent", func() {
		ResetManager()
		id1 := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		id2 := CreateOrder([]OrderItem{{ProductId: "PROD-2", Quantity: 1}}, 200.0)
		ConfirmOrder(id1)
		if GetOrderState(id1) != Confirmed {
			panic("id1 should be Confirmed")
		}
		if GetOrderState(id2) != Created {
			panic("id2 should still be Created")
		}
	})

	// Test 7: confirm non-existent order returns false
	test("test_nonexistent_order", func() {
		ResetManager()
		if ConfirmOrder("NONEXISTENT") {
			panic("confirm nonexistent should return false")
		}
	})

	fmt.Printf("PART1_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
