package main

import (
	"fmt"
	"time"
)

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

// ─── Observer Interface ───────────────────────────────────────────────────────

type OrderObserver interface {
	OnStateChange(orderId string, from, to OrderState)
}

// ─── OrderManager ─────────────────────────────────────────────────────────────

type OrderManager struct {
	orders    map[string]*Order
	inventory map[string]int
	history   map[string][]StateTransition
	observers []OrderObserver
	nextId    int
}

func NewOrderManager() *OrderManager {
	return &OrderManager{
		orders:    make(map[string]*Order),
		inventory: make(map[string]int),
		history:   make(map[string][]StateTransition),
		nextId:    1,
	}
}

func nowMs() int64 {
	return time.Now().UnixMilli()
}

func (m *OrderManager) notifyObservers(orderId string, from, to OrderState) {
	for _, obs := range m.observers {
		obs.OnStateChange(orderId, from, to)
	}
}

func (m *OrderManager) transition(orderId string, expected, next OrderState) bool {
	o, ok := m.orders[orderId]
	if !ok || o.State != expected {
		return false
	}
	from := o.State
	o.State = next

	// TODO: Record transition in history with timestamp
	m.history[orderId] = append(m.history[orderId], StateTransition{
		FromState: from,
		ToState:   next,
		Timestamp: nowMs(),
	})
	// TODO: Notify all observers
	m.notifyObservers(orderId, from, next)
	return true
}

func (m *OrderManager) SetInventory(productId string, qty int) {
	m.inventory[productId] = qty
}
func (m *OrderManager) GetInventory(productId string) int { return m.inventory[productId] }

func (m *OrderManager) CreateOrder(items []OrderItem, totalAmount float64) string {
	id := fmt.Sprintf("ORD-%d", m.nextId)
	m.nextId++
	for _, item := range items {
		m.inventory[item.ProductId] -= item.Quantity
	}
	m.orders[id] = &Order{Id: id, Items: items, TotalAmount: totalAmount, State: Created}
	// TODO: Record initial history entry for creation
	m.history[id] = append(m.history[id], StateTransition{
		FromState: Created,
		ToState:   Created,
		Timestamp: nowMs(),
	})
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
	if o.State != Created && o.State != Confirmed {
		return false
	}
	from := o.State
	for _, item := range o.Items {
		m.inventory[item.ProductId] += item.Quantity
	}
	o.State = Cancelled
	// TODO: Record cancellation in history
	m.history[orderId] = append(m.history[orderId], StateTransition{
		FromState: from,
		ToState:   Cancelled,
		Timestamp: nowMs(),
	})
	// TODO: Notify observers
	m.notifyObservers(orderId, from, Cancelled)
	return true
}

func (m *OrderManager) GetOrderState(orderId string) OrderState { return m.orders[orderId].State }

func (m *OrderManager) GetOrderHistory(orderId string) []StateTransition {
	// TODO: Return history for this order (empty slice if not found)
	return m.history[orderId]
}

func (m *OrderManager) AddObserver(obs OrderObserver) {
	m.observers = append(m.observers, obs)
}

// ─── Global Instance + Entry Points ──────────────────────────────────────────

var manager = NewOrderManager()

func CreateOrder(items []OrderItem, totalAmount float64) string {
	return manager.CreateOrder(items, totalAmount)
}
func ConfirmOrder(orderId string) bool                 { return manager.ConfirmOrder(orderId) }
func ShipOrder(orderId string) bool                    { return manager.ShipOrder(orderId) }
func DeliverOrder(orderId string) bool                 { return manager.DeliverOrder(orderId) }
func CancelOrder(orderId string) bool                  { return manager.CancelOrder(orderId) }
func GetOrderState(orderId string) OrderState          { return manager.GetOrderState(orderId) }
func SetInventory(productId string, qty int)           { manager.SetInventory(productId, qty) }
func GetInventory(productId string) int                { return manager.GetInventory(productId) }
func GetOrderHistory(orderId string) []StateTransition { return manager.GetOrderHistory(orderId) }
func AddObserver(obs OrderObserver)                    { manager.AddObserver(obs) }
func ResetManager()                                    { manager = NewOrderManager() }
