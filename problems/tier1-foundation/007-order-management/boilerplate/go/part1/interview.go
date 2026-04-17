package main

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

// ─── Your Design Starts Here ──────────────────────────────────────────────────
//
// Design and implement an OrderManager that:
//   1. Creates orders and stores them in a map
//   2. Enforces valid state transitions:
//      Created -> Confirmed -> Shipped -> Delivered
//   3. Rejects any invalid transition (no skipping, no backward)
//
// Think about:
//   - How do you validate that a transition is legal?
//   - What data structure gives O(1) order lookup by ID?
//   - What happens when Extension 1 (cancellation) is added?
//
// Entry points (must exist for tests):
//   func CreateOrder(items []OrderItem, totalAmount float64) string
//   func ConfirmOrder(orderId string) bool
//   func ShipOrder(orderId string) bool
//   func DeliverOrder(orderId string) bool
//   func GetOrderState(orderId string) OrderState
//
// ─────────────────────────────────────────────────────────────────────────────

func CreateOrder(items []OrderItem, totalAmount float64) string {
	return ""
}

func ConfirmOrder(orderId string) bool {
	return false
}

func ShipOrder(orderId string) bool {
	return false
}

func DeliverOrder(orderId string) bool {
	return false
}

func GetOrderState(orderId string) OrderState {
	return Created
}

func ResetManager() {}
