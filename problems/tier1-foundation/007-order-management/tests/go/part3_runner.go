package main

import "fmt"

// TestObserver records notifications for assertions
type TestObserver struct {
	Notifications []struct {
		OrderId string
		From    OrderState
		To      OrderState
	}
}

func (o *TestObserver) OnStateChange(orderId string, from, to OrderState) {
	o.Notifications = append(o.Notifications, struct {
		OrderId string
		From    OrderState
		To      OrderState
	}{orderId, from, to})
}

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

	// Test 1: history tracks full lifecycle
	test("test_history_full_lifecycle", func() {
		ResetManager()
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		ConfirmOrder(id)
		ShipOrder(id)
		DeliverOrder(id)
		hist := GetOrderHistory(id)
		// Should have: creation entry + 3 transitions = 4 entries
		if len(hist) != 4 {
			panic(fmt.Sprintf("expected 4 history entries, got %d", len(hist)))
		}
		if hist[1].FromState != Created || hist[1].ToState != Confirmed {
			panic("hist[1] should be Created->Confirmed")
		}
		if hist[2].FromState != Confirmed || hist[2].ToState != Shipped {
			panic("hist[2] should be Confirmed->Shipped")
		}
		if hist[3].FromState != Shipped || hist[3].ToState != Delivered {
			panic("hist[3] should be Shipped->Delivered")
		}
	})

	// Test 2: history timestamps are non-decreasing
	test("test_history_timestamps_ordered", func() {
		ResetManager()
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		ConfirmOrder(id)
		ShipOrder(id)
		hist := GetOrderHistory(id)
		for i := 1; i < len(hist); i++ {
			if hist[i].Timestamp < hist[i-1].Timestamp {
				panic(fmt.Sprintf("timestamps out of order at index %d", i))
			}
		}
	})

	// Test 3: failed transitions do not appear in history
	test("test_failed_transition_no_history", func() {
		ResetManager()
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		ShipOrder(id) // invalid — should fail
		hist := GetOrderHistory(id)
		if len(hist) != 1 {
			panic(fmt.Sprintf("expected 1 history entry, got %d", len(hist)))
		}
	})

	// Test 4: observer is notified on valid transitions
	test("test_observer_notified", func() {
		ResetManager()
		obs := &TestObserver{}
		AddObserver(obs)
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		ConfirmOrder(id)
		ShipOrder(id)
		if len(obs.Notifications) != 2 {
			panic(fmt.Sprintf("expected 2 notifications, got %d", len(obs.Notifications)))
		}
		if obs.Notifications[0].From != Created || obs.Notifications[0].To != Confirmed {
			panic("first notification should be Created->Confirmed")
		}
		if obs.Notifications[1].From != Confirmed || obs.Notifications[1].To != Shipped {
			panic("second notification should be Confirmed->Shipped")
		}
	})

	// Test 5: observer NOT notified on failed transitions
	test("test_observer_not_notified_on_failure", func() {
		ResetManager()
		obs := &TestObserver{}
		AddObserver(obs)
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		ShipOrder(id) // invalid
		if len(obs.Notifications) != 0 {
			panic(fmt.Sprintf("expected 0 notifications, got %d", len(obs.Notifications)))
		}
	})

	// Test 6: cancellation appears in history
	test("test_cancel_in_history", func() {
		ResetManager()
		SetInventory("PROD-1", 10)
		id := CreateOrder([]OrderItem{{ProductId: "PROD-1", Quantity: 1}}, 100.0)
		CancelOrder(id)
		hist := GetOrderHistory(id)
		if len(hist) != 2 {
			panic(fmt.Sprintf("expected 2 history entries, got %d", len(hist)))
		}
		if hist[1].FromState != Created || hist[1].ToState != Cancelled {
			panic("hist[1] should be Created->Cancelled")
		}
	})

	fmt.Printf("PART3_SUMMARY %d/%d\n", passed, passed+failed)
	return failed
}
