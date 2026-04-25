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

type OrderManager struct {
	orders map[string]*Order
	nextId int
}

func NewOrderManager() *OrderManager {
	return &OrderManager{orders: make(map[string]*Order), nextId: 1}
}

func (m *OrderManager) transition(orderId string, expected, next OrderState) bool {
	o, ok := m.orders[orderId]
	if !ok {
		return false
	}
	// TODO: Check if current state matches 'expected'
	// If yes, update to 'next' and return true
	// If no, return false (invalid transition)
	_ = o
	return false
}

func (m *OrderManager) CreateOrder(items []OrderItem, totalAmount float64) string {
	id := fmt.Sprintf("ORD-%d", m.nextId)
	m.nextId++
	// TODO: Create an Order with state Created and store it in the map
	return id
}

func (m *OrderManager) ConfirmOrder(orderId string) bool {
	// TODO: transition from Created to Confirmed
	return false
}

func (m *OrderManager) ShipOrder(orderId string) bool {
	// TODO: transition from Confirmed to Shipped
	return false
}

func (m *OrderManager) DeliverOrder(orderId string) bool {
	// TODO: transition from Shipped to Delivered
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

func ConfirmOrder(orderId string) bool {
	return manager.ConfirmOrder(orderId)
}

func ShipOrder(orderId string) bool {
	return manager.ShipOrder(orderId)
}

func DeliverOrder(orderId string) bool {
	return manager.DeliverOrder(orderId)
}

func GetOrderState(orderId string) OrderState {
	return manager.GetOrderState(orderId)
}

func ResetManager() {
	manager = NewOrderManager()
}
