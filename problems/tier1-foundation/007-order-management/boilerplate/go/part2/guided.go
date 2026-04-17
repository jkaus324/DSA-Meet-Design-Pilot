package main

import "fmt"

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

// ─── OrderManager ─────────────────────────────────────────────────────────────
// HINT: You now need TWO maps:
//   - orders: map[string]*Order     — order lookup by ID
//   - inventory: map[string]int     — stock count by product ID
//
// On CreateOrder: decrement inventory for each item
// On CancelOrder: restore inventory for each item (only if state is Created or Confirmed)

// HINT: CancelOrder should:
//   1. Check the order exists
//   2. Check state is Created or Confirmed
//   3. Iterate through order items and restore inventory
//   4. Set state to Cancelled

type OrderManager struct {
	orders    map[string]*Order
	inventory map[string]int
	nextId    int
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		orders:    make(map[string]*Order),
		inventory: make(map[string]int),
		nextId:    1,
	}
}

func (m *OrderManager) transition(orderId string, expected, next OrderState) bool {
	o, ok := m.orders[orderId]
	if !ok || o.State != expected {
		return false
	}
	o.State = next
	return true
}

func (m *OrderManager) SetInventory(productId string, qty int) {
	m.inventory[productId] = qty
}

func (m *OrderManager) GetInventory(productId string) int {
	return m.inventory[productId]
}

func (m *OrderManager) CreateOrder(items []OrderItem, totalAmount float64) string {
	id := fmt.Sprintf("ORD-%d", m.nextId)
	m.nextId++
	// TODO: Decrement inventory for each item
	// TODO: Store the order in the map with state Created
	return id
}

func (m *OrderManager) ConfirmOrder(orderId string) bool {
	return m.transition(orderId, Created, Confirmed)
}

func (m *OrderManager) ShipOrder(orderId string) bool {
	return m.transition(orderId, Confirmed, Shipped)
}

func (m *OrderManager) DeliverOrder(orderId string) bool {
	return m.transition(orderId, Shipped, Delivered)
}

func (m *OrderManager) CancelOrder(orderId string) bool {
	o, ok := m.orders[orderId]
	if !ok {
		return false
	}
	// TODO: Check if state is Created or Confirmed (otherwise return false)
	// TODO: Iterate through order items and restore inventory
	// TODO: Set state to Cancelled
	_ = o
	return false
}

func (m *OrderManager) GetOrderState(orderId string) OrderState {
	return m.orders[orderId].State
}

// ─── Global Instance + Entry Points ──────────────────────────────────────────

var manager = NewOrderManager()

func CreateOrder(items []OrderItem, totalAmount float64) string {
	return manager.CreateOrder(items, totalAmount)
}
func ConfirmOrder(orderId string) bool        { return manager.ConfirmOrder(orderId) }
func ShipOrder(orderId string) bool           { return manager.ShipOrder(orderId) }
func DeliverOrder(orderId string) bool        { return manager.DeliverOrder(orderId) }
func CancelOrder(orderId string) bool         { return manager.CancelOrder(orderId) }
func GetOrderState(orderId string) OrderState { return manager.GetOrderState(orderId) }
func SetInventory(productId string, qty int)  { manager.SetInventory(productId, qty) }
func GetInventory(productId string) int       { return manager.GetInventory(productId) }
func ResetManager()                           { manager = NewOrderManager() }
