package main

import "fmt"

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type OrderState int

const (
	Created   OrderState = iota
	Confirmed OrderState = iota
	Shipped   OrderState = iota
	Delivered OrderState = iota
)

type OrderItem struct {
	ProductId string
	Quantity  int
}

type Order struct {
	Id          string
	Items       []OrderItem
	TotalAmount float64
	State       OrderState
}

// ─── OrderManager ─────────────────────────────────────────────────────────────
// HINT: Use a map to store orders by ID.
// Each transition method should check the current state before changing it.

// type OrderManager struct {
//     orders map[string]*Order
//     nextId int
// }
//
// func (m *OrderManager) CreateOrder(items []OrderItem, totalAmount float64) string
// func (m *OrderManager) ConfirmOrder(orderId string) bool  // Created -> Confirmed
// func (m *OrderManager) ShipOrder(orderId string) bool     // Confirmed -> Shipped
// func (m *OrderManager) DeliverOrder(orderId string) bool  // Shipped -> Delivered
// func (m *OrderManager) GetOrderState(orderId string) OrderState
//
// HINT: Consider a private helper method like:
//   func (m *OrderManager) transition(orderId string, expected, next OrderState) bool
// This avoids duplicating the "check current state, then update" logic.

// ─── Test Entry Points ────────────────────────────────────────────────────────

var globalManager = &struct {
	orders map[string]*Order
	nextId int
}{orders: make(map[string]*Order), nextId: 1}

func CreateOrder(items []OrderItem, totalAmount float64) string {
	id := fmt.Sprintf("ORD-%d", globalManager.nextId)
	globalManager.nextId++
	// TODO: Create an Order with state Created and store it in the map
	_ = id
	return ""
}

func ConfirmOrder(orderId string) bool {
	// TODO: transition from Created to Confirmed
	return false
}

func ShipOrder(orderId string) bool {
	// TODO: transition from Confirmed to Shipped
	return false
}

func DeliverOrder(orderId string) bool {
	// TODO: transition from Shipped to Delivered
	return false
}

func GetOrderState(orderId string) OrderState {
	if o, ok := globalManager.orders[orderId]; ok {
		return o.State
	}
	return Created
}

func ResetManager() {
	globalManager.orders = make(map[string]*Order)
	globalManager.nextId = 1
}
