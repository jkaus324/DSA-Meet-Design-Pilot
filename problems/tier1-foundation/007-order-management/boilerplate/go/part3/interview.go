package main

// ─── Data Model (given — do not modify) ──────────────────────────────────────

type OrderState int

const (
	Created   OrderState = iota
	Confirmed OrderState = iota
	Shipped   OrderState = iota
	Delivered OrderState = iota
	Cancelled OrderState = iota
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

type StateTransition struct {
	FromState OrderState
	ToState   OrderState
	Timestamp int64
}

type OrderObserver interface {
	OnStateChange(orderId string, from, to OrderState)
}

// ─── NEW in Extension 2 ───────────────────────────────────────────────────────
//
// Track the full transition history for every order.
// Notify registered observers on every successful state change.
//
// Think about:
//   - Where do you store per-order history? (map of slices)
//   - How do you decouple notification logic from the state machine?
//   - What if you want logging, analytics, AND alerts — all independently?
//
// Entry points (must exist for tests — all previous plus):
//   func GetOrderHistory(orderId string) []StateTransition
//   func AddObserver(obs OrderObserver)
//
// ─────────────────────────────────────────────────────────────────────────────

func CreateOrder(items []OrderItem, totalAmount float64) string { return "" }
func ConfirmOrder(orderId string) bool                          { return false }
func ShipOrder(orderId string) bool                             { return false }
func DeliverOrder(orderId string) bool                          { return false }
func CancelOrder(orderId string) bool                           { return false }
func GetOrderState(orderId string) OrderState                   { return Created }
func SetInventory(productId string, quantity int)               {}
func GetInventory(productId string) int                         { return 0 }
func GetOrderHistory(orderId string) []StateTransition          { return nil }
func AddObserver(obs OrderObserver)                             {}
func ResetManager()                                             {}
