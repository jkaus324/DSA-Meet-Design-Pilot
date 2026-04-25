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

// ─── NEW in Extension 1 ───────────────────────────────────────────────────────
//
// Orders can now be CANCELLED from Created or Confirmed state.
// Cancellation from Shipped or Delivered is NOT allowed.
// When cancelled, inventory must be released (restored).
//
// Think about:
//   - How do you track inventory per product? (map)
//   - When should inventory be decremented? On order creation.
//   - When should inventory be restored? On cancellation.
//   - What if the order has multiple items?
//
// Entry points (must exist for tests):
//   func CreateOrder(items []OrderItem, totalAmount float64) string
//   func ConfirmOrder(orderId string) bool
//   func ShipOrder(orderId string) bool
//   func DeliverOrder(orderId string) bool
//   func CancelOrder(orderId string) bool
//   func GetOrderState(orderId string) OrderState
//   func SetInventory(productId string, quantity int)
//   func GetInventory(productId string) int
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
func ResetManager()                                             {}
