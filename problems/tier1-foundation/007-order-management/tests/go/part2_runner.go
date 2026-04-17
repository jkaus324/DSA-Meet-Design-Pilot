package main

import "fmt"

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

	// Test 1: cancel from Created state succeeds
	test("test_cancel_from_created", func() {
		ResetManager()
		SetInventory("PROD-1", 10)
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 3}}, 300.0)
		if GetInventory("PROD-1") != 7 {
			panic("inventory should be decremented on create")
		}
		if !CancelOrder(id) {
			panic("cancel from Created should succeed")
		}
		if GetOrderState(id) != Cancelled {
			panic("state should be Cancelled")
		}
		if GetInventory("PROD-1") != 10 {
			panic("inventory should be restored on cancel")
		}
	})

	// Test 2: cancel from Confirmed state succeeds
	test("test_cancel_from_confirmed", func() {
		ResetManager()
		SetInventory("PROD-1", 10)
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 2}}, 200.0)
		ConfirmOrder(id)
		if !CancelOrder(id) {
			panic("cancel from Confirmed should succeed")
		}
		if GetOrderState(id) != Cancelled {
			panic("state should be Cancelled")
		}
		if GetInventory("PROD-1") != 10 {
			panic("inventory should be restored")
		}
	})

	// Test 3: cancel from Shipped state fails
	test("test_cancel_from_shipped_fails", func() {
		ResetManager()
		SetInventory("PROD-1", 10)
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 2}}, 200.0)
		ConfirmOrder(id)
		ShipOrder(id)
		if CancelOrder(id) {
			panic("cancel from Shipped should fail")
		}
		if GetOrderState(id) != Shipped {
			panic("state should remain Shipped")
		}
		if GetInventory("PROD-1") != 8 {
			panic("inventory should not be restored")
		}
	})

	// Test 4: cancel from Delivered state fails
	test("test_cancel_from_delivered_fails", func() {
		ResetManager()
		SetInventory("PROD-1", 10)
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		ConfirmOrder(id)
		ShipOrder(id)
		DeliverOrder(id)
		if CancelOrder(id) {
			panic("cancel from Delivered should fail")
		}
		if GetOrderState(id) != Delivered {
			panic("state should remain Delivered")
		}
	})

	// Test 5: cancel restores inventory for multiple items
	test("test_cancel_multi_item_inventory", func() {
		ResetManager()
		SetInventory("PROD-A", 20)
		SetInventory("PROD-B", 15)
		id := CreateOrder([]OrderItem{
			{ProductId: "PROD-A", Quantity: 5},
			{ProductId: "PROD-B", Quantity: 3},
		}, 800.0)
		if GetInventory("PROD-A") != 15 {
			panic("PROD-A inventory should be 15")
		}
		if GetInventory("PROD-B") != 12 {
			panic("PROD-B inventory should be 12")
		}
		CancelOrder(id)
		if GetInventory("PROD-A") != 20 {
			panic("PROD-A inventory should be restored to 20")
		}
		if GetInventory("PROD-B") != 15 {
			panic("PROD-B inventory should be restored to 15")
		}
	})

	// Test 6: cannot cancel an already cancelled order
	test("test_double_cancel_fails", func() {
		ResetManager()
		SetInventory("PROD-1", 10)
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 2}}, 200.0)
		CancelOrder(id)
		if CancelOrder(id) {
			panic("double cancel should fail")
		}
	})

	fmt.Printf("PART2_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
